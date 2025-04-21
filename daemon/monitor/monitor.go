package monitor

import (
	"database/sql"
	"log"
	"time"

	"github.com/eduardo-escoto/gpu_request/daemon/monitor/gpu_metrics"

	_ "github.com/go-sql-driver/mysql"
)

func StartGPUMonitor(dsn string, interval time.Duration) error {
	// Connect to the database
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	defer db.Close()

	// Get the server name
	serverName, err := gpu_metrics.GetServerName()
	if err != nil {
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

func updateDatabase(db *sql.DB, serverName string, gpuUsages []gpu_metrics.GPU) error {
	for _, gpu := range gpuUsages {
		_, err := db.Exec(`
			INSERT INTO real_time_usage (server_name, gpu_number, utilization, memory_utilization, memory_used_mb, memory_available_mb, power_usage_watts, temperature_celsius, updated_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, NOW())
			ON DUPLICATE KEY UPDATE
				utilization = VALUES(utilization),
				memory_utilization = VALUES(memory_utilization),
				memory_used_mb = VALUES(memory_used_mb),
				memory_available_mb = VALUES(memory_available_mb),
				power_usage_watts = VALUES(power_usage_watts),
				temperature_celsius = VALUES(temperature_celsius),
				updated_at = NOW()`,
			serverName, gpu.Index, gpu.UtilizationGPU, gpu.UtilizationMemory, gpu.MemoryUsedMB, gpu.MemoryFreeMB, gpu.PowerDrawWatts, gpu.TemperatureCelsius,
		)
		if err != nil {
			return err
		}
	}
	return nil
}
