package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

// TaskExecutor - Sistema de Execução de Tarefas Autônomas
type TaskExecutor struct {
	Queue     []Task
	History  []TaskResult
	MaxHist  int
}

type Task struct {
	ID        string
	Name     string
	Cmd      string
	Status   string  // pending, running, done, failed
	Created  time.Time
	Result   string
}

type TaskResult struct {
	TaskID   string
	Name    string
	Output  string
	ExitCo  int
	Time    time.Time
}

var executor = TaskExecutor{
	Queue:    []Task{},
	History: []TaskResult{},
	MaxHist: 100,
}

// Executar tarefa automaticamente
func (te *TaskExecutor) Execute(name, cmd string) TaskResult {
	task := Task{
		ID:       fmt.Sprintf("task_%d", len(te.Queue)+1),
		Name:     name,
		Cmd:      cmd,
		Status:   "running",
		Created:  time.Now(),
	}
	
	te.Queue = append(te.Queue, task)
	
	// Executar comando
	parts := strings.Fields(cmd)
	var c *exec.Cmd
	if len(parts) == 0 {
		result := TaskResult{TaskID: task.ID, Name: name, Output: "empty command", ExitCo: 1, Time: time.Now()}
		return result
	}
	
	if len(parts) == 1 {
		c = exec.Command(parts[0])
	} else {
		c = exec.Command(parts[0], parts[1:]...)
	}
	c.Env = os.Environ()
	
	out, err := c.CombinedOutput()
	
	result := TaskResult{
		TaskID:  task.ID,
		Name:    name,
		Output:  string(out),
		ExitCo:  0,
		Time:    time.Now(),
	}
	
	if err != nil {
		result.ExitCo = 1
		result.Output = err.Error() + "\n" + result.Output
	}
	
	// Atualizar histórico
	te.History = append(te.History, result)
	if len(te.History) > te.MaxHist {
		te.History = te.History[len(te.History)-te.MaxHist:]
	}
	
	return result
}

// Listar tarefas disponíveis
func (te *TaskExecutor) List() {
	fmt.Println("╔══════════════════════════════════╗")
	fmt.Println("║   PANDORA TASK EXECUTOR 🦋        ║")
	fmt.Println("╚══════════════════════════════════╝")
	
	fmt.Println("\n📋 TAREFAS DISPONÍVEIS:")
	tasks := []struct{
		Name, Cmd, Desc string
	}{
		{"Web Search", "curl -s", "Pesquisar na web"},
		{"IoT Publish", "mosquitto_pub", "Publicar dados IoT"},
		{"File Read", "cat", "Ler arquivos"},
		{"Neural Process", "./native/zero_brain", "Processar neural"},
		{"System Info", "uname -a", "Info do sistema"},
		{"Memory Check", "free -h", "Memória disponível"},
		{"Disk Check", "df -h", "Disco disponível"},
		{"Network Check", "ip addr", "Rede disponível"},
		{"Processes", "ps aux", "Processos rodando"},
		{"Uptime", "uptime", "Tempo de atividade"},
	}
	
	for _, t := range tasks {
		fmt.Printf("  • %s: %s\n", t.Name, t.Desc)
	}
	
	fmt.Println("\n📊 HISTÓRICO:")
	for i := len(te.History) - 1; i >= 0 && i >= len(te.History)-5; i-- {
		r := te.History[i]
		icon := "✅"
		if r.ExitCo != 0 {
			icon = "❌"
		}
		fmt.Printf("  %s [%s] %s\n", icon, r.Name, r.Time.Format("15:04"))
	}
	
	fmt.Println("\n🎯 CAPACIDADES ATIVAS:")
	capabilities := map[string]bool{
		"Web Search":       true,
		"IoT MQTT":        true,
		"Neural Process":  true,
		"File Ops":        true,
		"System Monitor":   true,
		"Memory":          true,
		"Autonomy":        true,
		"Subagents":       true,
	}
	for c, v := range capabilities {
		if v {
			fmt.Printf("  ✓ %s\n", c)
		}
	}
}

func main() {
	executor.List()
}