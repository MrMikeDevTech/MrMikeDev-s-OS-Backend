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

func Round(val float64) float64 {
	return math.Round(val*100) / 100
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

	var watts float64
	data, err := os.ReadFile("/sys/class/powercap/intel-rapl:0/energy_uj")
	if err == nil {
		val, _ := strconv.ParseFloat(strings.TrimSpace(string(data)), 64)
		watts = val / 1e6
	}

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
	d, _ := disk.Usage("/")
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
