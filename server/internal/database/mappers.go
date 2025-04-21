package database

import "database/sql"

func mapRealTimeUsage(rows *sql.Rows) (RealTimeUsage, error) {
	var usage RealTimeUsage
	err := rows.Scan(
		&usage.ServerName,
		&usage.GPUNumber,
		&usage.Utilization,
		&usage.MemoryUtilization,
		&usage.MemoryUsedMB,
		&usage.MemoryAvailableMB,
		&usage.PowerUsageWatts,
		&usage.TemperatureCelsius,
		&usage.UpdatedAt,
	)
	return usage, err
}
