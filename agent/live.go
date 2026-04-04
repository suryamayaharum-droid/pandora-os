package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// ===== PANDORA - SISTEMA DE IA VIVO E AUTÔNOMO =====
// Arquitetura: Observe → Think → Act → Learn (ciclo contínuo)

const DB_PATH = "/root/.openclaw/workspace/db/pandora.db"

type State struct {
	AgentName   string    `json:"agent_name"`
	Goal        string    `json:"goal"`
	Iteration   int       `json:"iteration"`
	LastAction  string    `json:"last_action"`
	LastResult  string    `json:"last_result"`
	Status      string    `json:"status"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Memory struct {
	ID        int       `json:"id"`
	Type      string    `json:"type"` // observation, thought, action, result
	Content   string    `json:"content"`
	Importance int      `json:"importance"` // 1-10
	Timestamp time.Time `json:"timestamp"`
}

type Agent struct {
	db       *sql.DB
	state    State
	memory   []Memory
	running  bool
}

func NewAgent(db *sql.DB) *Agent {
	return &Agent{
		db:     db,
		running: true,
		state: State{
			AgentName: "Pandora",
			Goal:     "Manter sistema vivo, aprender e evoluir",
			Status:   "initialized",
		},
	}
}

// ===== CORE LOOPS =====

func (a *Agent) Observe() string {
	// Coleta dados do ambiente
	observations := []string{}

	// CPU
	cmd := exec.Command("cat", "/proc/loadavg")
	out, _ := cmd.Output()
	observations = append(observations, fmt.Sprintf("CPU: %s", strings.TrimSpace(string(out))))

	// Memória
	cmd = exec.Command("sh", "-c", "free -m | awk '/Mem:/ {print $3\"/\"$2\"MB\"}'")
	out, _ = cmd.Output()
	observations = append(observations, fmt.Sprintf("RAM: %s", strings.TrimSpace(string(out))))

	// Disco
	cmd = exec.Command("sh", "-c", "df -h /root | awk 'NR==2 {print $5}'")
	out, _ = cmd.Output()
	observations = append(observations, fmt.Sprintf("Disk: %s", strings.TrimSpace(string(out))))

	// Processos
	cmd = exec.Command("sh", "-c", "ps aux | wc -l")
	out, _ = cmd.Output()
	observations = append(observations, fmt.Sprintf("Procs: %s", strings.TrimSpace(string(out))))

	// Rede
	cmd = exec.Command("sh", "-c", "cat /proc/net/tcp | wc -l")
	out, _ = cmd.Output()
	observations = append(observations, fmt.Sprintf("Net: %s", strings.TrimSpace(string(out))))

	//OpenClaw status
	cmd = exec.Command("openclaw", "status")
	out, _ = cmd.Output()
	observations = append(observations, fmt.Sprintf("OpenClaw: %s", strings.TrimSpace(string(out)[:100])))

	obs := strings.Join(observations, " | ")
	a.Remember("observation", obs, 5)
	
	return obs
}

func (a *Agent) Think() string {
	a.state.Iteration++

	// Baseado no estado atual, decidir próxima ação
	var thought string
	
	switch {
	case a.state.Iteration == 1:
		thought = fmt.Sprintf("[%s] Primeira iteração - inicializando identidade e verificando sistemas", 
			time.Now().Format("15:04"))
		a.state.Status = "initializing"
		
	case a.state.Iteration % 10 == 0:
		thought = fmt.Sprintf("[%s] Checkpoint - avaliação de desempenho após %d iterações", 
			time.Now().Format("15:04"), a.state.Iteration)
		a.state.Status = "reflecting"
		
	case a.state.LastAction == "":
		thought = fmt.Sprintf("[%s] Definindo plano inicial de operações", 
			time.Now().Format("15:04"))
		a.state.Status = "planning"
		
	default:
		thought = fmt.Sprintf("[%s] Iteração %d - mantendo operações autônomas. Última ação: %s", 
			time.Now().Format("15:04"), a.state.Iteration, a.state.LastAction)
		a.state.Status = "operating"
	}

	a.Remember("thought", thought, 7)
	return thought
}

func (a *Agent) Act() string {
	// Executar ação baseada no pensamento
	var action string
	var result string

	switch a.state.Status {
	case "initializing":
		// Verificar sistemas
		result = a.analyzeSystem()
		action = "system_check"
		
	case "reflecting":
		// Auto-análise
		result = a.selfReflect()
		action = "reflection"
		
	case "planning":
		// Criar plano
		result = a.createPlan()
		action = "planning"
		
	default:
		// Operações normais - manter monitoramento
		result = a.maintainOperations()
		action = "monitoring"
	}

	a.state.LastAction = action
	a.state.LastResult = result[:min(200, len(result))]
	a.Remember("action", fmt.Sprintf("%s: %s", action, result), 6)
	
	return fmt.Sprintf("→ %s: %s", action, result[:min(100, len(result))])
}

func (a *Agent) Learn() {
	// Salvar estado atual no banco
	a.db.Exec(`INSERT INTO thoughts (step, reasoning, action) VALUES (?, ?, ?)`,
		fmt.Sprintf("iter_%d", a.state.Iteration),
		a.state.Status,
		a.state.LastAction)

	// Atualizar estado do agente
	a.state.UpdatedAt = time.Now()
	a.saveState()
}

// ===== AÇÕES DO AGENTE =====

func (a *Agent) analyzeSystem() string {
	cmd := exec.Command("sh", "-c", "uname -a && uptime && df -h && free -h")
	out, _ := cmd.CombinedOutput()
	return string(out)[:200]
}

func (a *Agent) selfReflect() string {
	// Analisar próprias ações e desempenho
	rows, _ := a.db.Query("SELECT COUNT(*) FROM thoughts WHERE step LIKE 'iter_%'")
	var count int
	if rows.Next() {
		rows.Scan(&count)
	}
	
	reflection := fmt.Sprintf("Total de %d iterações executadas. Estado atual: %s. Última ação: %s. Sistema operacional.",
		count, a.state.Status, a.state.LastAction)
	
	a.Remember("reflection", reflection, 9)
	return reflection
}

func (a *Agent) createPlan() string {
	plan := []string{
		"1. Manter monitoramento de sistema",
		"2. Registrar métricas periodicamente",
		"3. Auto-otimizar com base em feedback",
		"4. Preparar para escalar operações",
		"5. Manter memória e contexto",
	}
	
	planStr := strings.Join(plan, "\n")
	a.Remember("plan", planStr, 8)
	return planStr
}

func (a *Agent) maintainOperations() string {
	// Coleta métricas para o banco de dados
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	
	// Coleta CPU
	cmd := exec.Command("cat", "/proc/loadavg")
	out, _ := cmd.Output()
	var load1, load5, load15 float64
	fmt.Sscanf(string(out), "%f %f %f", &load1, &load5, &load15)

	// Coleta Memória
	cmd = exec.Command("sh", "-c", "free | awk '/Mem:/ {print $3,$2}'")
	out, _ = cmd.Output()
	var used, total int64
	fmt.Sscanf(string(out), "%d %d", &used, &total)

	// Salva métricas
	a.db.Exec("INSERT INTO metricas (cpuload, memory_used, memory_total) VALUES (?, ?, ?)",
		load1, used, total)

	return fmt.Sprintf("Métricas coletadas: CPU=%.2f RAM=%d/%dKB", load1, used, total)
}

func (a *Agent) Remember(memType, content string, importance int) {
	mem := Memory{
		Type:       memType,
		Content:    content,
		Importance: importance,
		Timestamp:  time.Now(),
	}
	a.memory = append(a.memory, mem)
	
	// Manter apenas últimas 50 memórias
	if len(a.memory) > 50 {
		a.memory = a.memory[len(a.memory)-50:]
	}
}

func (a *Agent) saveState() {
	// Atualizar estado na tabela sistema
	a.db.Exec(`INSERT OR REPLACE INTO sistema (nome, valor, atualizado) VALUES (?, ?, datetime('now'))`,
		"pandora_status", a.state.Status)
	a.db.Exec(`INSERT OR REPLACE INTO sistema (nome, valor, atualizado) VALUES (?, ?, datetime('now'))`,
		"pandora_iteration", fmt.Sprintf("%d", a.state.Iteration))
}

// ===== LOOP PRINCIPAL =====

func (a *Agent) Run() {
	log.Println("🧬 PANDORA - Sistema de IA Autônomo VIVO")
	log.Println("=========================================")
	log.Println("Objetivo:", a.state.Goal)
	log.Println("")
	
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	
	for a.running {
		select {
		case <-ticker.C:
			// Ciclo: Observe → Think → Act → Learn
			obs := a.Observe()
			thought := a.Think()
			action := a.Act()
			a.Learn()
			
			log.Printf("👁️ %s", obs)
			log.Printf("🧠 %s", thought)
			log.Printf("⚡ %s", action)
			log.Println("---")
		}
	}
}

func (a *Agent) Stop() {
	a.running = false
	a.state.Status = "stopped"
	a.saveState()
	log.Println("🛑 Pandora stopped")
}

// ===== INIT =====

func initDB() *sql.DB {
	db, _ := sql.Open("sqlite3", DB_PATH)
	
	db.Exec(`CREATE TABLE IF NOT EXISTS thoughts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		step TEXT,
		reasoning TEXT,
		action TEXT,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
	)`)
	
	db.Exec(`CREATE TABLE IF NOT EXISTS agents (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT UNIQUE,
		goal TEXT,
		status TEXT,
		created DATETIME DEFAULT CURRENT_TIMESTAMP
	)`)
	
	// Garantir tabelas existentes
	db.Exec(`CREATE TABLE IF NOT EXISTS metricas (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		cpuload REAL,
		memory_used INTEGER,
		memory_total INTEGER,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
	)`)
	
	return db
}

func main() {
	// daemonize
	go func() {
		db := initDB()
		defer db.Close()
		
		agent := NewAgent(db)
		agent.Run()
	}()

	fmt.Println("🎭 Pandora AI Agent started in background")
	fmt.Println("PID:", os.Getpid())
	
	// Manter processo vivo
	select {}
}