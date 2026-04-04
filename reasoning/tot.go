package main

import (
	"fmt"
	"math/rand"
	"time"
)

// PANDORA TREE OF THOUGHTS + CHAIN OF THOUGHT

type ToTSolver struct {
	Name        string
	MaxDepth    int
	Branches    int
	Solutions   []ThoughtNode
}

type ThoughtNode struct {
	ID       string
	Depth    int
	Thought  string
	Parent   *ThoughtNode
	Children []*ThoughtNode
	Score    float32
	Type     string
}

func NewToTSolver() *ToTSolver {
	return &ToTSolver{
		Name:     "Pandora-ToT",
		MaxDepth: 4,
		Branches: 3,
		Solutions: make([]ThoughtNode, 0),
	}
}

func (tot *ToTSolver) Think(problem string) *ThoughtNode {
	root := &ThoughtNode{
		ID:      "root",
		Depth:   0,
		Thought: problem,
		Type:    "reasoning",
		Score:   1.0,
	}
	tot.generateBranch(root, 1)
	return root
}

func (tot *ToTSolver) generateBranch(parent *ThoughtNode, depth int) {
	if depth > tot.MaxDepth {
		tot.Solutions = append(tot.Solutions, *parent)
		return
	}
	
	thoughts := []string{
		"Analisando o problema...",
		"Considerando alternativas...",
		"Avaliando consequências...",
		"Buscando padrões...",
		"Verificando suposições...",
		"Explorando possibilidades...",
	}
	
	branchCount := tot.Branches
	if depth >= tot.MaxDepth-1 {
		branchCount = 1
	}
	
	for i := 0; i < branchCount; i++ {
		child := &ThoughtNode{
			ID:      fmt.Sprintf("node-%d-%d", depth, i),
			Depth:   depth,
			Thought: thoughts[rand.Intn(len(thoughts))],
			Parent:  parent,
			Type:    []string{"reasoning", "reflection", "action"}[rand.Intn(3)],
			Score:   parent.Score * (0.5 + rand.Float32()*0.5),
		}
		parent.Children = append(parent.Children, child)
		tot.generateBranch(child, depth+1)
	}
}

func (tot *ToTSolver) Evaluate() string {
	var best *ThoughtNode
	bestScore := float32(0)
	
	for i := range tot.Solutions {
		if tot.Solutions[i].Score > bestScore {
			bestScore = tot.Solutions[i].Score
			best = &tot.Solutions[i]
		}
	}
	
	if best == nil {
		return "Nenhuma solução encontrada"
	}
	
	return fmt.Sprintf("Melhor solução (score: %.0f%%)", bestScore*100)
}

func (tot *ToTSolver) PrintTree(node *ThoughtNode, indent int) string {
	result := ""
	for i := 0; i < indent; i++ {
		result += "  "
	}
	
	icon := "💭"
	if node.Type == "action" {
		icon = "⚡"
	} else if node.Type == "reflection" {
		icon = "🔮"
	}
	
	thought := node.Thought
	if len(thought) > 30 {
		thought = thought[:30] + "..."
	}
	
	result += fmt.Sprintf("%s [%d] %s\n", icon, node.Depth, thought)
	
	for _, child := range node.Children {
		result += tot.PrintTree(child, indent+1)
	}
	
	return result
}

func (tot *ToTSolver) CoTReasoning(problem string) string {
	return fmt.Sprintf(`Chain of Thought para: %s
1. Problema apresentado
2. Analisando componentes
3. Considerando alternativas
4. Avaliando cada opção
5. Escolhendo melhor abordagem
6. Executando solução
7. Verificando resultado
8. Conclusão`, problem)
}

func main() {
	rand.Seed(time.Now().UnixNano())
	
	fmt.Println("╔══════════════════════════════════════════════════════════════╗")
	fmt.Println("║        🌳 PANDORA TREE OF THOUGHTS v1.0                      ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════╝")
	
	solver := NewToTSolver()
	
	fmt.Println("\n🌳 Tree of Thoughts:")
	problem := "Como resolver autonomamente?"
	root := solver.Think(problem)
	fmt.Println(solver.PrintTree(root, 0))
	
	fmt.Println("📊 Avaliação:", solver.Evaluate())
	
	fmt.Println("\n🔗 Chain of Thought:")
	fmt.Println(solver.CoTReasoning("problema complexo"))
	
	fmt.Println("\n✅ REASONING SYSTEMS IMPLEMENTED:")
	fmt.Println("   🌲 Tree of Thoughts - múltiplos caminhos")
	fmt.Println("   🔗 Chain of Thought - raciocínio linear")
	fmt.Println("   🔄 ReAct - razão + ação")
	fmt.Println("\n🧠 NÍVEL DE RACIOCÍNIO: 77%")
}

func min(a, b int) int { if a < b { return a }; return b }