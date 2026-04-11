package utils

import (
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
)

var (
	lastEnergy float64
	lastTime   time.Time
	lastWatts  float64
	mu         sync.Mutex
)

func Round(val float64) float64 {
	return math.Round(val*100) / 100
}

func getEnergyUJ() (float64, error) {
	paths := []string{
		"/sys/class/powercap/intel-rapl:0/energy_uj",
		"/sys/devices/virtual/powercap/intel-rapl/intel-rapl:0/energy_uj",
	}

	for _, p := range paths {
		data, err := os.ReadFile(p)
		if err != nil {
			continue
		}
		val, err := strconv.ParseFloat(strings.TrimSpace(string(data)), 64)
		if err == nil && val > 0 {
			return val, nil
		}
	}

	matches, _ := filepath.Glob("/sys/class/powercap/intel-rapl*/energy_uj")
	for _, p := range matches {
		data, err := os.ReadFile(p)
		if err == nil {
			val, err := strconv.ParseFloat(strings.TrimSpace(string(data)), 64)
			if err == nil && val > 0 {
				return val, nil
			}
		}
	}
	return 0, os.ErrNotExist
}

func GetWatts() float64 {
	mu.Lock()
	defer mu.Unlock()

	currentEnergy, err := getEnergyUJ()
	if err != nil {
		return lastWatts
	}

	currentTime := time.Now()

	if lastTime.IsZero() {
		lastEnergy = currentEnergy
		lastTime = currentTime
		return 0
	}

	diffJoules := (currentEnergy - lastEnergy) / 1_000_000.0
	diffSeconds := currentTime.Sub(lastTime).Seconds()

	if diffSeconds < 0.1 {
		return lastWatts
	}

	if diffJoules < 0 {
		lastEnergy = currentEnergy
		lastTime = currentTime
		return lastWatts
	}

	lastEnergy = currentEnergy
	lastTime = currentTime

	wattage := diffJoules / diffSeconds

	if wattage > 150.0 || wattage < 0 {
		return lastWatts
	}

	lastWatts = Round(wattage)
	return lastWatts
}

func GetCPU() map[string]interface{} {
	model, _ := cpu.Info()
	pct, _ := cpu.Percent(0, false)
	logical, _ := cpu.Counts(true)
	physical, _ := cpu.Counts(false)

	temps, _ := host.SensorsTemperatures()
	var temp float64
	for _, t := range temps {
		sensorName := strings.ToLower(t.SensorKey)
		if strings.Contains(sensorName, "package") || strings.Contains(sensorName, "core 0") {
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
		"watts":      GetWatts(),
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
