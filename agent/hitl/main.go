package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
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

type apiStatus struct {
	Service         string    `json:"service"`
	Version         string    `json:"version"`
	Time            time.Time `json:"time"`
	ModelConfigured bool      `json:"model_configured"`
}

type apiRunRequest struct {
	Goal     string `json:"goal"`
	Approved bool   `json:"approved"`
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

type pageData struct {
	Result *RunState
	Error  string
}

func main() {
	webMode := flag.Bool("web", false, "inicia interface HTML")
	autoMode := flag.Bool("auto", false, "executa agente autônomo em ciclos")
	addr := flag.String("addr", ":8080", "endereço do servidor web")
	autoGoal := flag.String("goal", "Evoluir continuamente o sistema", "objetivo do modo autônomo")
	iterations := flag.Int("iterations", 5, "número de ciclos autônomos")
	flag.Parse()

	if *webMode {
		startWebServer(*addr)
		return
	}
	if *autoMode {
		runAutonomous(*autoGoal, *iterations)
		return
	}

	runCLI()
}

func runCLI() {
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

	state = processExecution(state)

	fmt.Println("\n=== Relatório final ===")
	fmt.Printf("Objetivo: %s\n", state.Goal)
	fmt.Printf("Aprovado por humano: %t\n", state.ApprovedByHuman)
	fmt.Printf("Resultado: %s\n", state.ExecutionResult)

	persistState(state)
}

func processExecution(state RunState) RunState {
	if !state.ApprovedByHuman {
		state.CurrentStage = StageReport
		state.UpdatedAt = time.Now()
		state.ExecutionResult = "Execução bloqueada: ação externa requer aprovação humana."
		persistState(state)
		return state
	}

	state.CurrentStage = StageExecute
	state.UpdatedAt = time.Now()
	state.ExecutionResult = executePlan(state.Plan)

	state.CurrentStage = StageReport
	state.UpdatedAt = time.Now()
	persistState(state)
	return state
}

func startWebServer(addr string) {
	tmpl := template.Must(template.New("page").Parse(webPage))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			tmpl.Execute(w, pageData{})
		case http.MethodPost:
			if err := r.ParseForm(); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				tmpl.Execute(w, pageData{Error: "Falha ao ler formulário."})
				return
			}

			goal := strings.TrimSpace(r.FormValue("goal"))
			approved := r.FormValue("approved") == "on"
			if goal == "" {
				w.WriteHeader(http.StatusBadRequest)
				tmpl.Execute(w, pageData{Error: "Objetivo é obrigatório."})
				return
			}

			state := runGoal(goal, approved)
			tmpl.Execute(w, pageData{Result: &state})
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/api/status", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, apiStatus{
			Service:         "pandora-hitl-agent",
			Version:         "v0.2.0",
			Time:            time.Now(),
			ModelConfigured: os.Getenv("OPENAI_COMPAT_BASE_URL") != "" && os.Getenv("OPENAI_COMPAT_MODEL") != "",
		})
	})

	http.HandleFunc("/api/run", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "use POST"})
			return
		}
		var req apiRunRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "JSON inválido"})
			return
		}
		goal := strings.TrimSpace(req.Goal)
		if goal == "" {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "goal é obrigatório"})
			return
		}
		state := runGoal(goal, req.Approved)
		writeJSON(w, http.StatusOK, state)
	})

	fmt.Printf("🌐 Interface HTML disponível em http://localhost%s\n", addr)
	fmt.Printf("📡 API status: http://localhost%s/api/status\n", addr)
	fmt.Printf("▶️  API run: POST http://localhost%s/api/run\n", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		fmt.Printf("falha ao iniciar servidor: %v\n", err)
	}
}

func runGoal(goal string, approved bool) RunState {
	state := RunState{
		Goal:            goal,
		CurrentStage:    StagePlan,
		UpdatedAt:       time.Now(),
		Plan:            buildPlan(goal),
		ApprovedByHuman: approved,
	}
	state.CurrentStage = StageReview
	return processExecution(state)
}

func writeJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
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

const webPage = `<!doctype html>
<html lang="pt-BR">
<head>
  <meta charset="UTF-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1.0" />
  <title>Agent MVP HITL</title>
  <style>
    body { font-family: Arial, sans-serif; margin: 2rem; max-width: 860px; }
    textarea, input[type=text] { width: 100%; padding: 0.6rem; margin: 0.5rem 0 1rem; }
    .card { border: 1px solid #ddd; border-radius: 8px; padding: 1rem; margin-top: 1rem; }
    button { padding: 0.7rem 1.1rem; cursor: pointer; }
    .error { color: #b00020; }
    pre { white-space: pre-wrap; }
  </style>
</head>
<body>
  <h1>Agent MVP com HITL</h1>
  <p>Preencha o objetivo e aprove para permitir execução externa.</p>
  <p><small>Endpoints: <code>/api/status</code> e <code>/api/run</code>.</small></p>

  {{if .Error}}<p class="error">{{.Error}}</p>{{end}}

  <form method="post" action="/">
    <label>Objetivo</label>
    <input type="text" name="goal" placeholder="Ex: Criar documentação de API" required />

    <label>
      <input type="checkbox" name="approved" /> Aprovar execução externa (HITL)
    </label>
    <br/><br/>
    <button type="submit">Executar</button>
  </form>

  {{if .Result}}
  <div class="card">
    <h2>Resultado</h2>
    <p><strong>Objetivo:</strong> {{.Result.Goal}}</p>
    <p><strong>Etapa final:</strong> {{.Result.CurrentStage}}</p>
    <p><strong>Aprovado por humano:</strong> {{.Result.ApprovedByHuman}}</p>
    <h3>Plano</h3>
    <pre>{{.Result.Plan}}</pre>
    <h3>Execução</h3>
    <pre>{{.Result.ExecutionResult}}</pre>
  </div>
  {{end}}
</body>
</html>`
