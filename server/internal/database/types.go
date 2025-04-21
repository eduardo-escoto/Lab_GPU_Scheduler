package database

import "time"

type RealTimeUsage struct {
	ServerName         string
	GPUNumber          int
	Utilization        float32
	MemoryUtilization  float32
	MemoryUsedMB       int
	MemoryAvailableMB  int
	PowerUsageWatts    float32
	TemperatureCelsius float32
	UpdatedAt          time.Time
}
