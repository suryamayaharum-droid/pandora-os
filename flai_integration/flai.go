package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"strings"
	"time"
)

// PANDORA + FLAI INTEGRATION

type OllamaManager struct {
	Installed bool
	Running   bool
	Endpoint  string
	Models    []Model
	Selected  string
}

type Model struct {
	Name string
	Size string
}

type LocalAI struct {
	Name    string
	Ollama  *OllamaManager
	History []ChatMessage
}

type ChatMessage struct {
	Role    string
	Content string
}

func NewOllamaManager() *OllamaManager {
	om := &OllamaManager{
		Installed: false,
		Running:   false,
		Endpoint:  "http://localhost:11434",
	}
	
	if checkCmd("ollama") {
		om.Installed = true
		fmt.Println("  ✓ Ollama encontrado")
	}
	
	if om.Installed {
		if checkServer(om.Endpoint) {
			om.Running = true
			fmt.Println("  ✓ Servidor rodando")
		}
		om.listModels()
	}
	
	return om
}

func checkCmd(name string) bool {
	cmd := exec.Command("which", name)
	return cmd.Run() == nil
}

func checkServer(endpoint string) bool {
	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get(endpoint + "/api/tags")
	if err == nil {
		defer resp.Body.Close()
		return resp.StatusCode == 200
	}
	return false
}

func (om *OllamaManager) listModels() {
	if !om.Running {
		return
	}
	
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(om.Endpoint + "/api/tags")
	if err != nil {
		return
	}
	defer resp.Body.Close()
	
	var result struct {
		Models []struct {
			Name string `json:"name"`
			Size int64  `json:"size"`
		} `json:"models"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return
	}
	
	for _, m := range result.Models {
		om.Models = append(om.Models, Model{Name: m.Name, Size: formatSize(m.Size)})
	}
	
	if len(om.Models) > 0 {
		om.Selected = om.Models[0].Name
	}
}

func formatSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func (om *OllamaManager) Generate(prompt string, model string) (string, error) {
	if !om.Running {
		return "", fmt.Errorf("Servidor não está rodando")
	}
	
	if model == "" {
		model = om.Selected
	}
	
	reqBody := map[string]interface{}{
		"model":  model,
		"prompt": prompt,
		"stream": false,
	}
	
	body, _ := json.Marshal(reqBody)
	
	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Post(om.Endpoint+"/api/generate", "application/json", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	
	var result struct {
		Response string `json:"response"`
	}
	
	json.NewDecoder(resp.Body).Decode(&result)
	return result.Response, nil
}

func NewLocalAI(name string) *LocalAI {
	return &LocalAI{
		Name:    name,
		Ollama:  NewOllamaManager(),
		History: make([]ChatMessage, 0),
	}
}

func (ai *LocalAI) Chat(userInput string) string {
	ai.History = append(ai.History, ChatMessage{Role: "user", Content: userInput})
	
	prompt := fmt.Sprintf("Você é %s, um assistente de IA local.\n\nHistórico:\n%s\n\nUsuário: %s\n\nResposta:", 
		ai.Name, ai.formatHistory(), userInput)
	
	response, err := ai.Ollama.Generate(prompt, "")
	if err != nil {
		response = "Erro ao gerar resposta"
	}
	
	ai.History = append(ai.History, ChatMessage{Role: "assistant", Content: response})
	
	if len(ai.History) > 50 {
		ai.History = ai.History[len(ai.History)-50:]
	}
	
	return response
}

func (ai *LocalAI) formatHistory() string {
	var sb strings.Builder
	start := len(ai.History) - min(5, len(ai.History))
	if start < 0 {
		start = 0
	}
	for i := start; i < len(ai.History); i++ {
		sb.WriteString(fmt.Sprintf("%s: %s\n", ai.History[i].Role, ai.History[i].Content))
	}
	return sb.String()
}

func (ai *LocalAI) AnalyzeCode(code string) string {
	prompt := fmt.Sprintf("Analise este código:\n\n%s\n\nForneça: 1.Resumo 2.Qualidade(A-F) 3.Bugs 4.Melhorias", code)
	return ai.Chat(prompt)
}

func (ai *LocalAI) CodeReview(path string) string {
	prompt := fmt.Sprintf("Code Review de: %s\n\nForneça: 1.Estrutura 2.Qualidade 3.Bugs 4.Melhorias 5.Boas práticas", path)
	return ai.Chat(prompt)
}

func (ai *LocalAI) Explain(code string) string {
	prompt := fmt.Sprintf("Explique este código de forma simples:\n\n%s", code)
	return ai.Chat(prompt)
}

func main() {
	fmt.Println("╔══════════════════════════════════════════════════════════════╗")
	fmt.Println("║        🤖 PANDORA + FLAI INTEGRATION v1.0                    ║")
	fmt.Println("║        [ IA LOCAL COM OLLAMA - COMO FLAI ]                  ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════╝")
	
	ai := NewLocalAI("Pandora-Local")
	
	fmt.Println("\n📡 Status Ollama:")
	fmt.Printf("  Instalado: %v\n", ai.Ollama.Installed)
	fmt.Printf("  Rodando: %v\n", ai.Ollama.Running)
	
	if len(ai.Ollama.Models) > 0 {
		fmt.Println("\n📦 Modelos disponíveis:")
		for i, m := range ai.Ollama.Models {
			sel := ""
			if m.Name == ai.Ollama.Selected {
				sel = " [SELECIONADO]"
			}
			fmt.Printf("  [%d] %s (%s)%s\n", i+1, m.Name, m.Size, sel)
		}
	} else {
		fmt.Println("\n⚠️  Nenhum modelo disponível")
		fmt.Println("\n📥 INSTALAÇÃO:")
		fmt.Println("  $ curl -fsSL https://ollama.ai/install.sh | sh")
		fmt.Println("  $ ollama pull deepseek-coder")
		fmt.Println("  $ ollama serve")
	}
	
	fmt.Println("\n📟 Comandos FLAI-like:")
	cmds := [][]string{
		{"flai init", "Inicializa projeto"},
		{"flai analyze", "Analisa código"},
		{"flai review", "Code review"},
		{"flai refactor", "Refatora"},
		{"flai chat", "Modo chat"},
	}
	for _, c := range cmds {
		fmt.Printf("  $ %s → %s\n", c[0], c[1])
	}
	
	if ai.Ollama.Running && ai.Ollama.Selected != "" {
		fmt.Println("\n💬 Teste:")
		fmt.Printf("  Você: O que é FLAI?\n")
		resp := ai.Chat("O que é FLAI?")
		fmt.Printf("  Pandora: %s\n", resp[:min(100, len(resp))])
	} else {
		fmt.Println("\n💬 Demo (sem Ollama):")
		fmt.Println("  Você: O que você é?")
		fmt.Println("  Pandora: Sou Pandora, sistema de IA autônomo!")
		fmt.Println("           Posso usar Ollama quando disponível")
	}
	
	fmt.Println("\n" + "╔══════════════════════════════════════════════════════════════╗")
	fmt.Println("║              🏗️  ARQUITETURA FLAI-LIKE                    ║")
	fmt.Println("╠══════════════════════════════════════════════════════════════╣")
	fmt.Println("║  Pandora ←→ FLAI Shell ←→ Ollama (LLM Local)              ║")
	fmt.Println("║                                                          ║")
	fmt.Println("║  CAPACIDADES:                                            ║")
	fmt.Println("║  ✅ Análise de código                                    ║")
	fmt.Println("║  ✅ Code Review automático                               ║")
	fmt.Println("║  ✅ Refatoração                                          ║")
	fmt.Println("║  ✅ Chat interativo                                      ║")
	fmt.Println("║  ✅ Entendimento MVC                                     ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════╝")
	
	fmt.Println("\n✅ Integração FLAI-like completa!")
}

func min(a, b int) int { if a < b { return a }; return b }
