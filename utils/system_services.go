package utils

import (
	"os/exec"
	"runtime"
	"strings"
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
		"tailscale.service":     "Tailscale VPN",
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
