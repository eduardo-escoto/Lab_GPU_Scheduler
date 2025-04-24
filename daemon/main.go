package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/eduardo-escoto/gpu_request/daemon/monitor"
)

func main() {
	// Define command-line flags
	dsnFlag := flag.String("dsn", "", "Database DSN (e.g., user:password@tcp(localhost:3306)/gpu_scheduler)")
	intervalFlag := flag.String("interval", "", "Sleep interval in seconds between updates")
	verboseFlag := flag.Bool("verbose", false, "Enable verbose logging") // Add verbose flag

	// Parse command-line flags
	flag.Parse()

	// Load DSN from environment variable or command-line flag, fallback to default
	dsn := os.Getenv("DATABASE_DSN")
	if *dsnFlag != "" {
		dsn = *dsnFlag
	}
	if dsn == "" {
		dsn = "user:password@tcp(localhost:3306)/gpu_scheduler" // Default DSN
	}

	// Debug: Print the DSN
	if *verboseFlag {
		log.Printf("Using DSN: %s", dsn)
	}

	// Load sleep interval from environment variable or command-line flag, fallback to default
	sleepIntervalStr := os.Getenv("INTERVAL")
	if *intervalFlag != "" {
		sleepIntervalStr = *intervalFlag
	}
	sleepInterval := 10 * time.Second // Default interval
	if sleepIntervalStr != "" {
		interval, err := strconv.Atoi(sleepIntervalStr)
		if err != nil {
			log.Fatalf("Invalid interval value: %v", err)
		}
		sleepInterval = time.Duration(interval) * time.Second
	}

	// Debug: Print the sleep interval
	if *verboseFlag {
		log.Printf("Using sleep interval: %s", sleepInterval)
	}

	// Create channels for graceful shutdown
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	// Start GPU monitoring
	go func() {
		err := monitor.StartGPUMonitor(dsn, sleepInterval, *verboseFlag) // Pass verbose flag
		if err != nil {
			log.Fatalf("GPU Monitor failed: %v", err)
		}
	}()

	// Wait for shutdown signal
	<-stopChan
	log.Println("Shutting down daemon...")
}
