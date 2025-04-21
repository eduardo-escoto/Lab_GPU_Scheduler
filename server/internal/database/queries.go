package database

import (
	"database/sql"
	"log"
)

func QueryRealTimeUsage(db *sql.DB) ([]RealTimeUsage, error) {
	query := `
        SELECT server_name, gpu_number, utilization, memory_utilization, memory_used_mb,
               memory_available_mb, power_usage_watts, temperature_celsius, updated_at
        FROM gpu_scheduler.real_time_usage;
    `

	usages, err := QueryAndMap(db, query, nil, mapRealTimeUsage)
	if err != nil {
		log.Printf("Error querying real-time usage: %v", err)
		return nil, err
	}

	return usages, nil
}
