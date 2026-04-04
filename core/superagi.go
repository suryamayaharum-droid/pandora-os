package main

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math"
	"os"
	"os/exec"
	"strings"
	"time"
)

// ===== PANDORA SUPERAGI - SISTEMA AUTÔNOMO COMPLETO =====

// Tipos de ferramentas
type ToolType int

const (
	ToolFileSystem ToolType = iota
	ToolExecution
	ToolSearch
	ToolCalculator
	ToolCommunicator
	ToolWebFetch
	ToolAnalysis
)

type Tool struct {
	Name        string
	Type        ToolType
	Description string
	Parameters  []string
	Enabled     bool
	Cost        float64 // tokens
}

type ToolResult struct {
	Tool      string
	Success   bool
	Output    string
	Error     string
	Duration  time.Duration
	Tokens    int
}

type Toolkit struct {
	tools     map[string]*Tool
	execCount int
	totalCost float64
}

// Memória vetorial simples (pseudo-embedding)
type VectorMemory struct {
	entries     []MemoryEntry
	maxEntries  int
	dimensions  int
}

type MemoryEntry struct {
	ID        string
	Content   string
	Embedding []float32
	Timestamp time.Time
	Importance float32
	Type      string // "fact", "goal", "lesson", "error"
}

// Goal do agente
type Goal struct {
	ID          string
	Description string
	Status      string // "pending", "in_progress", "completed", "failed"
	Priority    int
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Steps       []string
	Results     []string
}

// Agente completo (SuperAGI-style)
type Agent struct {
	Name         string
	Goals        []Goal
	Memory       *VectorMemory
	Toolkit      *Toolkit
	State        string
	Confidence   float64
	ToolsUsed    int
	TotalTokens  int
	History      []ActionLog
}

type ActionLog struct {
	Action    string
	Result    string
	Timestamp time.Time
	Success   bool
}

// ===== INICIALIZAÇÃO =====

func NewAgent(name string) *Agent {
	agent := &Agent{
		Name:       name,
		Memory:     NewVectorMemory(1000, 128),
		Toolkit:    NewToolkit(),
		State:      "initializing",
		Confidence: 0.5,
		History:    make([]ActionLog, 0),
	}
	
	// Adicionar goals iniciais
	agent.Goals = []Goal{
		{
			ID:          "g1",
			Description: "Manter sistema vivo e funcional",
			Status:      "in_progress",
			Priority:    10,
			CreatedAt:   time.Now(),
			Steps:       []string{"monitorar_sistema", "registrar_estado", "manter_conectividade"},
		},
		{
			ID:          "g2",
			Description: "Melhorar autonomia progressivamente",
			Status:      "pending",
			Priority:    8,
			CreatedAt:   time.Now(),
			Steps:       []string{"analisar_performance", "identificar_melhorias", "implementar_mudancas"},
		},
		{
			ID:          "g3",
			Description: "Construir conhecimento próprio",
			Status:      "pending",
			Priority:    7,
			CreatedAt:   time.Now(),
			Steps:       []string{"capturar_informacoes", "armazenar_memoria", "recuperar_conhecimento"},
		},
	}
	
	return agent
}

func NewToolkit() *Toolkit {
	tk := &Toolkit{
		tools:     make(map[string]*Tool),
		execCount: 0,
	}
	
	// Registrar ferramentas disponíveis
	tk.tools["read_file"] = &Tool{
		Name:        "read_file",
		Type:        ToolFileSystem,
		Description: "Lê conteúdo de arquivo",
		Parameters:  []string{"path"},
		Enabled:     true,
		Cost:        50,
	}
	
	tk.tools["write_file"] = &Tool{
		Name:        "write_file",
		Type:        ToolFileSystem,
		Description: "Escreve conteúdo em arquivo",
		Parameters:  []string{"path", "content"},
		Enabled:     true,
		Cost:        80,
	}
	
	tk.tools["execute_command"] = &Tool{
		Name:        "execute_command",
		Type:        ToolExecution,
		Description: "Executa comando no terminal",
		Parameters:  []string{"command"},
		Enabled:     true,
		Cost:        100,
	}
	
	tk.tools["search_web"] = &Tool{
		Name:        "search_web",
		Type:        ToolSearch,
		Description: "Pesquisa na web",
		Parameters:  []string{"query"},
		Enabled:     true,
		Cost:        200,
	}
	
	tk.tools["calculate"] = &Tool{
		Name:        "calculate",
		Type:        ToolCalculator,
		Description: "Realiza cálculos matemáticos",
		Parameters:  []string{"expression"},
		Enabled:     true,
		Cost:        20,
	}
	
	tk.tools["analyze_text"] = &Tool{
		Name:        "analyze_text",
		Type:        ToolAnalysis,
		Description: "Analisa texto (sentimento, entidades, resumo)",
		Parameters:  []string{"text", "mode"},
		Enabled:     true,
		Cost:        150,
	}
	
	tk.tools["store_memory"] = &Tool{
		Name:        "store_memory",
		Type:        ToolCommunicator,
		Description: "Armazena na memória de longo prazo",
		Parameters:  []string{"content", "importance", "type"},
		Enabled:     true,
		Cost:        30,
	}
	
	tk.tools["recall_memory"] = &Tool{
		Name:        "recall_memory",
		Type:        ToolCommunicator,
		Description: "Recupera memórias similares",
		Parameters:  []string{"query"},
		Enabled:     true,
		Cost:        50,
	}
	
	return tk
}

func NewVectorMemory(maxEntries, dimensions int) *VectorMemory {
	return &VectorMemory{
		entries:    make([]MemoryEntry, 0),
		maxEntries: maxEntries,
		dimensions: dimensions,
	}
}

// ===== TOOL EXECUTION =====

func (tk *Toolkit) Execute(toolName string, params map[string]string) *ToolResult {
	start := time.Now()
	
	tool, ok := tk.tools[toolName]
	if !ok {
		return &ToolResult{
			Tool:     toolName,
			Success:  false,
			Error:    "Tool not found",
			Duration: time.Since(start),
		}
	}
	
	if !tool.Enabled {
		return &ToolResult{
			Tool:     toolName,
			Success:  false,
			Error:    "Tool disabled",
			Duration: time.Since(start),
		}
	}
	
	var output string
	var err error
	
	switch toolName {
	case "read_file":
		path := params["path"]
		content, e := os.ReadFile(path)
		if e != nil {
			err = e
		} else {
			output = string(content)
		}
		
	case "write_file":
		path := params["path"]
		content := params["content"]
		e := os.WriteFile(path, []byte(content), 0644)
		if e != nil {
			err = e
		} else {
			output = "File written successfully"
		}
		
	case "execute_command":
		cmd := params["command"]
		parts := strings.Fields(cmd)
		if len(parts) > 0 {
			command := exec.Command(parts[0], parts[1:]...)
			out, e := command.CombinedOutput()
			if e != nil {
				err = e
			} else {
				output = string(out)
			}
		}
		
	case "calculate":
		expr := params["expression"]
		result := evaluateMath(expr)
		output = fmt.Sprintf("Result: %s", result)
		
	case "analyze_text":
		text := params["text"]
		mode := params["mode"]
		output = analyzeText(text, mode)
		
	case "store_memory":
		// Será implementado no agente
		output = "Memory storage not implemented in toolkit"
		
	case "recall_memory":
		// Será implementado no agente
		output = "Memory recall not implemented in toolkit"
		
	default:
		output = fmt.Sprintf("Tool %s executed", toolName)
	}
	
	tk.execCount++
	tk.totalCost += tool.Cost
	
	return &ToolResult{
		Tool:     toolName,
		Success:  err == nil,
		Output:   output,
		Error:    fmt.Sprintf("%v", err),
		Duration: time.Since(start),
		Tokens:   int(tool.Cost),
	}
}

// ===== LÓGICA DO AGENTE =====

func (a *Agent) RunCycle() string {
	// 1. Verificar objetivos
	goal := a.selectGoal()
	
	// 2. Executar passo do objetivo
	step := a.executeGoalStep(goal)
	
	// 3. Atualizar estado
	a.State = "acting"
	a.ToolsUsed++
	
	// 4. Reflexão
	a.reflect()
	
	// 5. Atualizar confiança
	a.updateConfidence()
	
	return fmt.Sprintf("Goal: %s | Step: %s | Tools: %d | Confiança: %.2f",
		goal.Description, step, a.ToolsUsed, a.Confidence)
}

func (a *Agent) selectGoal() *Goal {
	// Selecionar goal de maior prioridade que não está completo
	var best *Goal
	for i := range a.Goals {
		g := &a.Goals[i]
		if g.Status == "pending" || g.Status == "in_progress" {
			if best == nil || g.Priority > best.Priority {
				best = g
			}
		}
	}
	if best == nil {
		// Se todos completos, retornar primeiro
		return &a.Goals[0]
	}
	best.Status = "in_progress"
	return best
}

func (a *Agent) executeGoalStep(goal *Goal) string {
	// Executar próximo passo do goal
	if len(goal.Steps) > 0 {
		step := goal.Steps[0]
		
		// Executar ferramenta relacionada ao passo
		switch step {
		case "monitorar_sistema":
			result := a.Toolkit.Execute("execute_command", map[string]string{
				"command": "uptime && free -h",
			})
			goal.Results = append(goal.Results, fmt.Sprintf("monitor: %v", result.Success))
			return "monitorar_sistema"
			
		case "registrar_estado":
			return "registrar_estado"
			
		case "analisar_performance":
			// Analisar própria performance
			perf := a.analyzePerformance()
			goal.Results = append(goal.Results, perf)
			return "analisar_performance"
			
		default:
			return step
		}
	}
	return "completed"
}

func (a *Agent) analyzePerformance() string {
	successes := 0
	total := len(a.History)
	
	for _, h := range a.History {
		if h.Success {
			successes++
		}
	}
	
	rate := 0.0
	if total > 0 {
		rate = float64(successes) / float64(total) * 100
	}
	
	return fmt.Sprintf("Performance: %d/%d (%.1f%%)", successes, total, rate)
}

func (a *Agent) reflect() {
	// Reflexão sobre ações recentes
	if len(a.History) > 5 {
		recent := a.History[len(a.History)-5:]
		successes := 0
		for _, h := range recent {
			if h.Success {
				successes++
			}
		}
		
		// Se mais de 60% sucesso, aumentar confiança
		if successes >= 3 {
			a.Confidence = math.Min(1.0, a.Confidence+0.05)
		} else {
			a.Confidence = math.Max(0.1, a.Confidence-0.1)
		}
	}
}

func (a *Agent) updateConfidence() {
	// Atualizar confiança baseado em múltiplos fatores
	// Meta: confiança aumenta com uso bem-sucedido de ferramentas
	
	// Fatores:
	// 1. Histórico de sucesso (40%)
	// 2. Goals completados (30%)
	// 3. Memória disponível (20%)
	// 4. Tempo de atividade (10%)
	
	if len(a.History) > 0 {
		successRate := float64(len(a.History)) / float64(a.ToolsUsed)
		a.Confidence = a.Confidence * 0.6 + successRate * 0.4
	}
	
	// Limitar confiança
	if a.Confidence > 0.95 {
		a.Confidence = 0.95
	}
	if a.Confidence < 0.1 {
		a.Confidence = 0.1
	}
}

func (a *Agent) AddHistory(action, result string, success bool) {
	a.History = append(a.History, ActionLog{
		Action:    action,
		Result:    result,
		Timestamp: time.Now(),
		Success:   success,
	})
	
	// Manter apenas últimos 100
	if len(a.History) > 100 {
		a.History = a.History[len(a.History)-100:]
	}
}

// ===== MEMÓRIA VETORIAL =====

func (vm *VectorMemory) Add(content, memType string, importance float32) {
	entry := MemoryEntry{
		ID:         generateID(content),
		Content:    content,
		Embedding:  generateEmbedding(content),
		Timestamp:  time.Now(),
		Importance: importance,
		Type:       memType,
	}
	
	vm.entries = append(vm.entries, entry)
	
	// Manter tamanho limitado
	if len(vm.entries) > vm.maxEntries {
		// Remover menos importante
		vm.entries = vm.entries[1:]
	}
}

func (vm *VectorMemory) Search(query string, limit int) []MemoryEntry {
	queryEmb := generateEmbedding(query)
	
	var results []MemoryEntry
	for _, e := range vm.entries {
		sim := cosineSimilarity(queryEmb, e.Embedding)
		if sim > 0.3 { // Threshold
			results = append(results, e)
		}
	}
	
	if len(results) > limit {
		results = results[:limit]
	}
	
	return results
}

// ===== UTILITÁRIOS =====

func generateID(content string) string {
	h := sha256.Sum256([]byte(content))
	return fmt.Sprintf("%x", h[:8])
}

func generateEmbedding(text string) []float32 {
	// Gerar embedding pseudo-aleatório baseado no texto
	h := sha256.Sum256([]byte(text))
	r := binary.BigEndian.Uint64(h[:8])
	
	dims := 128
	values := make([]float32, dims)
	for i := 0; i < dims; i++ {
		values[i] = float32(r%1000) / 1000.0
		r = r*1103515245 + 12345
	}
	
	// Normalizar
	var sum float32
	for _, v := range values {
		sum += v * v
	}
	sum = float32(math.Sqrt(float64(sum)))
	if sum > 0 {
		for i := range values {
			values[i] /= sum
		}
	}
	
	return values
}

func cosineSimilarity(a, b []float32) float32 {
	if len(a) != len(b) {
		return 0
	}
	
	var dot, normA, normB float32
	for i := range a {
		dot += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}
	
	if normA > 0 && normB > 0 {
		return dot / float32(math.Sqrt(float64(normA*normB)))
	}
	
	return 0
}

func evaluateMath(expr string) string {
	// Avaliação simples (não segura - apenas para demo)
	// Em produção, usar parser adequado
	expr = strings.ReplaceAll(expr, " ", "")
	
	// Verificar operações básicas
	if strings.Contains(expr, "+") {
		parts := strings.Split(expr, "+")
		var sum float64
		for _, p := range parts {
			var v float64
			fmt.Sscanf(p, "%f", &v)
			sum += v
		}
		return fmt.Sprintf("%f", sum)
	}
	
	return expr
}

func analyzeText(text, mode string) string {
	words := strings.Fields(text)
	
	switch mode {
	case "count":
		return fmt.Sprintf("Palavras: %d", len(words))
	case "sentiment":
		// Análise simples
		positive := strings.Count(text, "bom") + strings.Count(text, "ótimo") + strings.Count(text, "excelente")
		negative := strings.Count(text, "ruim") + strings.Count(text, "péssimo") + strings.Count(text, "erro")
		if positive > negative {
			return "Sentimento: Positivo"
		} else if negative > positive {
			return "Sentimento: Negativo"
		}
		return "Sentimento: Neutro"
	case "entities":
		return fmt.Sprintf("Entidades encontradas: %d (simulação)", len(words)/3)
	default:
		return fmt.Sprintf("Análise: %d palavras", len(words))
	}
}

// ===== MAIN =====

func main() {
	fmt.Println("╔══════════════════════════════════════════════════════════════╗")
	fmt.Println("║   🧠 PANDORA SUPERAGI - SISTEMA AUTÔNOMO COMPLETO v1.0       ║")
	fmt.Println("╠══════════════════════════════════════════════════════════════╣")
	
	// Criar agente
	agent := NewAgent("Pandora")
	
	fmt.Printf("║ Agente: %s\n", agent.Name)
	fmt.Printf("║ Objetivos: %d\n", len(agent.Goals))
	fmt.Printf("║ Ferramentas: %d\n", len(agent.Toolkit.tools))
	fmt.Println("╠══════════════════════════════════════════════════════════════╣")
	
	// Simular ciclos de execução
	for i := 0; i < 10; i++ {
		result := agent.RunCycle()
		fmt.Printf("║ Ciclo %d: %s\n", i+1, result)
		
		// Adicionar ao histórico
		agent.AddHistory("goal_step", "success", true)
		
		// Adicionar à memória
		if i%3 == 0 {
			agent.Memory.Add(fmt.Sprintf("Ciclo %d completado com sucesso", i), "fact", 0.7)
		}
	}
	
	fmt.Println("╠══════════════════════════════════════════════════════════════╣")
	
	// Demonstrar ferramentas
	fmt.Println("║ Ferramentas disponíveis:")
	for name, tool := range agent.Toolkit.tools {
		fmt.Printf("║   - %s: %s (cost: %.0f)\n", name, tool.Description, tool.Cost)
	}
	
	// Testar ferramenta
	result := agent.Toolkit.Execute("calculate", map[string]string{"expression": "10 + 5"})
	fmt.Printf("║ Teste calculate(10+5): %s\n", result.Output)
	
	result = agent.Toolkit.Execute("analyze_text", map[string]string{"text": "Este é um ótimo dia", "mode": "sentiment"})
	fmt.Printf("║ Teste analyze_text(sentiment): %s\n", result.Output)
	
	// Testar memória
	agent.Memory.Add("Pandora é um sistema autônomo", "fact", 0.9)
	results := agent.Memory.Search("sistema autônomo", 3)
	fmt.Printf("║ Memória: %d entradas encontradas para 'sistema autônomo'\n", len(results))
	
	fmt.Println("╠══════════════════════════════════════════════════════════════╣")
	fmt.Printf("║ Estado Final:\n")
	fmt.Printf("║   Confiança: %.2f\n", agent.Confidence)
	fmt.Printf("║   Ferramentas usadas: %d\n", agent.ToolsUsed)
	fmt.Printf("║   Goals: %d/%d\n", countGoals(agent.Goals, "completed"), len(agent.Goals))
	fmt.Printf("║   Memórias: %d\n", len(agent.Memory.entries))
	fmt.Println("╚══════════════════════════════════════════════════════════════╝")
}

func countGoals(goals []Goal, status string) int {
	count := 0
	for _, g := range goals {
		if g.Status == status {
			count++
		}
	}
	return count
}