package main

import (
	"database/sql"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

// ImportGPUsFromCSV reads a CSV file and inserts data into the GPUs table.
func ImportGPUsFromCSV(db *sql.DB, filePath string, overwrite bool) error {
	if overwrite {
		_, err := db.Exec("TRUNCATE TABLE gpus")
		if err != nil {
			return fmt.Errorf("failed to truncate GPUs table: %v", err)
		}
		log.Println("GPUs table truncated successfully")
	}

	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to read CSV file: %v", err)
	}

	// Skip the header row and insert each record
	for i, record := range records {
		if i == 0 {
			continue // Skip header
		}
		if len(record) != 5 {
			return fmt.Errorf("invalid record format at line %d", i+1)
		}
		log.Println(record)
		// Insert into the GPUs table
		_, err := db.Exec(`
            INSERT INTO gpus (id, server_name, gpu_number, manufacturer, model_name, vram_size_mb)
            VALUES (MD5(CONCAT(?, ?)), ?, ?, ?, ?, ?)`,
			record[0], record[1], record[0], record[1], record[2], record[3], record[4],
		)
		if err != nil {
			return fmt.Errorf("failed to insert record at line %d: %v", i+1, err)
		}
	}
	log.Println("GPUs data imported successfully")
	return nil
}

// ImportUsersFromCSV reads a CSV file and inserts data into the Users table.
func ImportUsersFromCSV(db *sql.DB, filePath string, overwrite bool) error {
	if overwrite {
		_, err := db.Exec("TRUNCATE TABLE users")
		if err != nil {
			return fmt.Errorf("failed to truncate Users table: %v", err)
		}
		log.Println("Users table truncated successfully")
	}

	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to read CSV file: %v", err)
	}

	// Skip the header row and insert each record
	for i, record := range records {
		if i == 0 {
			continue // Skip header
		}
		if len(record) != 3 {
			return fmt.Errorf("invalid record format at line %d", i+1)
		}

		// Insert into the Users table
		_, err := db.Exec(`
            INSERT INTO users (email, name, password, is_admin, is_whitelisted)
            VALUES (?, ?, ?, ?, ?)`,
			record[0], record[1], record[2], false, false,
		)
		if err != nil {
			return fmt.Errorf("failed to insert record at line %d: %v", i+1, err)
		}
	}
	log.Println("Users data imported successfully")
	return nil
}

func main() {
	// Define command-line flags
	filePath := flag.String("file", "", "Path to the CSV file")
	table := flag.String("table", "", "Table to update (gpus or users)")
	mode := flag.String("mode", "insert", "Mode of operation (insert or overwrite)")
	dsn := flag.String("dsn", "user:password@tcp(localhost:3306)/gpu_scheduler", "Database DSN")

	flag.Parse()

	// Validate flags
	if *filePath == "" || *table == "" {
		log.Fatalf("Usage: go run csv_import.go -file=<path_to_csv> -table=<gpus|users> -mode=<insert|overwrite>")
	}

	// Connect to the database
	db, err := sql.Open("mysql", *dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Determine the mode
	overwrite := *mode == "overwrite"

	// Call the appropriate function based on the table flag
	switch *table {
	case "gpus":
		err = ImportGPUsFromCSV(db, *filePath, overwrite)
	case "users":
		err = ImportUsersFromCSV(db, *filePath, overwrite)
	default:
		log.Fatalf("Invalid table specified: %s. Use 'gpus' or 'users'.", *table)
	}

	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	log.Println("Operation completed successfully")
}
