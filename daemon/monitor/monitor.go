package monitor

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// StartGPUMonitor starts the GPU monitoring daemon
func StartGPUMonitor(dsn string, interval time.Duration, verbose bool) error {
	// Debug: Print the DSN being used
	if verbose {
		log.Printf("Connecting to database with DSN: %s", dsn)
	}

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
		if verbose {
			log.Println("Fetching GPU usage data...")
		}

		gpuUsages, err := GetGPUMetrics(verbose) // Pass verbose flag to GetGPUMetrics
		if err != nil {
			log.Printf("Error fetching GPU metrics: %v", err)
			time.Sleep(interval)
			continue
		}

		if verbose {
			log.Println("Updating database with GPU usage data...")
		}

		err = updateDatabase(db, serverName, gpuUsages, verbose) // Pass verbose flag to updateDatabase
		if err != nil {
			log.Printf("Error updating database: %v", err)
		}

		time.Sleep(interval)
	}
}

// updateDatabase inserts new records into the real_time_usage and gpu_processes tables
func updateDatabase(db *sql.DB, serverName string, gpuUsages []GPU, verbose bool) error {
	// Get the current timestamp
	timestamp := time.Now()

	for _, gpu := range gpuUsages {
		// Insert a new record into the real_time_usage table
		_, err := db.Exec(`
			INSERT INTO gpu_scheduler.real_time_usage (gpu_uuid, gpu_name, server_name, gpu_number, utilization, memory_utilization, memory_used_mb, 
				memory_available_mb, power_usage_watts, temperature_celsius, reported_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			gpu.UUID, gpu.Name, serverName, gpu.Index, gpu.UtilizationGPU, gpu.UtilizationMemory, gpu.MemoryUsedMB,
			gpu.MemoryFreeMB, gpu.PowerDrawWatts, gpu.TemperatureCelsius, timestamp,
		)
		if err != nil {
			return fmt.Errorf("failed to insert real_time_usage for GPU %d: %v", gpu.Index, err)
		}

		if verbose {
			log.Printf("Inserted GPU %d (%s) into real_time_usage table: Utilization=%.2f%%, Memory Used=%d MB, Power=%.2f W, Temp=%.2f°C, Timestamp=%s",
				gpu.Index, gpu.UUID, gpu.UtilizationGPU, gpu.MemoryUsedMB, gpu.PowerDrawWatts, gpu.TemperatureCelsius, timestamp)
		}

		// Insert records into the gpu_processes table for each process running on the GPU
		for _, process := range gpu.Processes {
			// Calculate GPU utilization as a percentage and round to two decimal places
			gpuUtilizationPercentage := float64(process.UsedGPUMemoryMB) / float64(gpu.MemoryTotalMB) * 100
			gpuUtilizationPercentageRounded := fmt.Sprintf("%.2f", gpuUtilizationPercentage)

			_, err := db.Exec(`
				INSERT INTO gpu_scheduler.gpu_processes (gpu_uuid, process_id, process_name, user_name, gpu_utilization, used_gpu_memory, reported_at)
				VALUES (?, ?, ?, ?, ?, ?, ?)`,
				gpu.UUID, process.PID, process.ProcessName, process.UserName, gpuUtilizationPercentageRounded, process.UsedGPUMemoryMB, timestamp,
			)
			if err != nil {
				return fmt.Errorf("failed to insert gpu_processes for GPU %d, PID %d: %v", gpu.Index, process.PID, err)
			}

			if verbose {
				log.Printf("Inserted Process %d (%s) by user %s into gpu_processes table: GPU %d (%s), GPU Utilization=%s%%, Memory Used=%d MB, Timestamp=%s",
					process.PID, process.ProcessName, process.UserName, gpu.Index, gpu.UUID, gpuUtilizationPercentageRounded, process.UsedGPUMemoryMB, timestamp)
			}
		}
	}

	return nil
}
