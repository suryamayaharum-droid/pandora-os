package main

import (
	"crypto/sha256"
	"fmt"
	"math"
	"math/rand"
	"time"
)

// PANDORA NATIVE AI ENGINE - 100% Local

type NeuralNet struct {
	Config   *NetConfig
	Weights  [][][]float32
	Biases   [][]float32
}

type NetConfig struct {
	InputSize, HiddenSize, OutputSize int
	LearningRate float32
}

func NewNeuralNet(input, hidden, output int) *NeuralNet {
	net := &NeuralNet{
		Config: &NetConfig{InputSize: input, HiddenSize: hidden, OutputSize: output, LearningRate: 0.01},
	}
	net.Weights = [][][]float32{randomMatrix(input, hidden), randomMatrix(hidden, output)}
	net.Biases = [][]float32{randomVector(hidden), randomVector(output)}
	return net
}

func randomMatrix(rows, cols int) [][]float32 {
	m := make([][]float32, rows)
	for i := range m {
		m[i] = make([]float32, cols)
		for j := range m[i] {
			m[i][j] = (rand.Float32()*2 - 1) * 0.1
		}
	}
	return m
}

func randomVector(size int) []float32 {
	v := make([]float32, size)
	for i := range v {
		v[i] = (rand.Float32()*2 - 1) * 0.1
	}
	return v
}

func sigmoid(x float32) float32 {
	return float32(1.0 / (1.0 + math.Exp(-float64(x))))
}

func sigmoidDeriv(x float32) float32 {
	return x * (1 - x)
}

func (net *NeuralNet) Forward(input []float32) []float32 {
	hidden := make([]float32, net.Config.HiddenSize)
	for j := 0; j < net.Config.HiddenSize; j++ {
		sum := net.Biases[0][j]
		for i := 0; i < net.Config.InputSize; i++ {
			sum += input[i] * net.Weights[0][i][j]
		}
		hidden[j] = sigmoid(sum)
	}
	
	output := make([]float32, net.Config.OutputSize)
	for j := 0; j < net.Config.OutputSize; j++ {
		sum := net.Biases[1][j]
		for i := 0; i < net.Config.HiddenSize; i++ {
			sum += hidden[i] * net.Weights[1][i][j]
		}
		output[j] = sigmoid(sum)
	}
	return output
}

func (net *NeuralNet) Train(input, target []float32) {
	hidden := make([]float32, net.Config.HiddenSize)
	for j := 0; j < net.Config.HiddenSize; j++ {
		sum := net.Biases[0][j]
		for i := 0; i < net.Config.InputSize; i++ {
			sum += input[i] * net.Weights[0][i][j]
		}
		hidden[j] = sigmoid(sum)
	}
	
	output := make([]float32, net.Config.OutputSize)
	for j := 0; j < net.Config.OutputSize; j++ {
		sum := net.Biases[1][j]
		for i := 0; i < net.Config.HiddenSize; i++ {
			sum += hidden[i] * net.Weights[1][i][j]
		}
		output[j] = sigmoid(sum)
	}
	
	outputError := make([]float32, net.Config.OutputSize)
	for i := 0; i < net.Config.OutputSize; i++ {
		outputError[i] = (target[i] - output[i]) * sigmoidDeriv(output[i])
	}
	
	hiddenError := make([]float32, net.Config.HiddenSize)
	for i := 0; i < net.Config.HiddenSize; i++ {
		sum := float32(0)
		for j := 0; j < net.Config.OutputSize; j++ {
			sum += outputError[j] * net.Weights[1][i][j]
		}
		hiddenError[i] = sum * sigmoidDeriv(hidden[i])
	}
	
	lr := net.Config.LearningRate
	for i := 0; i < net.Config.HiddenSize; i++ {
		for j := 0; j < net.Config.OutputSize; j++ {
			net.Weights[1][i][j] += lr * outputError[j] * hidden[i]
		}
	}
	for i := 0; i < net.Config.InputSize; i++ {
		for j := 0; j < net.Config.HiddenSize; j++ {
			net.Weights[0][i][j] += lr * hiddenError[j] * input[i]
		}
	}
	for i := 0; i < net.Config.OutputSize; i++ {
		net.Biases[1][i] += lr * outputError[i]
	}
	for i := 0; i < net.Config.HiddenSize; i++ {
		net.Biases[0][i] += lr * hiddenError[i]
	}
}

func (net *NeuralNet) EmbedText(text string) []float32 {
	hash := sha256.Sum256([]byte(text))
	embedding := make([]float32, net.Config.InputSize)
	for i := 0; i < net.Config.InputSize; i++ {
		embedding[i] = float32(hash[i%len(hash)]) / 255.0
	}
	sum := float32(0)
	for _, v := range embedding {
		sum += v * v
	}
	sum = float32(math.Sqrt(float64(sum)))
	if sum > 0 {
		for i := range embedding {
			embedding[i] /= sum
		}
	}
	return embedding
}

type TrainingPair struct {
	Input, Output []float32
}

func main() {
	rand.Seed(time.Now().UnixNano())
	
	fmt.Println("╔══════════════════════════════════════════════════════════════╗")
	fmt.Println("║        🧠 PANDORA NATIVE AI ENGINE v1.0                     ║")
	fmt.Println("║        [ 100% LOCAL - SEM APIs, SEM VMs ]                   ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════╝")
	
	net := NewNeuralNet(10, 8, 3)
	fmt.Println("\n📊 Rede Neural:")
	fmt.Printf("  Input: %d | Hidden: %d | Output: %d\n", net.Config.InputSize, net.Config.HiddenSize, net.Config.OutputSize)
	
	data := []TrainingPair{
		{[]float32{1,0,0,0,0,0,0,0,0,0}, []float32{1,0,0}},
		{[]float32{0,1,0,0,0,0,0,0,0,0}, []float32{1,0,0}},
		{[]float32{0,0,1,0,0,0,0,0,0,0}, []float32{1,0,0}},
		{[]float32{0,0,0,1,0,0,0,0,0,0}, []float32{0,1,0}},
		{[]float32{0,0,0,0,1,0,0,0,0,0}, []float32{0,1,0}},
		{[]float32{0,0,0,0,0,1,0,0,0,0}, []float32{0,1,0}},
		{[]float32{0,0,0,0,0,0,1,0,0,0}, []float32{0,0,1}},
		{[]float32{0,0,0,0,0,0,0,1,0,0}, []float32{0,0,1}},
		{[]float32{0,0,0,0,0,0,0,0,1,0}, []float32{0,0,1}},
		{[]float32{0,0,0,0,0,0,0,0,0,1}, []float32{0,0,1}},
	}
	
	fmt.Println("\n📚 Treinando...")
	for epoch := 0; epoch < 1000; epoch++ {
		for _, p := range data {
			net.Train(p.Input, p.Output)
		}
	}
	fmt.Println("  ✓ Treinamento completo!")
	
	fmt.Println("\n🧪 Testando:")
	labels := []string{"greeting", "question", "command"}
	for i := 0; i < 9; i++ {
		output := net.Forward(data[i].Input)
		best := 0
		for j, s := range output {
			if s > output[best] {
				best = j
			}
		}
		fmt.Printf("  Input #%d → %s (%.0f%%)\n", i, labels[best], output[best]*100)
	}
	
	fmt.Println("\n📝 Embedding:")
	emb := net.EmbedText("Pandora AI")
	fmt.Printf("  \"Pandora AI\" → [%d dimensões]\n", len(emb))
	fmt.Printf("  Primeiros 5: [%.2f, %.2f, %.2f, %.2f, %.2f]\n", emb[0], emb[1], emb[2], emb[3], emb[4])
	
	fmt.Println("\n✅ NATIVE AI - 100% LOCAL")
	fmt.Println("   ✓ Sem APIs externas")
	fmt.Println("   ✓ Sem VMs ou containers")
	fmt.Println("   ✓ Rede neural do zero")
	fmt.Println("   ✓ Treino supervisionado")
	fmt.Println("   ✓ Embedding de texto")
}
