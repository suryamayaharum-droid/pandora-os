package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

type AutonomousCycle struct {
	Cycle       int       `json:"cycle"`
	Goal        string    `json:"goal"`
	Observation string    `json:"observation"`
	Plan        string    `json:"plan"`
	Action      string    `json:"action"`
	Result      string    `json:"result"`
	Timestamp   time.Time `json:"timestamp"`
}

func runAutonomous(goal string, iterations int) {
	if iterations < 1 {
		iterations = 1
	}

	fmt.Println("🤖 Modo autônomo iniciado")
	fmt.Printf("🎯 Objetivo: %s\n", goal)
	fmt.Printf("🔁 Ciclos: %d\n\n", iterations)

	for i := 1; i <= iterations; i++ {
		observation := autonomousObservation(i)
		plan := buildPlan(fmt.Sprintf("%s | contexto: %s", goal, observation))
		action := deriveAction(plan)
		result := executePlan(plan)

		cycle := AutonomousCycle{
			Cycle:       i,
			Goal:        goal,
			Observation: observation,
			Plan:        plan,
			Action:      action,
			Result:      result,
			Timestamp:   time.Now(),
		}

		appendCycle(cycle)
		fmt.Printf("[ciclo %d] ação=%s resultado=%s\n", i, action, result)
		time.Sleep(500 * time.Millisecond)
	}

	fmt.Println("\n✅ Execução autônoma finalizada. Registro salvo em autonomous_cycles.jsonl")
}

func autonomousObservation(cycle int) string {
	return fmt.Sprintf("ciclo=%d; estado=operacional; necessidade=progresso incremental", cycle)
}

func deriveAction(plan string) string {
	firstLine := strings.Split(strings.TrimSpace(plan), "\n")[0]
	if firstLine == "" {
		return "executar_proximo_passo"
	}
	return strings.ToLower(strings.ReplaceAll(firstLine, " ", "_"))
}

func appendCycle(c AutonomousCycle) {
	f, err := os.OpenFile("autonomous_cycles.jsonl", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Printf("falha ao salvar ciclo: %v\n", err)
		return
	}
	defer f.Close()

	b, err := json.Marshal(c)
	if err != nil {
		fmt.Printf("falha ao serializar ciclo: %v\n", err)
		return
	}
	f.Write(append(b, '\n'))
}
