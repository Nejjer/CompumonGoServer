package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// ---- CPU ----

func getCPUInfo() CPUInfo {
	load := getCPULoad()
	temp := getCPUTemperature()
	fan := getCPUFanSpeed()

	return CPUInfo{Load: load, Temperature: temp, FanSpeed: fan}
}

func getCPULoad() float64 {
	// читаем /proc/loadavg
	data, err := os.ReadFile("/proc/loadavg")
	if err != nil {
		return 0
	}
	fields := strings.Fields(string(data))
	if len(fields) == 0 {
		return 0
	}
	load, _ := strconv.ParseFloat(fields[0], 64)
	return load
}

func getCPUTemperature() float64 {
	out, err := exec.Command("sensors").Output()
	if err != nil {
		return 0
	}

	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		if strings.Contains(line, "Tctl:") {
			parts := strings.Fields(line)
			if len(parts) < 2 {
				continue
			}
			tempStr := strings.TrimPrefix(parts[1], "+")
			tempStr = strings.TrimSuffix(tempStr, "°C")
			temp, err := strconv.ParseFloat(tempStr, 64)
			if err != nil {
				continue
			}
			return temp
		}
	}

	return 0
}

func getCPUFanSpeed() float64 {
	files, err := os.ReadDir("/sys/class/hwmon/")
	if err != nil {
		return 0
	}
	for _, f := range files {
		path := fmt.Sprintf("/sys/class/hwmon/%s/fan1_input", f.Name())
		data, err := os.ReadFile(path)
		if err == nil {
			fan, _ := strconv.ParseFloat(strings.TrimSpace(string(data)), 64)
			return fan
		}
	}
	return 0
}

// ---- RAM ----

func getRAMInfo() RAMInfo {
	data, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		return RAMInfo{}
	}
	lines := strings.Split(string(data), "\n")
	var total, free float64
	for _, l := range lines {
		if strings.HasPrefix(l, "MemTotal:") {
			fields := strings.Fields(l)
			total, _ = strconv.ParseFloat(fields[1], 64)
		}
		if strings.HasPrefix(l, "MemAvailable:") {
			fields := strings.Fields(l)
			free, _ = strconv.ParseFloat(fields[1], 64)
		}
	}
	used := total - free
	// /proc/meminfo в kB
	return RAMInfo{
		Total: floor1(total / (1024 * 1000)),
		Used:  floor1(used / (1024 * 1000)),
	}
}

// ---- GPU ----

func getGPUInfo() *GPUInfo {
	out, err := exec.Command("nvidia-smi", "--query-gpu=utilization.gpu,temperature.gpu,fan.speed,memory.total,memory.used", "--format=csv,noheader,nounits").Output()
	if err != nil {
		return nil
	}
	fields := strings.Split(strings.TrimSpace(string(out)), ", ")
	if len(fields) < 5 {
		return nil
	}

	load, _ := strconv.ParseFloat(fields[0], 64)
	temp, _ := strconv.ParseFloat(fields[1], 64)
	fan, _ := strconv.ParseFloat(fields[2], 64)
	memTotal, _ := strconv.ParseFloat(fields[3], 64)
	memUsed, _ := strconv.ParseFloat(fields[4], 64)

	return &GPUInfo{
		Load:        floor1(load),
		Temperature: temp,
		FanSpeed:    fan,
		Memory: GPUMemory{
			Total: floor1(memTotal / 1024),
			Used:  floor1(memUsed / 1024),
		},
	}
}

func getLinuxPCInfo() IPCInfo {
	info := IPCInfo{
		CPU: getCPUInfo(),
		RAM: getRAMInfo(),
	}

	// if gpu := getGPUInfo(); gpu != nil {
	// 	info.GPU = gpu
	// }

	return info
}
