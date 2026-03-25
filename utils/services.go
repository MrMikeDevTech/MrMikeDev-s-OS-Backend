package utils

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

type ServiceStatus struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

func GetSystemServices() []ServiceStatus {
	monitoredServices := map[string]string{
		"mrmikedevs-os.service": "OS Backend",
		"nginx.service":         "Nginx Proxy",
		"cloudflared.service":   "Cloudflare Tunnel",
		"ssh.service":           "SSH Server",
		"ufw.service":           "UFW Firewall",
		"docker.service":        "Docker Engine",
		"docker.socket":         "Docker Socket",
		"tailscaled.service":    "Tailscale VPN",
	}

	var results []ServiceStatus

	isWindows := runtime.GOOS == "windows"

	for id, friendlyName := range monitoredServices {
		var status string

		if isWindows {
			status = "not-available"
		} else {
			cmd := exec.Command("systemctl", "is-active", id)
			out, _ := cmd.Output()
			status = strings.TrimSpace(string(out))
		}

		results = append(results, ServiceStatus{
			ID:     id,
			Name:   friendlyName,
			Status: status,
		})
	}
	return results
}

func HandleServiceAction(serviceID string, action string) error {
	validActions := map[string]bool{"start": true, "stop": true, "restart": true}
	if !validActions[action] {
		return fmt.Errorf("acción no permitida: %s", action)
	}

	if runtime.GOOS == "windows" {
		fmt.Printf("Mock: %s %s\n", action, serviceID)
		return nil
	}

	cmd := exec.Command("sudo", "systemctl", action, serviceID)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error: %s - %v", string(out), err)
	}

	return nil
}

type LogEntry struct {
	Timestamp *time.Time `json:"timestamp"`
	Message   string     `json:"message"`
}

func GetServiceLogs(serviceID string) ([]LogEntry, error) {
	if runtime.GOOS == "windows" {
		return []LogEntry{{Message: "Logs no disponibles en Windows"}}, nil
	}

	cmd := exec.Command("sudo", "journalctl", "-u", serviceID, "-n", "250", "--no-pager", "--output=short-iso")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("error al obtener logs: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	var logEntries []LogEntry

	for _, line := range lines {
		line = strings.TrimRight(line, " \r")
		if line == "" {
			continue
		}

		parts := strings.SplitN(line, " ", 3)

		var entry LogEntry

		t, err := time.Parse(time.RFC3339, parts[0])
		if err == nil {
			entry.Timestamp = &t

			if len(parts) >= 3 {
				msgPart := parts[2]
				if idx := strings.Index(msgPart, ": "); idx != -1 {
					entry.Message = msgPart[idx+2:]
				} else {
					entry.Message = msgPart
				}
			} else {
				entry.Message = ""
			}
		} else {
			entry.Timestamp = nil
			entry.Message = line
		}

		logEntries = append(logEntries, entry)
	}

	return logEntries, nil
}
