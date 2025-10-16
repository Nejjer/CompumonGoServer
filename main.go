package main

import (
	"encoding/json"
	"fmt"
	"net/http"
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

	// --- Эндпоинты ---
	http.HandleFunc("/getLinuxStat", func(w http.ResponseWriter, r *http.Request) {
		info := getLinuxPCInfo()
		data, _ := json.MarshalIndent(info, "", "  ")
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	})

	http.HandleFunc("/getWinStat", func(w http.ResponseWriter, r *http.Request) {
		info := getWinPCInfo(cfg)
		data, _ := json.MarshalIndent(info, "", "  ")
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	})

	fmt.Printf("Server running at http://localhost:%d\n", cfg.Port)
	http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), nil)
}
