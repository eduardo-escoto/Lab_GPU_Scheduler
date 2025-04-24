package monitor

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

// GPU represents the metrics for a single GPU
type GPU struct {
	Index              int          // GPU index
	Name               string       // GPU name
	UUID               string       // GPU UUID
	MemoryTotalMB      int          // Total memory in MB
	MemoryUsedMB       int          // Used memory in MB
	MemoryFreeMB       int          // Free memory in MB
	PowerDrawWatts     float64      // Power draw in watts
	PowerLimitWatts    float64      // Power limit in watts
	TemperatureCelsius float64      // Temperature in Celsius
	UtilizationGPU     float64      // GPU utilization percentage
	UtilizationMemory  float64      // Memory utilization percentage
	Processes          []GPUProcess // List of processes running on the GPU
}

// GPUProcess represents a single process running on a GPU
type GPUProcess struct {
	PID             int    // Process ID
	ProcessName     string // Name of the process
	UserName        string // User running the process
	UsedGPUMemoryMB int    // GPU memory used by the process in MB
}

// GetGPUMetrics fetches GPU metrics using nvidia-smi
func GetGPUMetrics(verbose bool) ([]GPU, error) {
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

		// Fetch processes running on this GPU
		processes, err := GetGPUProcesses(index, verbose)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch processes for GPU %d: %v", index, err)
		}

		// Log GPU metrics if verbose mode is enabled
		if verbose {
			log.Printf("GPU %d (%s): UUID=%s, Total Memory=%d MB, Used Memory=%d MB, Free Memory=%d MB, Power Draw=%.2f W, Power Limit=%.2f W, Temperature=%.2f°C, GPU Utilization=%.2f%%, Memory Utilization=%.2f%%",
				index, name, uuid, memoryTotalMB, memoryUsedMB, memoryFreeMB, powerDrawWatts, powerLimitWatts, temperatureCelsius, utilizationGPU, utilizationMemory)
			for _, process := range processes {
				log.Printf("  Process %d (%s) by user %s: %d MB used",
					process.PID, process.ProcessName, process.UserName, process.UsedGPUMemoryMB)
			}
		}

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
			Processes:          processes,
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to parse nvidia-smi output: %v", err)
	}

	return gpus, nil
}

// GetGPUProcesses fetches the processes running on a specific GPU
func GetGPUProcesses(gpuIndex int, verbose bool) ([]GPUProcess, error) {
	// Command to query GPU processes
	cmd := exec.Command("nvidia-smi",
		"--query-compute-apps=pid,process_name,used_gpu_memory",
		"--format=csv,noheader,nounits",
		"--id="+strconv.Itoa(gpuIndex))

	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("failed to execute nvidia-smi for processes: %v", err)
	}

	var processes []GPUProcess
	scanner := bufio.NewScanner(&out)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Split(line, ",")
		if len(fields) != 3 {
			return nil, fmt.Errorf("unexpected nvidia-smi process output format")
		}

		// Parse the fields
		pid, _ := strconv.Atoi(strings.TrimSpace(fields[0]))
		processName := strings.TrimSpace(fields[1])
		usedGPUMemoryMB, _ := strconv.Atoi(strings.TrimSpace(fields[2]))

		// Fetch the username for the process
		userName, err := getProcessUser(pid)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch username for PID %d: %v", pid, err)
		}

		// Log process details if verbose mode is enabled
		if verbose {
			log.Printf("Process %d (%s) by user %s: %d MB used",
				pid, processName, userName, usedGPUMemoryMB)
		}

		// Append process information to the list
		processes = append(processes, GPUProcess{
			PID:             pid,
			ProcessName:     processName,
			UserName:        userName,
			UsedGPUMemoryMB: usedGPUMemoryMB,
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to parse nvidia-smi process output: %v", err)
	}

	return processes, nil
}

// getProcessUser fetches the username of the user running a specific process
func getProcessUser(pid int) (string, error) {
	// Command to query the username for a specific PID
	cmd := exec.Command("ps", "-o", "user=", "-p", strconv.Itoa(pid))

	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("failed to execute ps command for PID %d: %v", pid, err)
	}

	// Trim the output to get the username
	username := strings.TrimSpace(out.String())
	if username == "" {
		return "", fmt.Errorf("failed to fetch username for PID %d", pid)
	}

	return username, nil
}
