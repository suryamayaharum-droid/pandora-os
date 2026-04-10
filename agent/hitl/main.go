package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type Stage string

const (
	StageInput   Stage = "input"
	StagePlan    Stage = "plan"
	StageReview  Stage = "human_review"
	StageExecute Stage = "execute"
	StageReport  Stage = "report"
)

type RunState struct {
	Goal            string    `json:"goal"`
	Plan            string    `json:"plan"`
	ApprovedByHuman bool      `json:"approved_by_human"`
	ExecutionResult string    `json:"execution_result"`
	CurrentStage    Stage     `json:"current_stage"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type chatRequest struct {
	Model    string        `json:"model"`
	Messages []chatMessage `json:"messages"`
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatResponse struct {
	Choices []struct {
		Message chatMessage `json:"message"`
	} `json:"choices"`
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	state := RunState{CurrentStage: StageInput, UpdatedAt: time.Now()}

	fmt.Println("=== Agent MVP com HITL ===")
	fmt.Print("Digite o objetivo da tarefa: ")
	goal, _ := reader.ReadString('\n')
	state.Goal = strings.TrimSpace(goal)

	state.CurrentStage = StagePlan
	state.UpdatedAt = time.Now()
	state.Plan = buildPlan(state.Goal)
	fmt.Println("\nPlano proposto pelo agente:")
	fmt.Println(state.Plan)

	state.CurrentStage = StageReview
	state.UpdatedAt = time.Now()
	fmt.Print("\nAprovar execução externa? (sim/nao): ")
	answer, _ := reader.ReadString('\n')
	answer = strings.TrimSpace(strings.ToLower(answer))
	state.ApprovedByHuman = answer == "sim" || answer == "s" || answer == "yes"

	if !state.ApprovedByHuman {
		fmt.Println("\nExecução bloqueada: ação externa requer aprovação humana.")
		persistState(state)
		return
	}

	state.CurrentStage = StageExecute
	state.UpdatedAt = time.Now()
	state.ExecutionResult = executePlan(state.Plan)

	state.CurrentStage = StageReport
	state.UpdatedAt = time.Now()
	fmt.Println("\n=== Relatório final ===")
	fmt.Printf("Objetivo: %s\n", state.Goal)
	fmt.Printf("Aprovado por humano: %t\n", state.ApprovedByHuman)
	fmt.Printf("Resultado: %s\n", state.ExecutionResult)

	persistState(state)
}

func buildPlan(goal string) string {
	prompt := fmt.Sprintf("Crie um plano curto em 3 passos para atingir este objetivo: %s", goal)
	if generated, err := callOpenAICompat(prompt); err == nil && strings.TrimSpace(generated) != "" {
		return generated
	}
	return "1) Entender requisitos\n2) Executar passos com validação humana\n3) Entregar relatório com métricas"
}

func executePlan(plan string) string {
	return "Plano executado em modo seguro (simulação controlada): " + strings.Split(plan, "\n")[0]
}

func callOpenAICompat(prompt string) (string, error) {
	baseURL := os.Getenv("OPENAI_COMPAT_BASE_URL")
	model := os.Getenv("OPENAI_COMPAT_MODEL")
	apiKey := os.Getenv("OPENAI_COMPAT_API_KEY")
	if baseURL == "" || model == "" {
		return "", fmt.Errorf("base URL/modelo não configurados")
	}

	payload := chatRequest{
		Model: model,
		Messages: []chatMessage{
			{Role: "system", Content: "Você é um agente operacional objetivo."},
			{Role: "user", Content: prompt},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodPost, strings.TrimRight(baseURL, "/")+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	if apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+apiKey)
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		content, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("erro %d: %s", resp.StatusCode, strings.TrimSpace(string(content)))
	}

	var out chatResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", err
	}
	if len(out.Choices) == 0 {
		return "", fmt.Errorf("resposta sem choices")
	}
	return out.Choices[0].Message.Content, nil
}

func persistState(state RunState) {
	file, err := os.Create("hitl_run_state.json")
	if err != nil {
		fmt.Printf("Falha ao persistir estado: %v\n", err)
		return
	}
	defer file.Close()

	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")
	if err := enc.Encode(state); err != nil {
		fmt.Printf("Falha ao serializar estado: %v\n", err)
	}
}
