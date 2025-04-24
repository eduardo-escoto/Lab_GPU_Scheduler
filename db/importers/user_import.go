package importers

import (
	"database/sql"
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

// ImportSurveyResponsesFromCSV downloads a CSV from a URL, parses it, and populates the survey_responses table.
func ImportSurveyResponsesFromCSV(db *sql.DB, csvURL string, mode string, verbose bool) error {
	// Step 1: Download the CSV file
	resp, err := http.Get(csvURL)
	if err != nil {
		return fmt.Errorf("failed to download CSV: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download CSV: received status code %d", resp.StatusCode)
	}

	// Step 2: Parse the CSV file
	reader := csv.NewReader(resp.Body)
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to parse CSV: %w", err)
	}

	if len(records) < 2 {
		return errors.New("CSV file is empty or missing data rows")
	}

	// Step 3: Validate and insert data into the database
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Handle modes
	switch mode {
	case "remake":
		if verbose {
			log.Println("Clearing survey_responses table...")
		}
		_, err := tx.Exec("DELETE FROM survey_responses")
		if err != nil {
			return fmt.Errorf("failed to clear survey_responses table: %w", err)
		}
		if verbose {
			log.Println("survey_responses table cleared successfully.")
		}

	case "update":
		if verbose {
			log.Println("Running in update mode...")
		}

	case "insert":
		if verbose {
			log.Println("Running in insert mode...")
		}

	default:
		return fmt.Errorf("invalid mode: %s. Use 'remake', 'update', or 'insert'", mode)
	}

	// Prepare statement for inserting survey responses
	insertStmt, err := tx.Prepare(`
        INSERT INTO survey_responses (
            email, full_name, desired_username, ssh_key, remark, user_type, lab_join_year, submitted_at, 
            granted_access_at, revoked_access_at, revoke_scheduled_at, approving_party, revoking_party
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
        ON DUPLICATE KEY UPDATE
            full_name = VALUES(full_name),
            desired_username = VALUES(desired_username),
            ssh_key = VALUES(ssh_key),
            remark = VALUES(remark),
            user_type = VALUES(user_type),
            lab_join_year = VALUES(lab_join_year),
            submitted_at = VALUES(submitted_at),
            granted_access_at = VALUES(granted_access_at),
            revoked_access_at = VALUES(revoked_access_at),
            revoke_scheduled_at = VALUES(revoke_scheduled_at),
            approving_party = VALUES(approving_party),
            revoking_party = VALUES(revoking_party)
    `)
	if err != nil {
		return fmt.Errorf("failed to prepare insert statement: %w", err)
	}
	defer insertStmt.Close()

	// Helper function to parse and format datetime values
	parseDatetime := func(value string) (interface{}, error) {
		if value == "" {
			return nil, nil // Return NULL for empty values
		}
		// Try parsing with non-zero-padded format
		parsedTime, err := time.Parse("1/2/2006 15:04:05", value)
		if err != nil {
			parsedTime, err = time.Parse("01/02/2006 15:04:05", value)
			if err != nil {
				return nil, fmt.Errorf("failed to parse datetime: %w", err)
			}
			return parsedTime.Format("2006-01-02 15:04:05"), nil
		}
		return parsedTime.Format("2006-01-02 15:04:05"), nil
	}

	// Helper function to normalize string values
	normalizeValue := func(value string) interface{} {
		if value == "" {
			return nil // Return NULL for empty values
		}
		return value
	}

	// Process each record, skipping the header row
	for i := 1; i < len(records); i++ {
		record := records[i]
		if len(record) < 13 { // Ensure there are enough columns
			return fmt.Errorf("invalid record on line %d: %v", i+1, record)
		}

		// Extract and normalize fields from the CSV
		submittedAt, err := parseDatetime(strings.TrimSpace(record[0]))
		if err != nil {
			return fmt.Errorf("failed to parse submitted_at on line %d: %w", i+1, err)
		}

		fullName := normalizeValue(strings.TrimSpace(record[1]))
		email := normalizeValue(strings.TrimSpace(record[2]))
		desiredUsername := normalizeValue(strings.TrimSpace(record[3]))
		sshKey := normalizeValue(strings.TrimSpace(record[4]))
		remark := normalizeValue(strings.TrimSpace(record[5]))
		userType := normalizeValue(strings.TrimSpace(record[6]))
		labJoinYear := normalizeValue(strings.TrimSpace(record[7]))

		grantedAccessAt, err := parseDatetime(strings.TrimSpace(record[8]))
		if err != nil {
			return fmt.Errorf("failed to parse granted_access_at on line %d: %w", i+1, err)
		}

		revokedAccessAt, err := parseDatetime(strings.TrimSpace(record[9]))
		if err != nil {
			return fmt.Errorf("failed to parse revoked_access_at on line %d: %w", i+1, err)
		}

		revokeScheduledAt, err := parseDatetime(strings.TrimSpace(record[10]))
		if err != nil {
			return fmt.Errorf("failed to parse revoke_scheduled_at on line %d: %w", i+1, err)
		}

		approvingParty := normalizeValue(strings.TrimSpace(record[11]))
		revokingParty := normalizeValue(strings.TrimSpace(record[12]))

		// Insert or update survey response
		_, err = insertStmt.Exec(email, fullName, desiredUsername, sshKey, remark, userType, labJoinYear, submittedAt,
			grantedAccessAt, revokedAccessAt, revokeScheduledAt, approvingParty, revokingParty)
		if err != nil {
			return fmt.Errorf("failed to insert/update survey response on line %d: %w", i+1, err)
		}

		if verbose {
			log.Printf("Processed survey response for email: %s", email)
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	log.Println("Survey responses imported successfully.")
	return nil
}

// UpdateUsersFromSurveyResponses updates the users table based on the latest survey responses.
func UpdateUsersFromSurveyResponses(db *sql.DB, verbose bool) error {
	// Query to get the relevant values from survey_responses
	query := `
        WITH ranked_responses AS (
            SELECT *,
                   ROW_NUMBER() OVER (PARTITION BY email ORDER BY submitted_at DESC) AS survey_rn,
                   FIRST_VALUE(submitted_at) OVER (PARTITION BY email ORDER BY submitted_at ASC) AS first_submission
            FROM survey_responses
        )
        SELECT
            email,
            desired_username AS user_name,
            NULL AS password,
            remark AS comment,
            user_type,
            lab_join_year,
            first_submission AS access_survey_submitted_at,
            submitted_at AS access_survey_updated_at,
            granted_access_at,
            revoked_access_at,
            revoke_scheduled_at,
            approving_party,
            revoking_party,
            full_name AS name, -- Include the name field
            FALSE AS is_admin,
            TRUE AS is_whitelisted
        FROM ranked_responses
        WHERE survey_rn = 1;
    `

	// Execute the query
	rows, err := db.Query(query)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	// Prepare the insert/update statement for the users table
	insertStmt, err := db.Prepare(`
        INSERT INTO users (
            email, user_name, password, comment, user_type, lab_join_year, access_survey_submitted_at, 
            access_survey_updated_at, granted_access_at, revoked_access_at, revoke_scheduled_at, 
            approving_party, revoking_party, name, is_admin, is_whitelisted
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
        ON DUPLICATE KEY UPDATE
            user_name = VALUES(user_name),
            password = VALUES(password),
            comment = VALUES(comment),
            user_type = VALUES(user_type),
            lab_join_year = VALUES(lab_join_year),
            access_survey_submitted_at = VALUES(access_survey_submitted_at),
            access_survey_updated_at = VALUES(access_survey_updated_at),
            granted_access_at = VALUES(granted_access_at),
            revoked_access_at = VALUES(revoked_access_at),
            revoke_scheduled_at = VALUES(revoke_scheduled_at),
            approving_party = VALUES(approving_party),
            revoking_party = VALUES(revoking_party),
            name = VALUES(name),
            is_admin = VALUES(is_admin),
            is_whitelisted = VALUES(is_whitelisted)
    `)
	if err != nil {
		return fmt.Errorf("failed to prepare insert statement for users table: %w", err)
	}
	defer insertStmt.Close()

	// Iterate through the query results and insert/update the users table
	for rows.Next() {
		var (
			email, userName, password, comment, userType, approvingParty, revokingParty, name                   sql.NullString
			labJoinYear                                                                                         sql.NullInt64
			accessSurveySubmittedAt, accessSurveyUpdatedAt, grantedAccessAt, revokedAccessAt, revokeScheduledAt sql.NullTime
			isAdmin, isWhitelisted                                                                              bool
		)

		err := rows.Scan(
			&email, &userName, &password, &comment, &userType, &labJoinYear,
			&accessSurveySubmittedAt, &accessSurveyUpdatedAt, &grantedAccessAt,
			&revokedAccessAt, &revokeScheduledAt, &approvingParty, &revokingParty,
			&name, &isAdmin, &isWhitelisted,
		)
		if err != nil {
			return fmt.Errorf("failed to scan row: %w", err)
		}

		_, err = insertStmt.Exec(
			email, userName, password, comment, userType, labJoinYear,
			accessSurveySubmittedAt, accessSurveyUpdatedAt, grantedAccessAt,
			revokedAccessAt, revokeScheduledAt, approvingParty, revokingParty,
			name, isAdmin, isWhitelisted,
		)
		if err != nil {
			return fmt.Errorf("failed to insert/update user: %w", err)
		}

		if verbose {
			log.Printf("Processed user with email: %s", email.String)
		}
	}

	// Check for errors during iteration
	if err := rows.Err(); err != nil {
		return fmt.Errorf("error iterating over rows: %w", err)
	}

	log.Println("Users table updated successfully.")
	return nil
}
