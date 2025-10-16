package main

import (
	"encoding/json"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"
	"strings"
)

type Config struct {
	URL            string `json:"url"`
	CpuLoad        int    `json:"cpuLoad"`
	CpuTemp        int    `json:"cpuTemp"`
	FanSpeed       int    `json:"fanSpeed"`
	RAMUsed        int    `json:"ramUsed"`
	RAMAvailable   int    `json:"RAMAvailable"`
	GpuLoad        int    `json:"gpuLoad"`
	GpuTemp        int    `json:"gpuTemp"`
	GpuFanSpeed    int    `json:"gpuFanSpeed"`
	GpuMemoryTotal int    `json:"gpuMemoryTotal"`
	GpuMemoryUsed  int    `json:"gpuMemoryUsed"`
	Port           int    `json:"port"`
}

type Sensor struct {
	ID       int      `json:"id"`
	Text     string   `json:"Text"`
	Min      string   `json:"Min"`
	Value    string   `json:"Value"`
	Max      string   `json:"Max"`
	ImageURL string   `json:"ImageURL"`
	Children []Sensor `json:"Children"`
}

// --- Парсинг строки Value в float64 ---
func parseFloat(val string) float64 {
	// Убираем пробелы и единицы измерения
	clean := strings.ReplaceAll(val, "°C", "")
	clean = strings.ReplaceAll(clean, "%", "")
	clean = strings.ReplaceAll(clean, "GB", "")
	clean = strings.ReplaceAll(clean, "RPM", "")
	clean = strings.ReplaceAll(clean, "MB", "")
	clean = strings.TrimSpace(clean)
	// Меняем запятую на точку
	clean = strings.ReplaceAll(clean, ",", ".")
	// Конвертим в float64
	f, err := strconv.ParseFloat(clean, 64)
	if err != nil {
		return 0
	}
	return f
}

// Рекурсивная функция поиска Value по Text
func findValueByText(sensors []Sensor, target int) (float64, bool) {
	for _, sensor := range sensors {
		if sensor.ID == target {
			return parseFloat(sensor.Value), true
		}
		if val, found := findValueByText(sensor.Children, target); found {
			return val, true
		}
	}
	return 0, false
}

// Возвращает Value по Text сразу как float64
func findFloatValueByText(sensors []Sensor, target int) float64 {
	val, found := findValueByText(sensors, target)
	if !found {
		return 0
	}
	return val
}

func fetchInfo(url string) []Sensor {
	resp, err := http.Get(url)
	if err != nil {
		return []Sensor{}
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var root Sensor
	if err := json.Unmarshal(body, &root); err != nil {
		panic(err)
	}

	return []Sensor{root}
}

func floor1(f float64) float64 {
	return math.Floor(f*10) / 10
}

func getWinPCInfo(cfg Config) IPCInfo {
	rawInfo := fetchInfo(cfg.URL)

	ipc := IPCInfo{
		CPU: CPUInfo{
			Load:        findFloatValueByText(rawInfo, cfg.CpuLoad),
			Temperature: findFloatValueByText(rawInfo, cfg.CpuTemp),
			FanSpeed:    findFloatValueByText(rawInfo, cfg.FanSpeed),
		},
		GPU: &GPUInfo{
			Load:        findFloatValueByText(rawInfo, cfg.GpuLoad),
			Temperature: findFloatValueByText(rawInfo, cfg.GpuTemp),
			FanSpeed:    findFloatValueByText(rawInfo, cfg.GpuFanSpeed),
			Memory: GPUMemory{
				Total: floor1(findFloatValueByText(rawInfo, cfg.GpuMemoryTotal) / 1024),
				Used:  floor1(findFloatValueByText(rawInfo, cfg.GpuMemoryUsed) / 1024),
			},
		},
		RAM: RAMInfo{
			Total: floor1(findFloatValueByText(rawInfo, cfg.RAMAvailable) + findFloatValueByText(rawInfo, cfg.RAMUsed)),
			Used:  findFloatValueByText(rawInfo, cfg.RAMUsed),
		},
		// GPU и RAM можно искать другими targetText из конфига или оставить nil
	}

	return ipc

}
