package monitor

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func StartGPUMonitor(dsn string, interval time.Duration) error {
	// Debug: Print the DSN being used
	log.Printf("Connecting to database with DSN: %s", dsn)

	// Connect to the database
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Printf("Error connecting to database: %v", err)
		return err
	}
	defer db.Close()

	// Test the connection
	err = db.Ping()
	if err != nil {
		log.Printf("Error pinging database: %v", err)
		return err
	}

	// Get the server name
	serverName, err := GetServerName()
	if err != nil {
		log.Printf("Error getting server name: %v", err)
		return err
	}

	// Daemon loop
	for {
		log.Println("Fetching GPU usage data...")
		gpuUsages, err := GetGPUMetrics()
		if err != nil {
			log.Printf("Error fetching GPU metrics: %v", err)
			time.Sleep(interval)
			continue
		}

		log.Println("Updating database with GPU usage data...")
		err = updateDatabase(db, serverName, gpuUsages)
		if err != nil {
			log.Printf("Error updating database: %v", err)
		}

		time.Sleep(interval)
	}
}

func updateDatabase(db *sql.DB, serverName string, gpuUsages []GPU) error {
	for _, gpu := range gpuUsages {
		// Check if the GPU entry already exists
		var exists bool
		err := db.QueryRow(`
			SELECT EXISTS(
				SELECT 1 FROM real_time_usage WHERE server_name = ? AND gpu_number = ?
			)`, serverName, gpu.Index).Scan(&exists)
		if err != nil {
			return err
		}

		if exists {
			// Update the existing record
			_, err = db.Exec(`
				UPDATE gpu_scheduler.real_time_usage
				SET utilization = ?, memory_utilization = ?, memory_used_mb = ?, memory_available_mb = ?, 
					power_usage_watts = ?, temperature_celsius = ?, updated_at = NOW()
				WHERE server_name = ? AND gpu_number = ?`,
				gpu.UtilizationGPU, gpu.UtilizationMemory, gpu.MemoryUsedMB, gpu.MemoryFreeMB,
				gpu.PowerDrawWatts, gpu.TemperatureCelsius, serverName, gpu.Index,
			)
			if err != nil {
				return err
			}
		} else {
			// Insert a new record
			_, err = db.Exec(`
				INSERT INTO gpu_scheduler.real_time_usage (server_name, gpu_number, utilization, memory_utilization, memory_used_mb, 
					memory_available_mb, power_usage_watts, temperature_celsius, updated_at)
				VALUES (?, ?, ?, ?, ?, ?, ?, ?, NOW())`,
				serverName, gpu.Index, gpu.UtilizationGPU, gpu.UtilizationMemory, gpu.MemoryUsedMB,
				gpu.MemoryFreeMB, gpu.PowerDrawWatts, gpu.TemperatureCelsius,
			)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
