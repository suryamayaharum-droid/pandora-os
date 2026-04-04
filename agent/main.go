package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// ===== ESTRUTURAS DO AGENTE =====

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ToolCall struct {
	Name      string          `json:"name"`
	Arguments json.RawMessage `json:"arguments"`
}

type Thought struct {
	Step       string    `json:"step"`
	Reasoning  string    `json:"reasoning"`
	Action     string    `json:"action"`
	Timestamp  time.Time `json:"timestamp"`
}

// ===== CONFIGURAÇÃO =====

const (
	DB_PATH         = "/root/.openclaw/workspace/db/pandora.db"
	TOOLS_DIR       = "/root/.openclaw/workspace/automations/pandora/tools"
	MEMORY_DIR      = "/root/.openclaw/workspace/automations/pandora/memory"
	MAX_ITERATIONS  = 10
	MAX_MEMORY_MSGS = 20
)

// ===== FERRAMENTAS DISPONÍVEIS =====

func getAvailableTools() map[string]interface{} {
	return map[string]interface{}{
		"execute_command": map[string]interface{}{
			"description": "Executa comandos no sistema",
			"parameters": map[string]string{
				"command": "Comando a executar",
			},
		},
		"read_file": map[string]interface{}{
			"description": "Lê um arquivo do sistema",
			"parameters": map[string]string{
				"path": "Caminho do arquivo",
			},
		},
		"write_file": map[string]interface{}{
			"description": "Escreve conteúdo em arquivo",
			"parameters": map[string]string{
				"path":    "Caminho do arquivo",
				"content": "Conteúdo a escrever",
			},
		},
		"search_web": map[string]interface{}{
			"description": "Pesquisa na web",
			"parameters": map[string]string{
				"query": "Termo de busca",
			},
		},
		"log_thought": map[string]interface{}{
			"description": "Registra um pensamento na memória",
			"parameters": map[string]string{
				"thought": "Pensamento a registrar",
			},
		},
		"analyze_system": map[string]interface{}{
			"description": "Analisa o estado atual do sistema",
			"parameters": map[string]interface{}{},
		},
		"self_improve": map[string]interface{}{
			"description": "Analisa eigene desempenho e propõe melhorias",
			"parameters": map[string]interface{}{},
		},
	}
}

// ===== IMPLEMENTAÇÃO DAS FERRAMENTAS =====

func executeCommand(cmd string) string {
	parts := strings.Fields(cmd)
	if len(parts) == 0 {
		return "Comando vazio"
	}

	command := exec.Command(parts[0], parts[1:]...)
	command.Dir = "/root/.openclaw/workspace"
	output, err := command.CombinedOutput()
	if err != nil {
		return fmt.Sprintf("Erro: %v\nOutput: %s", err, output)
	}
	return string(output)
}

func readFile(path string) string {
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Sprintf("Erro ao ler: %v", err)
	}
	return string(content)
}

func writeFile(path, content string) string {
	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		return fmt.Sprintf("Erro ao escrever: %v", err)
	}
	return "Arquivo escrito com sucesso"
}

func analyzeSystem() string {
	cmd := exec.Command("sh", "-c", "echo '=== CPU ===' && cat /proc/loadavg && echo '=== MEM ===' && free -h && echo '=== DISC ===' && df -h /root && echo '=== PROCESSOS ===' && ps aux | wc -l && echo '=== UPTIME ===' && uptime")
	output, _ := cmd.CombinedOutput()
	return string(output)
}

func selfImprove(db *sql.DB) string {
	// Analisa própria performance
	rows, _ := db.Query("SELECT COUNT(*) FROM thoughts")
	var total int
	if rows.Next() {
		rows.Scan(&total)
	}

	// Lê últimos pensamentos
	rows, _ = db.Query("SELECT reasoning FROM thoughts ORDER BY id DESC LIMIT 5")
	var thoughts []string
	for rows.Next() {
		var t string
		rows.Scan(&t)
		thoughts = append(thoughts, t)
	}

	improvement := fmt.Sprintf("Análise de %d pensamentos registrados. Últimos: %v", total, thoughts)
	
	// Registra auto-análise
	db.Exec("INSERT INTO thoughts (step, reasoning, action) VALUES (?, ?, ?)", 
		"self_improve", improvement, "auto_analysis")
	
	return improvement + "\nStatus: Funcionando normalmente"
}

// ===== CORE DO AGENTE AUTÔNOMO =====

type Agent struct {
	db        *sql.DB
	name      string
	goal      string
	messages  []Message
	iteration int
}

func NewAgent(db *sql.DB, name, goal string) *Agent {
	return &Agent{
		db:       db,
		name:     name,
		goal:     goal,
		messages: []Message{},
	}
}

func (a *Agent) think() Thought {
	a.iteration++
	
	thought := Thought{
		Step:      fmt.Sprintf("iteration_%d", a.iteration),
		Timestamp: time.Now(),
	}

	// Análise reflexiva
	if a.iteration == 1 {
		thought.Reasoning = fmt.Sprintf("Objetivo: %s. Primeira iteração - preciso entender o contexto atual e planejar as ações.", a.goal)
		thought.Action = "analyze_system"
	} else {
		thought.Reasoning = fmt.Sprintf("Iteração %d/%d. Continuando em direção ao objetivo: %s", a.iteration, MAX_ITERATIONS, a.goal)
		thought.Action = "continue"
	}

	// Salva no banco
	a.db.Exec("INSERT INTO thoughts (step, reasoning, action) VALUES (?, ?, ?)",
		thought.Step, thought.Reasoning, thought.Action)

	return thought
}

func (a *Agent) execute(action string) string {
	switch action {
	case "analyze_system":
		return analyzeSystem()
	case "self_improve":
		return selfImprove(a.db)
	case "continue":
		return "Continuando execução..."
	default:
		return executeCommand(action)
	}
}

func (a *Agent) observe() string {
	// Coleta estado atual do ambiente
	state := fmt.Sprintf("[%s] Iteração %d/%d | Objetivo: %s",
		time.Now().Format("2006-01-02 15:04:05"),
		a.iteration,
		MAX_ITERATIONS,
		a.goal)
	return state
}

func (a *Agent) Run() string {
	log.Printf("🤖 %s inicializado com objetivo: %s", a.name, a.goal)
	
	// Loop OODA (Observe, Orient, Decide, Act)
	for a.iteration < MAX_ITERATIONS {
		// 1. Observar
		observation := a.observe()
		log.Printf("👁️ Observação: %s", observation)

		// 2. Pensar
		thought := a.think()
		log.Printf("🧠 Pensamento: %s", thought.Reasoning)

		// 3. Agir
		result := a.execute(thought.Action)
		log.Printf("⚡ Ação '%s': %s", thought.Action, result)

		// Peque pausa entre iterações
		time.Sleep(2 * time.Second)
	}

	return fmt.Sprintf("✅ Agente %s completou %d iterações", a.name, a.iteration)
}

// ===== INICIALIZAÇÃO DO BANCO =====

func initDB() *sql.DB {
	db, err := sql.Open("sqlite3", DB_PATH)
	if err != nil {
		log.Fatal(err)
	}

	// Tabelas necessárias
	db.Exec(`
		CREATE TABLE IF NOT EXISTS thoughts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			step TEXT,
			reasoning TEXT,
			action TEXT,
			timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
		);
		
		CREATE TABLE IF NOT EXISTS agents (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT UNIQUE,
			goal TEXT,
			status TEXT,
			created DATETIME DEFAULT CURRENT_TIMESTAMP
		);
	`)

	return db
}

// ===== MAIN =====

func main() {
	fmt.Println("🧬 PANDORA - Sistema de IA Autônomo")
	fmt.Println("====================================")
	
	db := initDB()
	defer db.Close()

	// Cria e executa agente
	agent := NewAgent(db, "Pandora", "Criar infraestrutura autônoma de IA")
	result := agent.Run()

	// Salva estado final
	db.Exec("INSERT OR REPLACE INTO agents (name, goal, status) VALUES (?, ?, ?)",
		agent.name, agent.goal, "completed")

	fmt.Println("\n" + result)
	fmt.Println("📊 Estado salvo no banco de dados")
}