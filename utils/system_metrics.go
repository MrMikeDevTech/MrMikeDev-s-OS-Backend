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

func init() {
	if _, err := os.Stat("/host/proc"); err == nil {
		os.Setenv("HOST_PROC", "/host/proc")
		os.Setenv("HOST_SYS", "/host/sys")
		os.Setenv("HOST_ETC", "/host/etc")
	}
}

func Round(val float64) float64 {
	return math.Round(val*100) / 100
}

func GetWatts() float64 {
	data, err := os.ReadFile("/host/sys/class/powercap/intel-rapl:0/energy_uj")

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
	pct, _ := cpu.Percent(0, false)
	logical, _ := cpu.Counts(true)
	physical, _ := cpu.Counts(false)

	temps, _ := host.SensorsTemperatures()
	var temp float64
	if len(temps) > 0 {
		temp = temps[0].Temperature
	}

	watts := GetWatts()

	return map[string]interface{}{
		"percentage": Round(pct[0]),
		"cores":      physical,
		"threads":    logical,
		"temp_c":     Round(temp),
		"watts":      Round(watts),
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

func GetDisk() map[string]interface{} {
	d, _ := disk.Usage("/host")
	return map[string]interface{}{
		"total_mb":   d.Total / 1024 / 1024,
		"used_mb":    d.Used / 1024 / 1024,
		"free_mb":    d.Free / 1024 / 1024,
		"percentage": Round(d.UsedPercent),
	}
}

func GetNetwork(interval time.Duration) map[string]interface{} {
	n1, _ := net.IOCounters(false)
	time.Sleep(interval)
	n2, _ := net.IOCounters(false)

	down := float64(n2[0].BytesRecv-n1[0].BytesRecv) / 1024 / 1024
	up := float64(n2[0].BytesSent-n1[0].BytesSent) / 1024 / 1024

	return map[string]interface{}{
		"download_mbps": Round(down),
		"upload_mbps":   Round(up),
	}
}
