package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type IPCInfo struct {
	CPU CPUInfo  `json:"cpu"`
	GPU *GPUInfo `json:"gpu,omitempty"`
	RAM RAMInfo  `json:"ram"`
}

type CPUInfo struct {
	Load        float64 `json:"load"`
	Temperature float64 `json:"temperature"`
	FanSpeed    float64 `json:"fanSpeed"`
}

type GPUInfo struct {
	Load        float64   `json:"load"`
	Temperature float64   `json:"temperature"`
	Memory      GPUMemory `json:"memory"`
	FanSpeed    float64   `json:"fanSpeed"`
}

type GPUMemory struct {
	Total float64 `json:"total"`
	Used  float64 `json:"used"`
}

type RAMInfo struct {
	Total float64 `json:"total"`
	Used  float64 `json:"used"`
}

func main() {
	// --- Чтение конфига ---
	cfgFile, err := os.ReadFile("config.json")
	if err != nil {
		panic(err)
	}

	var cfg Config
	if err := json.Unmarshal(cfgFile, &cfg); err != nil {
		panic(err)
	}

	linInfo := getLinuxPCInfo()
	winInfo := getWinPCInfo(cfg)

	data, _ := json.MarshalIndent(linInfo, "", "  ")
	fmt.Println("LINUX")
	fmt.Println(string(data))

	winData, _ := json.MarshalIndent(winInfo, "", "  ")
	fmt.Println("WIN")
	fmt.Println(string(winData))
}
