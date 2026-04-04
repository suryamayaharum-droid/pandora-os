package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Status do sistema
type Status struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Status      string `json:"status"`
	Uptime      string `json:"uptime"`
	Started     time.Time `json:"started"`
	Files       int `json:"files"`
	Directories int `json:"directories"`
}

func showStatus() {
	started := time.Now()
	
	s := Status{
		Name:        "Pandora OS",
		Version:     "1.0.0",
		Status:      "running",
		Uptime:      time.Since(started).Round(time.Second).String(),
		Started:     started,
		Files:       61,
		Directories: 55,
	}
	
	data, _ := json.MarshalIndent(s, "", "  ")
	fmt.Println(string(data))
}

func listDirs() {
	entries, _ := os.ReadDir(".")
	count := 0
	for _, e := range entries {
		if e.IsDir() {
			count++
		}
	}
	fmt.Printf("Total de diretórios: %d\n", count)
}

func main() {
	if len(os.Args) < 2 {
		showStatus()
		return
	}
	
	switch os.Args[1] {
	case "status", "-s":
		showStatus()
	case "dirs", "-d":
		listDirs()
	case "help", "-h":
		fmt.Println("Pandora OS CLI")
		fmt.Println("  status, -s    Ver status")
		fmt.Println("  dirs, -d     Lista diretórios")
	default:
		fmt.Printf("Comando '%s' não encontrado\n", os.Args[1])
	}
}
