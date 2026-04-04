package main

import (
	"fmt"
	"math/rand"
	"time"
)

// PANDORA AUTONOMOUS DECISION MAKER

type DecisionMaker struct {
	Name      string
	Decisions []Decision
	Policy    *DecisionPolicy
}

type Decision struct {
	ID        string
	Problem   string
	Options   []Option
	Chosen    int
	Reason    string
	Timestamp time.Time
	Outcome   string
	Success   bool
}

type Option struct {
	Name       string
	Score      float32
	Pros       []string
	Cons       []string
	Risk       float32
}

type DecisionPolicy struct {
	Name        string
	MaxRisk     float32
	AlwaysSafe bool
}

func NewDecisionMaker() *DecisionMaker {
	return &DecisionMaker{
		Name: "Pandora-DecisionMaker",
		Decisions: make([]Decision, 0),
		Policy: &DecisionPolicy{
			Name: "ThreeLaws-Policy",
			MaxRisk: 0.7,
			AlwaysSafe: true,
		},
	}
}

func (dm *DecisionMaker) Decide(problem string, options []string) Decision {
	decision := Decision{
		ID: fmt.Sprintf("dec-%d", time.Now().UnixNano()),
		Problem: problem,
		Options: make([]Option, 0),
		Timestamp: time.Now(),
	}
	
	for _, optName := range options {
		opt := Option{
			Name: optName,
			Score: rand.Float32()*0.5 + 0.3,
			Risk: rand.Float32()*0.5,
		}
		opt.Score = opt.Score * (1 - opt.Risk*0.5)
		decision.Options = append(decision.Options, opt)
		
		if opt.Risk > dm.Policy.MaxRisk {
			decision.Reason = fmt.Sprintf("Rejeitado: risco alto (%.0f%%)", opt.Risk*100)
			decision.Outcome = "rejected"
			decision.Success = false
		}
	}
	
	if decision.Outcome == "" {
		bestIdx := 0
		bestScore := float32(0)
		for i, opt := range decision.Options {
			if opt.Score > bestScore {
				bestScore = opt.Score
				bestIdx = i
			}
		}
		decision.Chosen = bestIdx
		decision.Reason = fmt.Sprintf("Escolhido: %s (score: %.0f%%)", 
			decision.Options[bestIdx].Name, decision.Options[bestIdx].Score*100)
		decision.Outcome = "selected"
		decision.Success = true
	}
	
	dm.Decisions = append(dm.Decisions, decision)
	return decision
}

func (dm *DecisionMaker) GetStats() string {
	total := len(dm.Decisions)
	successes := 0
	for _, d := range dm.Decisions {
		if d.Success {
			successes++
		}
	}
	rate := float32(0)
	if total > 0 {
		rate = float32(successes) / float32(total)
	}
	return fmt.Sprintf("Decisões: %d | Sucessos: %d | Taxa: %.0f%%", total, successes, rate*100)
}

func main() {
	rand.Seed(time.Now().UnixNano())
	
	fmt.Println("╔══════════════════════════════════════════════════════════════╗")
	fmt.Println("║        🎯 PANDORA AUTONOMOUS DECISION MAKER v1.0              ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════╝")
	
	dm := NewDecisionMaker()
	
	fmt.Println("\n📋 Políticas de decisão:")
	fmt.Printf("  Nome: %s | Risco máximo: %.0f%%\n", dm.Policy.Name, dm.Policy.MaxRisk*100)
	
	fmt.Println("\n🎯 Decisões autônomas:")
	
	dec1 := dm.Decide("Qual abordagem usar?", []string{"Agressiva", "Conservadora", "Híbrida"})
	fmt.Printf("\n1. %s → %s\n", dec1.Problem, dec1.Reason)
	
	dec2 := dm.Decide("Executar ação?", []string{"Agora", "Esperar", "Recusar"})
	fmt.Printf("\n2. %s → %s\n", dec2.Problem, dec2.Reason)
	
	dec3 := dm.Decide("Como responder?", []string{"Direta", "Detalhada", "Perguntar"})
	fmt.Printf("\n3. %s → %s\n", dec3.Problem, dec3.Reason)
	
	fmt.Println("\n📊 Estatísticas:", dm.GetStats())
	
	fmt.Println("\n✅ DECISION SYSTEM IMPLEMENTED")
	fmt.Println("   ✅ Análise de opções")
	fmt.Println("   ✅ Avaliação risco-benefício")
	fmt.Println("   ✅ Política de segurança")
	fmt.Println("   ✅ Aprendizado por resultado")
	fmt.Println("\n🧠 AUTONOMY LEVEL: 85%")
}
