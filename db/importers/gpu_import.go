package importers

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

// deleteDependentRecords deletes records from all dependent tables based on the provided WHERE clause.
func deleteDependentRecords(db *sql.DB, whereClause string, args ...interface{}) error {
	// Tables to delete from, in order of dependency
	tables := []string{
		"request_gpu_assignments",
		"real_time_usage",
		"gpu_processes",
		"real_time_usage_hourly_historical",
		"gpu_processes_hourly_historical",
	}

	for _, table := range tables {
		query := fmt.Sprintf("DELETE FROM %s %s", table, whereClause)
		_, err := db.Exec(query, args...)
		if err != nil {
			return fmt.Errorf("failed to delete records from %s: %v", table, err)
		}
		log.Printf("Dependent records deleted successfully from %s", table)
	}

	return nil
}

// ImportGPUsFromNvidiaSMI runs the `nvidia-smi` command, parses the output, and updates the GPUs table.
func ImportGPUsFromNvidiaSMI(db *sql.DB, mode string, verbose bool) error {
	// Get the server name from the hostname
	serverName, err := os.Hostname()
	if err != nil {
		return fmt.Errorf("failed to get hostname: %v", err)
	}

	// Handle the "remake" mode by deleting all dependent records and truncating the table
	if mode == "remake" {
		err := deleteDependentRecords(db, "WHERE gpu_uuid IN (SELECT gpu_uuid FROM gpus)")
		if err != nil {
			return err
		}

		// Truncate the gpus table
		_, err = db.Exec("DELETE FROM gpus")
		if err != nil {
			return fmt.Errorf("failed to truncate GPUs table: %v", err)
		}
		log.Println("GPUs table deleted successfully")
	}

	// Handle the "remake_for_server" mode by deleting rows for the current server
	if mode == "remake_for_server" {
		err := deleteDependentRecords(db, "WHERE gpu_uuid IN (SELECT gpu_uuid FROM gpus WHERE server_name = ?)", serverName)
		if err != nil {
			return err
		}

		// Delete rows from the gpus table for this server
		_, err = db.Exec("DELETE FROM gpus WHERE server_name = ?", serverName)
		if err != nil {
			return fmt.Errorf("failed to delete GPUs for server %s: %v", serverName, err)
		}
		log.Printf("GPUs for server %s deleted successfully", serverName)
	}

	// Run the `nvidia-smi` command with the updated query string
	cmd := exec.Command("nvidia-smi", "--query-gpu=gpu_uuid,index,name,memory.total,gpu_serial,gpu_bus_id", "--format=csv,noheader,nounits")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to run nvidia-smi: %v", err)
	}

	// Parse the CSV output
	reader := csv.NewReader(strings.NewReader(string(output)))
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to parse nvidia-smi output: %v", err)
	}

	// Insert or update each GPU record
	for _, record := range records {
		if len(record) != 6 {
			return fmt.Errorf("invalid record format: expected 6 fields, got %d", len(record))
		}

		// Extract GPU information and trim whitespace
		gpuUUID := strings.TrimSpace(record[0])
		gpuNumber := strings.TrimSpace(record[1])
		gpuName := strings.TrimSpace(record[2])
		vramSizeMB := strings.TrimSpace(record[3])
		gpuSerial := strings.TrimSpace(record[4])
		busID := strings.TrimSpace(record[5])

		// Prepare the SQL query based on the mode
		var query string
		switch mode {
		case "remake", "update", "remake_for_server":
			query = `
                INSERT INTO gpus (gpu_uuid, server_name, gpu_number, model_name, vram_size_mb, gpu_serial, gpu_bus_id, created_at, updated_at)
                VALUES (?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
                ON DUPLICATE KEY UPDATE
                    server_name = VALUES(server_name),
                    gpu_number = VALUES(gpu_number),
                    model_name = VALUES(model_name),
                    vram_size_mb = VALUES(vram_size_mb),
                    gpu_serial = VALUES(gpu_serial),
                    gpu_bus_id = VALUES(gpu_bus_id),
                    updated_at = NOW()`
		case "insert":
			query = `
                INSERT IGNORE INTO gpus (gpu_uuid, server_name, gpu_number, model_name, vram_size_mb, gpu_serial, gpu_bus_id, created_at, updated_at)
                VALUES (?, ?, ?, ?, ?, ?, ?, NOW(), NOW())`
		default:
			return fmt.Errorf("invalid mode: %s. Use 'remake', 'update', 'insert', or 'remake_for_server'", mode)
		}

		// Log verbose output if enabled
		if verbose {
			log.Printf("Inserting/Updating GPU Record: UUID=%s, Server=%s, Number=%s, Name=%s, VRAM=%s, Serial=%s, BusID=%s",
				gpuUUID, serverName, gpuNumber, gpuName, vramSizeMB, gpuSerial, busID)
			log.Printf("Executing Query: %s", query)
		}

		// Execute the query
		_, err := db.Exec(query, gpuUUID, serverName, gpuNumber, gpuName, vramSizeMB, gpuSerial, busID)
		if err != nil {
			return fmt.Errorf("failed to insert or update GPU record for UUID %s: %v", gpuUUID, err)
		}
	}

	log.Println("GPUs table updated successfully")
	return nil
}
