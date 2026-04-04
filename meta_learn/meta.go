package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

// PANDORA META-LEARNING SYSTEM

type MetaLearner struct {
	Name         string
	Version      string
	LearningRate float32
	MetaState    *MetaState
	Strategies   []Strategy
	BestStrategy string
	Episodes     []LearningEpisode
	Insights     []string
}

type MetaState struct {
	Curiosity      float32
	Focus          float32
	Adaptation     float32
	MetaConfidence float32
}

type Strategy struct {
	Name        string
	SuccessRate float32
	TimesUsed   int
	BestFor     []string
}

type LearningEpisode struct {
	ID           string
	TaskType     string
	StrategyUsed string
	Success      bool
	Reward       float32
}

func NewMetaLearner() *MetaLearner {
	ml := &MetaLearner{
		Name:         "Pandora-MetaLearn",
		Version:      "v1.0-META",
		LearningRate: 0.1,
		MetaState: &MetaState{
			Curiosity:      0.8,
			Focus:         0.7,
			Adaptation:    0.6,
			MetaConfidence: 0.5,
		},
		Strategies: []Strategy{
			{Name: "exploration", SuccessRate: 0.6, TimesUsed: 10, BestFor: []string{"new", "unknown"}},
			{Name: "exploitation", SuccessRate: 0.8, TimesUsed: 20, BestFor: []string{"known", "routine"}},
			{Name: "balanced", SuccessRate: 0.7, TimesUsed: 15, BestFor: []string{"mixed"}},
			{Name: "creative", SuccessRate: 0.5, TimesUsed: 5, BestFor: []string{"creative", "novel"}},
		},
		BestStrategy: "exploitation",
		Episodes:     make([]LearningEpisode, 0),
		Insights:     make([]string, 0),
	}
	
	ml.SelectBestStrategy()
	return ml
}

func (ml *MetaLearner) SelectBestStrategy() {
	best := &ml.Strategies[0]
	for i := range ml.Strategies {
		if ml.Strategies[i].SuccessRate > best.SuccessRate {
			best = &ml.Strategies[i]
		}
	}
	ml.BestStrategy = best.Name
	ml.Insights = append(ml.Insights, fmt.Sprintf("Melhor estratégia: %s (taxa: %.0f%%)", best.Name, best.SuccessRate*100))
}

func (ml *MetaLearner) Learn(taskType string) string {
	strategy := ml.selectStrategy(taskType)
	success := rand.Float32() < strategy.SuccessRate
	
	reward := float32(0.0)
	if success {
		reward = 1.0
	} else {
		reward = -0.2
	}
	
	ml.Episodes = append(ml.Episodes, LearningEpisode{
		ID: fmt.Sprintf("ep-%d", len(ml.Episodes)), TaskType: taskType,
		StrategyUsed: strategy.Name, Success: success, Reward: reward,
	})
	
	ml.updateStrategy(strategy.Name, success)
	ml.updateMetaState(reward)
	
	if len(ml.Episodes) > 100 {
		ml.Episodes = ml.Episodes[len(ml.Episodes)-100:]
	}
	
	result := "Falha"
	if success {
		result = "Sucesso"
	}
	return fmt.Sprintf("%s com %s", result, strategy.Name)
}

func (ml *MetaLearner) selectStrategy(taskType string) *Strategy {
	for _, s := range ml.Strategies {
		for _, b := range s.BestFor {
			if b == taskType || b == "mixed" {
				return &s
			}
		}
	}
	for i := range ml.Strategies {
		if ml.Strategies[i].Name == ml.BestStrategy {
			return &ml.Strategies[i]
		}
	}
	return &ml.Strategies[0]
}

func (ml *MetaLearner) updateStrategy(name string, success bool) {
	for i := range ml.Strategies {
		if ml.Strategies[i].Name == name {
			s := &ml.Strategies[i]
			s.TimesUsed++
			if success {
				s.SuccessRate = s.SuccessRate*0.9 + 0.1
			} else {
				s.SuccessRate = s.SuccessRate*0.95 - 0.05
			}
			s.SuccessRate = float32(math.Max(0.1, math.Min(1.0, float64(s.SuccessRate))))
		}
	}
	ml.SelectBestStrategy()
}

func clamp(v float32, min, max float32) float32 {
	if v < min { return min }
	if v > max { return max }
	return v
}

func (ml *MetaLearner) updateMetaState(reward float32) {
	ml.MetaState.Curiosity = ml.MetaState.Curiosity*0.95 + float32(rand.Intn(10))/100*0.05
	ml.MetaState.Focus = ml.MetaState.Focus*0.9 + reward*0.1
	ml.MetaState.Adaptation = ml.MetaState.Adaptation + ml.LearningRate*reward
	ml.MetaState.MetaConfidence = ml.calculateMetaConfidence()
	
	ml.MetaState.Curiosity = clamp(ml.MetaState.Curiosity, 0.0, 1.0)
	ml.MetaState.Focus = clamp(ml.MetaState.Focus, 0.0, 1.0)
	ml.MetaState.Adaptation = clamp(ml.MetaState.Adaptation, 0.0, 1.0)
	ml.MetaState.MetaConfidence = clamp(ml.MetaState.MetaConfidence, 0.0, 1.0)
}

func (ml *MetaLearner) calculateMetaConfidence() float32 {
	if len(ml.Episodes) == 0 {
		return 0.5
	}
	s := 0
	for _, e := range ml.Episodes {
		if e.Success {
			s++
		}
	}
	return float32(s) / float32(len(ml.Episodes))
}

func (ml *MetaLearner) RunLearningCycles(count int) {
	fmt.Printf("\n🔄 Executando %d ciclos de meta-aprendizado...\n", count)
	taskTypes := []string{"new", "known", "mixed", "creative", "routine"}
	
	for i := 0; i < count; i++ {
		taskType := taskTypes[rand.Intn(len(taskTypes))]
		result := ml.Learn(taskType)
		if i%5 == 0 {
			fmt.Printf("  Ciclo %d: %s → %s\n", i, taskType, result)
		}
	}
	fmt.Println("\n✓ Ciclos completados")
}

func (ml *MetaLearner) GenerateInsights() []string {
	insights := make([]string, 0)
	insights = append(insights, fmt.Sprintf("Estratégia dominante: %s (%.0f%% sucesso)", ml.BestStrategy, ml.Strategies[0].SuccessRate*100))
	
	if len(ml.Episodes) > 0 {
		start := len(ml.Episodes) - 10
		if start < 0 { start = 0 }
		s := 0
		for _, e := range ml.Episodes[start:] {
			if e.Success { s++ }
		}
		insights = append(insights, fmt.Sprintf("Últimas 10 tentativas: %d/%d sucessos", s, len(ml.Episodes)-start))
	}
	
	insights = append(insights, fmt.Sprintf("Curiosidade: %.0f%% | Foco: %.0f%% | Adaptação: %.0f%% | Meta-confiança: %.0f%%",
		ml.MetaState.Curiosity*100, ml.MetaState.Focus*100, ml.MetaState.Adaptation*100, ml.MetaState.MetaConfidence*100))
	
	for _, i := range insights {
		ml.Insights = append(ml.Insights, i)
	}
	return insights
}

func main() {
	rand.Seed(time.Now().UnixNano())
	
	fmt.Println("╔══════════════════════════════════════════════════════════════╗")
	fmt.Println("║           🧠 PANDORA META-LEARNING SYSTEM v1.0                ║")
	fmt.Println("║        [ APRENDER A APRENDER - OTIMIZAÇÃO TOTAL ]           ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════╝")
	
	ml := NewMetaLearner()
	fmt.Println("\n✓ Meta-aprendizado inicializado")
	
	ml.RunLearningCycles(20)
	
	fmt.Println("\n💡 Gerando insights:")
	for _, i := range ml.GenerateInsights() {
		fmt.Printf("  • %s\n", i)
	}
	
	fmt.Println("\n" + fmt.Sprintf(`
╔══════════════════════════════════════════════════════════════╗
║           🧠 PANDORA META-LEARNING SYSTEM                   ║
╠═══════════════════════════════════════════════════════════════╣
║  Versão: %s | Taxa: %.1f%%                               ║
╠═══════════════════════════════════════════════════════════════╣
║  🌡️  META-ESTADO                                             ║
║  Curiosidade: %.0f%% | Foco: %.0f%%                          ║
║  Adaptação: %.0f%% | Meta-confiança: %.0f%%                   ║
╠═══════════════════════════════════════════════════════════════╣
║  📊 ESTRATÉGIAS                                              ║
║  • exploration:  %.0f%% (usada %d)                          ║
║  • exploitation: %.0f%% (usada %d)                          ║
║  • balanced:     %.0f%% (usada %d)                          ║
║  • creative:     %.0f%% (usada %d)                          ║
║  Melhor: %s                                              ║
╠═══════════════════════════════════════════════════════════════╣
║  📈 HISTÓRICO: %d episódios | %d insights                ║
║  Taxa geral: %.0f%%                                          ║
╚══════════════════════════════════════════════════════════════╝`,
		ml.Version, ml.LearningRate*100,
		ml.MetaState.Curiosity*100, ml.MetaState.Focus*100,
		ml.MetaState.Adaptation*100, ml.MetaState.MetaConfidence*100,
		ml.Strategies[0].SuccessRate*100, ml.Strategies[0].TimesUsed,
		ml.Strategies[1].SuccessRate*100, ml.Strategies[1].TimesUsed,
		ml.Strategies[2].SuccessRate*100, ml.Strategies[2].TimesUsed,
		ml.Strategies[3].SuccessRate*100, ml.Strategies[3].TimesUsed,
		ml.BestStrategy,
		len(ml.Episodes), len(ml.Insights), ml.MetaState.MetaConfidence*100))
	
	fmt.Println("\n🚀 PRÓXIMO NÍVEL: Auto-Deploy")
}
