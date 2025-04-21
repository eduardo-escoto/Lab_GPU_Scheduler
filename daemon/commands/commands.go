package commands

import (
	"database/sql"
	"log"
	"os/exec"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Command struct {
	ID          int
	TargetNode  string
	CommandType string
	Parameters  string
	Status      string
}

func StartCommandMonitor(dsn string, interval time.Duration) error {
	// Connect to the database
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	defer db.Close()

	// Get the server name
	serverName, err := exec.Command("hostname").Output()
	if err != nil {
		return err
	}
	nodeName := strings.TrimSpace(string(serverName))

	// Daemon loop
	for {
		log.Println("Checking for new commands...")
		commands, err := fetchPendingCommands(db, nodeName)
		if err != nil {
			log.Printf("Error fetching commands: %v", err)
			time.Sleep(interval)
			continue
		}

		for _, cmd := range commands {
			log.Printf("Executing command: %v", cmd)
			err := executeCommand(db, cmd)
			if err != nil {
				log.Printf("Error executing command %d: %v", cmd.ID, err)
			}
		}

		time.Sleep(interval)
	}
}

func fetchPendingCommands(db *sql.DB, nodeName string) ([]Command, error) {
	rows, err := db.Query("SELECT id, command_type, parameters FROM commands WHERE target_node = ? AND status = 'pending'", nodeName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var commands []Command
	for rows.Next() {
		var cmd Command
		err := rows.Scan(&cmd.ID, &cmd.CommandType, &cmd.Parameters)
		if err != nil {
			return nil, err
		}
		commands = append(commands, cmd)
	}
	return commands, nil
}

func executeCommand(db *sql.DB, cmd Command) error {
	// Mark command as in progress
	_, err := db.Exec("UPDATE commands SET status = 'in_progress' WHERE id = ?", cmd.ID)
	if err != nil {
		return err
	}

	// Simulate command execution (replace with actual logic)
	log.Printf("Executing command: %s with parameters: %s", cmd.CommandType, cmd.Parameters)
	time.Sleep(2 * time.Second) // Simulate work

	// Mark command as completed
	_, err = db.Exec("UPDATE commands SET status = 'completed' WHERE id = ?", cmd.ID)
	return err
}
