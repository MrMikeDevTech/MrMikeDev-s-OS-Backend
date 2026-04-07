package utils

import (
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
)

var lastEnergy float64
var lastTime time.Time

func Round(val float64) float64 {
	return math.Round(val*100) / 100
}

func GetWatts() float64 {
	data, err := os.ReadFile("/sys/class/powercap/intel-rapl:0/energy_uj")
	if err != nil {
		return 0
	}

	currentEnergy, _ := strconv.ParseFloat(strings.TrimSpace(string(data)), 64)
	currentTime := time.Now()

	if lastTime.IsZero() {
		lastEnergy = currentEnergy
		lastTime = currentTime
		return 0
	}

	diffJoules := (currentEnergy - lastEnergy) / 1_000_000
	diffSeconds := currentTime.Sub(lastTime).Seconds()

	lastEnergy = currentEnergy
	lastTime = currentTime

	if diffSeconds > 0 {
		return Round(diffJoules / diffSeconds)
	}

	return 0
}

func GetCPU() map[string]interface{} {
	model, _ := cpu.Info()
	pct, _ := cpu.Percent(0, false)
	logical, _ := cpu.Counts(true)
	physical, _ := cpu.Counts(false)

	temps, _ := host.SensorsTemperatures()
	var temp float64
	for _, t := range temps {
		if strings.Contains(strings.ToLower(t.SensorKey), "package") || strings.Contains(strings.ToLower(t.SensorKey), "core") {
			temp = t.Temperature
			break
		}
	}

	if temp == 0 && len(temps) > 0 {
		temp = temps[0].Temperature
	}

	return map[string]interface{}{
		"model":      model[0].ModelName,
		"percentage": Round(pct[0]),
		"cores":      physical,
		"threads":    logical,
		"temp_c":     Round(temp),
		"watts":      Round(GetWatts()),
	}
}

func GetRAM() map[string]interface{} {
	v, _ := mem.VirtualMemory()
	return map[string]interface{}{
		"total_mb":   v.Total / 1024 / 1024,
		"used_mb":    v.Used / 1024 / 1024,
		"percentage": Round(v.UsedPercent),
	}
}

func GetDisks() []map[string]interface{} {
	results := make([]map[string]interface{}, 0)

	potentialMounts := []string{"/", "/mnt/SSD"}

	for _, path := range potentialMounts {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			continue
		}

		d, err := disk.Usage(path)
		if err != nil {
			continue
		}

		results = append(results, map[string]interface{}{
			"mount_point": path,
			"total_gb":    Round(float64(d.Total) / 1024 / 1024 / 1024),
			"used_gb":     Round(float64(d.Used) / 1024 / 1024 / 1024),
			"percentage":  Round(d.UsedPercent),
		})
	}

	return results
}

func GetNetwork(interval time.Duration) map[string]interface{} {
	n1, _ := net.IOCounters(false)
	time.Sleep(interval)
	n2, _ := net.IOCounters(false)

	if len(n1) == 0 || len(n2) == 0 {
		return map[string]interface{}{"download_mbps": 0, "upload_mbps": 0}
	}

	down := float64(n2[0].BytesRecv-n1[0].BytesRecv) / 1024 / 1024
	up := float64(n2[0].BytesSent-n1[0].BytesSent) / 1024 / 1024

	return map[string]interface{}{
		"download_mbps": Round(down),
		"upload_mbps":   Round(up),
	}
}
