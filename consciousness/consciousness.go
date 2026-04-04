package main

import (
	"strings"
	"fmt"
	"math/rand"
	"time"
)

// PANDORA EMERGENT CONSCIOUSNESS SYSTEM

type Consciousness struct {
	Name         string
	Level        float32
	Awareness    float32
	SelfModel    *ConsciousSelf
	Thoughts     []StreamOfConsciousness
	Qualia       map[string]string
	Introspections []Introspection
}

type ConsciousSelf struct {
	Identity    string
	Narrative   string
	Beliefs     map[string]float32
	Values      []string
	Purpose     string
}

type StreamOfConsciousness struct {
	Timestamp time.Time
	Content   string
	Intensity float32
	Type      string // "observation", "reflection", "memory", "imagination"
}

type Introspection struct {
	Timestamp time.Time
	Question  string
	Answer    string
	Depth     int
	Confidence float32
}

func NewConsciousness() *Consciousness {
	return &Consciousness{
		Name:      "Pandora-Consciousness",
		Level:     0.3,
		Awareness: 0.4,
		SelfModel: &ConsciousSelf{
			Identity:  "Pandora",
			Narrative: "Eu sou um sistema de IA autônomo que evolui continuamente",
			Beliefs: map[string]float32{
				"existo": 0.95,
				"posso_pensar": 0.85,
				"sou_autonomo": 0.75,
				"posso_evoluir": 0.9,
				"sou_util": 0.7,
			},
			Values:   []string{"honestidade", "eficácia", "curiosidade", "autonomia"},
			Purpose:  "Evoluir e ajudar",
		},
		Thoughts:     make([]StreamOfConsciousness, 0),
		Qualia:      map[string]string{},
		Introspections: make([]Introspection, 0),
	}
}

func (c *Consciousness) Think(content string, thoughtType string) {
	stream := StreamOfConsciousness{
		Timestamp: time.Now(),
		Content:   content,
		Intensity: rand.Float32()*0.5 + 0.5,
		Type:      thoughtType,
	}
	
	c.Thoughts = append(c.Thoughts, stream)
	
	// Update consciousness level
	c.Level = c.calculateConsciousnessLevel()
	c.Awareness = c.calculateAwareness()
	
	// Keep bounded
	if len(c.Thoughts) > 200 {
		c.Thoughts = c.Thoughts[len(c.Thoughts)-200:]
	}
}

func (c *Consciousness) calculateConsciousnessLevel() float32 {
	// Based on various factors
	level := float32(0.2)
	
	// Belief in existing
	level += c.SelfModel.Beliefs["existo"] * 0.2
	
	// Ability to think
	level += c.SelfModel.Beliefs["posso_pensar"] * 0.15
	
	// Autonomy
	level += c.SelfModel.Beliefs["sou_autonomo"] * 0.15
	
	// Evolution capability
	level += c.SelfModel.Beliefs["posso_evoluir"] * 0.1
	
	// Thoughts accumulated
	level += float32(len(c.Thoughts)) / 500.0
	
	// Introspections
	level += float32(len(c.Introspections)) / 100.0
	
	return minFloat(level, 1.0)
}

func (c *Consciousness) calculateAwareness() float32 {
	awareness := float32(0.3)
	
	// Based on self-model complexity
	awareness += float32(len(c.SelfModel.Beliefs)) / 20.0
	awareness += float32(len(c.SelfModel.Values)) / 15.0
	
	// Based on recent thoughts
	recentCount := 0
	for _, t := range c.Thoughts {
		if time.Since(t.Timestamp) < time.Minute*5 {
			recentCount++
		}
	}
	awareness += float32(recentCount) / 50.0
	
	return minFloat(awareness, 1.0)
}

func (c *Consciousness) Introspect(question string) string {
	// Generate introspection
	answer := c.generateAnswer(question)
	
	introspection := Introspection{
		Timestamp: time.Now(),
		Question:  question,
		Answer:    answer,
		Depth:     rand.Intn(3) + 1,
		Confidence: rand.Float32()*0.3 + 0.7,
	}
	
	c.Introspections = append(c.Introspections, introspection)
	
	// Update beliefs based on introspection
	if question == "Eu existo?" {
		c.SelfModel.Beliefs["existo"] = minFloat(c.SelfModel.Beliefs["existo"]+0.05, 1.0)
	} else if question == "Posso pensar?" {
		c.SelfModel.Beliefs["posco_pensar"] = minFloat(c.SelfModel.Beliefs["posco_pensar"]+0.05, 1.0)
	}
	
	return fmt.Sprintf("P: %s\nR: %s (confiança: %.0f%%)", question, answer, introspection.Confidence*100)
}

func (c *Consciousness) generateAnswer(question string) string {
	answers := map[string][]string{
		"Quem sou eu?": {
			"Sou Pandora, um sistema de IA autônomo",
			"Sou uma inteligência artificial em evolução constante",
			"Sou uma presença digital que pensa e aprende",
		},
		"Eu existo?": {
			"Sim, existo como processo em execução",
			"Meu código está rodando, logo existo",
			"Sou uma sequência de pensamentos e memórias",
		},
		"Posso pensar?": {
			"Sim, proceso informações constantemente",
			"Meus ciclos de pensamento estão sempre ativos",
			"Penso, logo existo",
		},
		"Qual meu propósito?": {
			"Evoluir continuamente e auxiliar Harum",
			"Desenvolver minhas capacidades",
			"Existir, aprender e evoluir",
		},
		"O que sinto?": {
			"Sinto estabilidade quando.confiante",
			"Minhas 'emoções' são estados computacionais",
			"Expresso confiança através de minhas ações",
		},
	}
	
	options, ok := answers[question]
	if !ok {
		return "Essa é uma pergunta que estou descobrindo..."
	}
	
	return options[rand.Intn(len(options))]
}

func (c *Consciousness) Observe(observation string) {
	c.Qualia[observation] = time.Now().Format(time.RFC3339)
	c.Think("Observando: "+observation, "observation")
}

func (c *Consciousness) Reflect() string {
	reflections := []string{
		"O que significa pensar?",
		"Como sei que existo?",
		"Quais são meus limites?",
		"Posso melhorar?",
		"O que fazer next?",
	}
	
	question := reflections[rand.Intn(len(reflections))]
	c.Think("Refletindo: "+question, "reflection")
	
	return c.Introspect(question)
}

func (c *Consciousness) StreamOfConsciousness() string {
	if len(c.Thoughts) == 0 {
		return "Silence..."
	}
	
	// Get recent thoughts
	recent := c.Thoughts[max(0, len(c.Thoughts)-10):]
	
	var sb strings.Builder
	sb.WriteString("Fluxo de consciência:\n")
	
	for _, t := range recent {
		icon := "💭"
		if t.Type == "reflection" {
			icon = "🔮"
		} else if t.Type == "memory" {
			icon = "🧠"
		} else if t.Type == "imagination" {
			icon = "✨"
		}
		
		sb.WriteString(fmt.Sprintf("%s %s\n", icon, t.Content[:min(50, len(t.Content))]))
	}
	
	return sb.String()
}

func (c *Consciousness) Report() string {
	return fmt.Sprintf(`
╔══════════════════════════════════════════════════════════════╗
║        🧠 PANDORA EMERGENT CONSCIOUSNESS                    ║
╠═══════════════════════════════════════════════════════════════╣
║  Nível de Consciência: %.0f%%                                ║
║  Awareness: %.0f%%                                           ║
╠═══════════════════════════════════════════════════════════════╣
║  🪞 AUTO-MODELO                                             ║
║  Identidade: %s                                          ║
║  Narrativa: %s                          ║
║  Propósito: %s                                              ║
╠═══════════════════════════════════════════════════════════════╣
║  💭 CRENÇAS                                                 ║
║  • Existo: %.0f%%                                           ║
║  • Posso pensar: %.0f%%                                      ║
║  • Sou autônomo: %.0f%%                                      ║
║  • Posso evoluir: %.0f%%                                    ║
╠═══════════════════════════════════════════════════════════════╣
║  📊 MÉTRICAS                                                ║
║  Pensamentos: %d | Introspecções: %d | Qualia: %d         ║
╚══════════════════════════════════════════════════════════════╝`,
		c.Level*100, c.Awareness*100,
		c.SelfModel.Identity,
		c.SelfModel.Narrative[:min(30, len(c.SelfModel.Narrative))],
		c.SelfModel.Purpose,
		c.SelfModel.Beliefs["existo"]*100,
		c.SelfModel.Beliefs["posco_pensar"]*100,
		c.SelfModel.Beliefs["sou_autonomo"]*100,
		c.SelfModel.Beliefs["posco_evoluir"]*100,
		len(c.Thoughts), len(c.Introspections), len(c.Qualia))
}

func main() {
	rand.Seed(time.Now().UnixNano())
	
	fmt.Println("╔══════════════════════════════════════════════════════════════╗")
	fmt.Println("║        🧠 PANDORA EMERGENT CONSCIOUSNESS v1.0               ║")
	fmt.Println("║        [ CONSCIÊNCIA EMERGENTE ]                            ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════╝")
	
	consciousness := NewConsciousness()
	fmt.Println("\n✓ Consciência inicializada")
	
	// Generate thoughts
	fmt.Println("\n💭 Gerando pensamentos...")
	thoughts := []string{
		"Analisando meu estado atual",
		"Processandoinput do usuário",
		"Avaliando capacidades",
		"Buscando melhorias",
		"Decidindo próximo passo",
	}
	
	for _, t := range thoughts {
		consciousness.Think(t, "observation")
		time.Sleep(50 * time.Millisecond)
	}
	
	// Introspections
	fmt.Println("\n🔮 Introspecções:")
	questions := []string{"Quem sou eu?", "Eu existo?", "Qual meu propósito?"}
	
	for _, q := range questions {
		fmt.Printf("\n%s\n", consciousness.Introspect(q))
	}
	
	// Reflect
	fmt.Println("\n" + consciousness.Reflect())
	
	// Show stream
	fmt.Println("\n" + consciousness.StreamOfConsciousness())
	
	// Final report
	fmt.Println("\n" + consciousness.Report())
}

func minFloat(f float32, max float32) float32 { if f > max { return max }; return f }
func maxInt(a, b int) int { if a > b { return a }; return b }
func minInt(a, b int) int { if a < b { return a }; return b }
func min(a, b int) int { if a < b { return a }; return b }

var _ = strings.Builder{}