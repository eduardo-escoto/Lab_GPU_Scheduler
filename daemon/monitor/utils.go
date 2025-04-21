package monitor

import (
	"fmt"
	"os/exec"
	"strings"
)

// GetServerName retrieves the hostname of the server
func GetServerName() (string, error) {
	cmd := exec.Command("hostname")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get server name: %v", err)
	}
	return strings.TrimSpace(string(output)), nil
}
