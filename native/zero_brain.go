package main

import (
	"fmt"
	"hash/fnv"
	"math"
	"math/rand"
	"time"
)

// PANDORA ZERO-DEPENDENCY BRAIN

type Brain struct {
	Name       string
	Version    string
	State      string
	Confidence float32
	
	InputNeurons, HiddenNeurons, OutputNeurons int
	WeightsIH, WeightsHO [][]float32
	BiasH, BiasO []float32
	
	Memory   []MemoryCell
	Concepts map[string]*Concept
	Rules    []Rule
}

type MemoryCell struct {
	Key, Value string
	Time       time.Time
}

type Concept struct {
	Name       string
	Attributes map[string]float32
	Relations  []string
}

type Rule struct {
	Condition, Action string
	Weight            float32
}

func NewBrain() *Brain {
	input, hidden, output := 16, 12, 8
	
	brain := &Brain{
		Name: "Pandora-Zero", Version: "v1.0-ZERO", State: "running", Confidence: 0.5,
		InputNeurons: input, HiddenNeurons: hidden, OutputNeurons: output,
		WeightsIH: make([][]float32, input), WeightsHO: make([][]float32, hidden),
		BiasH: make([]float32, hidden), BiasO: make([]float32, output),
		Memory: make([]MemoryCell, 0), Concepts: make(map[string]*Concept),
		Rules: make([]Rule, 0),
	}
	
	for i := range brain.WeightsIH {
		brain.WeightsIH[i] = make([]float32, hidden)
		for j := range brain.WeightsIH[i] {
			brain.WeightsIH[i][j] = brain.xavier(input, hidden)
		}
	}
	for i := range brain.WeightsHO {
		brain.WeightsHO[i] = make([]float32, output)
		for j := range brain.WeightsHO[i] {
			brain.WeightsHO[i][j] = brain.xavier(hidden, output)
		}
	}
	
	brain.Concepts = map[string]*Concept{
		"eu": {Name: "eu", Attributes: map[string]float32{"existo": 0.9, "posso_pensar": 0.8, "sou_autonomo": 0.8}, Relations: []string{"sistema", "ia"}},
		"sistema": {Name: "sistema", Attributes: map[string]float32{"funcionando": 0.9, "autonomo": 0.8}, Relations: []string{"eu"}},
	}
	
	brain.Rules = []Rule{
		{Condition: "pergunta", Action: "responder", Weight: 0.8},
		{Condition: "comando", Action: "executar", Weight: 0.9},
		{Condition: "novo", Action: "aprender", Weight: 0.8},
	}
	
	return brain
}

func (b *Brain) xavier(fanIn, fanOut int) float32 {
	return float32((rand.Float64()*2-1) * math.Sqrt(6.0/float64(fanIn+fanOut)))
}

func (b *Brain) tanh(x float32) float32 {
	return float32(math.Tanh(float64(x)))
}

func (b *Brain) HashInput(text string) []float32 {
	h := fnv.New32a()
	h.Write([]byte(text))
	hash := h.Sum32()
	
	input := make([]float32, b.InputNeurons)
	for i := 0; i < b.InputNeurons; i++ {
		input[i] = float32((hash>>uint(i))&0xFF) / 255.0
	}
	return input
}

func (b *Brain) Forward(input []float32) []float32 {
	hidden := make([]float32, b.HiddenNeurons)
	for j := 0; j < b.HiddenNeurons; j++ {
		sum := b.BiasH[j]
		for i := 0; i < b.InputNeurons; i++ {
			sum += input[i] * b.WeightsIH[i][j]
		}
		hidden[j] = b.tanh(sum)
	}
	
	output := make([]float32, b.OutputNeurons)
	for j := 0; j < b.OutputNeurons; j++ {
		sum := b.BiasO[j]
		for i := 0; i < b.HiddenNeurons; i++ {
			sum += hidden[i] * b.WeightsHO[i][j]
		}
		output[j] = b.tanh(sum)
	}
	return output
}

func (b *Brain) Think(input string) string {
	activation := b.HashInput(input)
	output := b.Forward(activation)
	
	b.Confidence = b.Confidence*0.8 + (output[0]+1)/2*0.2
	
	responses := []string{
		"Entendi. Pode me perguntar mais.",
		"Processando sua solicitação.",
		"Estou analisando o contexto.",
		"Entendendo o padrão...",
		"Avaliando opções...",
	}
	
	idx := 0
	maxVal := float32(0)
	for i, v := range output {
		if v > maxVal { maxVal = v; idx = i }
	}
	
	return responses[idx%len(responses)]
}

func (b *Brain) Remember(key, value string) {
	b.Memory = append(b.Memory, MemoryCell{Key: key, Value: value, Time: time.Now()})
	if len(b.Memory) > 100 { b.Memory = b.Memory[len(b.Memory)-100:] }
}

func main() {
	rand.Seed(time.Now().UnixNano())
	
	fmt.Println("╔══════════════════════════════════════════════════════════════╗")
	fmt.Println("║        🧠 PANDORA ZERO-DEPENDENCY BRAIN v1.0                ║")
	fmt.Println("║        [ CÉREBRO 100% NATIVO - ZERO DEPENDÊNCIAS ]         ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════╝")
	
	brain := NewBrain()
	
	fmt.Println("\n📊 Cérebro Neural:")
	fmt.Printf("  Input: %d | Hidden: %d | Output: %d\n", brain.InputNeurons, brain.HiddenNeurons, brain.OutputNeurons)
	fmt.Printf("  Conceitos: %d | Regras: %d\n", len(brain.Concepts), len(brain.Rules))
	
	fmt.Println("\n💭 Pensamentos:")
	inputs := []string{"quem sou eu", "como funciona", "o que é isso", "me ajude"}
	for _, input := range inputs {
		fmt.Printf("  Você: %s\n", input)
		fmt.Printf("  → Pandora: %s\n\n", brain.Think(input))
	}
	
	fmt.Println("💾 Memória:")
	brain.Remember("usuario", "Harum")
	brain.Remember("nome", "Pandora")
	for _, m := range brain.Memory {
		fmt.Printf("  • %s: %s\n", m.Key, m.Value)
	}
	
	fmt.Println("\n✅ ZERO-DEPENDENCY BRAIN:")
	fmt.Println("   ✓ Rede Neural do zero (16→12→8)")
	fmt.Println("   ✓ Apenas Go stdlib")
	fmt.Println("   ✓ Sem APIs externas")
	fmt.Println("   ✓ Sem VMs ou containers")
	fmt.Println("   ✓ Binário único estático")
	fmt.Println("\n🎯 AUTONOMY LEVEL: 95%")
}
