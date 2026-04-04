// Pandora OS - Sistema Real
// Uso: go run app.go <comando>

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
	
	_ "pandora-os/autonomy"
	_ "pandora-os/core"
	_ "pandora-os/agent"
)

type Status struct {
	Agent      string    `json:"agent"`
	Version   string    `json:"version"`
	Status    string    `json:"status"`
	Uptime    string    `json:"uptime"`
	TasksRun  int       `json:"tasks_run"`
	Started   time.Time `json:"started"`
}

func main() {
	start := time.Now()
	
	status := Status{
		Agent:    "Pandora OS",
		Version: "1.0",
		Status:  "running",
		Uptime:  time.Since(start).Round(time.Second).String(),
		TasksRun: 0,
		Started: start,
	}
	
	// Parsear argumento
	if len(os.Args) < 2 {
		showStatus(status)
		return
	}
	
	cmd := os.Args[1]
	
	switch cmd {
	case "status":
		showStatus(status)
	case "tasks":
		fmt.Println("Tarefas: sistema operacional ativo")
	case "help":
		showHelp()
	default:
		fmt.Printf("Comando '%s' não reconhecido\n", cmd)
		showHelp()
	}
}

func showStatus(s Status) {
	data, _ := json.MarshalIndent(s, "", "  ")
	fmt.Println(string(data))
}

func showHelp() {
	fmt.Println("Comandos disponíveis:")
	fmt.Println("  status  - Ver status do sistema")
	fmt.Println("  tasks   - Lista de tarefas ativas")
	fmt.Println("  help    - Mostrar esta ajuda")
}
