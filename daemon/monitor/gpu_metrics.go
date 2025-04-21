package monitor

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// GPU represents the metrics for a single GPU
type GPU struct {
	Index              int     // GPU index
	Name               string  // GPU name
	UUID               string  // GPU UUID
	MemoryTotalMB      int     // Total memory in MB
	MemoryUsedMB       int     // Used memory in MB
	MemoryFreeMB       int     // Free memory in MB
	PowerDrawWatts     float64 // Power draw in watts
	PowerLimitWatts    float64 // Power limit in watts
	TemperatureCelsius float64 // Temperature in Celsius
	UtilizationGPU     float64 // GPU utilization percentage
	UtilizationMemory  float64 // Memory utilization percentage
}

// GetGPUMetrics fetches GPU metrics using nvidia-smi
func GetGPUMetrics() ([]GPU, error) {
	// Command to query GPU metrics
	cmd := exec.Command("nvidia-smi",
		"--query-gpu=index,name,uuid,memory.total,memory.used,memory.free,power.draw,power.limit,temperature.gpu,utilization.gpu,utilization.memory",
		"--format=csv,noheader,nounits")

	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("failed to execute nvidia-smi: %v", err)
	}

	var gpus []GPU
	scanner := bufio.NewScanner(&out)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Split(line, ",")
		if len(fields) != 11 {
			return nil, fmt.Errorf("unexpected nvidia-smi output format")
		}

		// Parse the fields
		index, _ := strconv.Atoi(strings.TrimSpace(fields[0]))
		name := strings.TrimSpace(fields[1])
		uuid := strings.TrimSpace(fields[2])
		memoryTotalMB, _ := strconv.Atoi(strings.TrimSpace(fields[3]))
		memoryUsedMB, _ := strconv.Atoi(strings.TrimSpace(fields[4]))
		memoryFreeMB, _ := strconv.Atoi(strings.TrimSpace(fields[5]))
		powerDrawWatts, _ := strconv.ParseFloat(strings.TrimSpace(fields[6]), 64)
		powerLimitWatts, _ := strconv.ParseFloat(strings.TrimSpace(fields[7]), 64)
		temperatureCelsius, _ := strconv.ParseFloat(strings.TrimSpace(fields[8]), 64)
		utilizationGPU, _ := strconv.ParseFloat(strings.TrimSpace(fields[9]), 64)
		utilizationMemory, _ := strconv.ParseFloat(strings.TrimSpace(fields[10]), 64)

		// Append GPU metrics to the list
		gpus = append(gpus, GPU{
			Index:              index,
			Name:               name,
			UUID:               uuid,
			MemoryTotalMB:      memoryTotalMB,
			MemoryUsedMB:       memoryUsedMB,
			MemoryFreeMB:       memoryFreeMB,
			PowerDrawWatts:     powerDrawWatts,
			PowerLimitWatts:    powerLimitWatts,
			TemperatureCelsius: temperatureCelsius,
			UtilizationGPU:     utilizationGPU,
			UtilizationMemory:  utilizationMemory,
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to parse nvidia-smi output: %v", err)
	}

	return gpus, nil
}
