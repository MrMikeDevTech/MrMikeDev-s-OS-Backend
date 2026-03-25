package utils

import (
	"os"
	"os/exec"
	"runtime"
)

func getNginxPath() string {
	if runtime.GOOS == "windows" {
		return "./test_local/nginx_mock.conf"
	}
	return "/etc/nginx/nginx.conf"
}

func ReadNginxConfig() (string, error) {
	content, err := os.ReadFile(getNginxPath())
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func TestNginxConfig(content string) (string, error) {
	tempFile := "/tmp/nginx_temp.conf"
	if runtime.GOOS == "windows" {
		tempFile = "./test_local/nginx_temp.conf"
	}

	err := os.WriteFile(tempFile, []byte(content), 0644)
	if err != nil {
		return "", err
	}

	if runtime.GOOS == "windows" {
		return "Syntax OK", nil
	}

	cmd := exec.Command("sudo", "nginx", "-t", "-c", tempFile)
	out, err := cmd.CombinedOutput()

	return string(out), err
}

func SaveNginxConfig(content string) error {
	path := getNginxPath()

	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		return err
	}

	if runtime.GOOS != "windows" {
		exec.Command("sudo", "systemctl", "reload", "nginx").Run()
	}
	return nil
}
