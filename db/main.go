package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/eduardo-escoto/gpu_request/db/importers"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// Define command-line flags
	table := flag.String("table", "", "Table to update (gpus or users)")
	mode := flag.String("mode", "update", "Mode of operation (remake, update, insert, or remake_for_server)")
	dsn := flag.String("dsn", "", "Database DSN (can also be set via the DATABASE_DSN environment variable)")
	verbose := flag.Bool("verbose", false, "Enable verbose logging")
	fileID := flag.String("file-id", "", "Google Drive file ID for importing users (only applicable for the 'users' table)")

	flag.Parse()

	// Validate flags
	if *table == "" {
		log.Fatalf("Usage: go run main.go -table=<gpus|users> -mode=<remake|update|insert|remake_for_server> -dsn=<database_dsn> [-file-id=<file_id>] [-verbose]")
	}

	// Restrict `remake_for_server` mode to the `gpus` table
	if *mode == "remake_for_server" && *table != "gpus" {
		log.Fatalf("The 'remake_for_server' mode is only valid for the 'gpus' table.")
	}

	// Check for DATABASE_DSN environment variable if DSN is not provided as a flag
	if *dsn == "" {
		*dsn = os.Getenv("DATABASE_DSN")
		if *dsn == "" {
			log.Fatalf("No DSN provided. Use the -dsn flag or set the DATABASE_DSN environment variable.")
		}
	}

	// Connect to the database
	db, err := sql.Open("mysql", *dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Call the appropriate function based on the table flag
	switch *table {
	case "gpus":
		err = importers.ImportGPUsFromNvidiaSMI(db, *mode, *verbose)
	case "users":
		// Ensure the file ID is provided for the users table
		if *fileID == "" {
			log.Fatalf("File ID must be provided when the table is 'users'. Use the -file-id flag.")
		}

		// Construct the Google Sheets export URL
		csvURL := fmt.Sprintf("https://docs.google.com/spreadsheets/d/%s/gviz/tq?tqx=out:csv", *fileID)

		// Import users from the constructed CSV URL
		err = importers.ImportSurveyResponsesFromCSV(db, csvURL, *mode, *verbose)
		if err != nil {
			log.Fatalf("Error: %v", err)
		}
		err = importers.UpdateUsersFromSurveyResponses(db, *verbose)
	default:
		log.Fatalf("Invalid table specified: %s. Use 'gpus' or 'users'.", *table)
	}

	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	log.Println("Operation completed successfully")
}
