package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gpu_request/daemon/commands"
	"gpu_request/daemon/monitor"
)

func main() {
	// Load configuration (e.g., DSN, interval) from environment variables or flags
	dsn := "user:password@tcp(localhost:3306)/gpu_scheduler"
	sleepInterval := 10 * time.Second

	// Create channels for graceful shutdown
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	// Start GPU monitoring
	go func() {
		err := monitor.StartGPUMonitor(dsn, sleepInterval)
		if err != nil {
			log.Fatalf("GPU Monitor failed: %v", err)
		}
	}()

	// Start command monitoring
	go func() {
		err := commands.StartCommandMonitor(dsn, sleepInterval)
		if err != nil {
			log.Fatalf("Command Monitor failed: %v", err)
		}
	}()

	// Wait for shutdown signal
	<-stopChan
	log.Println("Shutting down daemon...")
}
