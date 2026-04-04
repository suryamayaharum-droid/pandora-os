package main

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"strings"
	"time"
)

// ═══════════════════════════════════════════════════════════════
// PANDORA NEURAL CORE - Sistema Neural Leve e Proprietário
// Tecnologia de vetores própria + Linguagem de Máquina Universal
// ═══════════════════════════════════════════════════════════════

// === SISTEMA DE VETORES PROPRIETÁRIOS ===

// Pandora Vector - nosso formato proprietário
type PVector struct {
	ID        string
	Values    []float32
	Metadata  map[string]interface{}
	Type      VectorType
	Checksum  string
}

type VectorType int

const (
	VectorSemantic VectorType = iota
	VectorAction
	VectorMemory
	VectorGoal
	VectorCode
)

// Sistema de Encoding proprietário
type PandoraEncoding struct {
	Alphabet  string
	Base      int
	Precision int
}

const PANDORA_ALPHABET = "⏣∇◊◇□△○☆★◆◇●○◆◈▣▤▥▦▧▨▩▪▫▬▭▮▯▰▱▲△▼▽◆◇□■□●○◐◑◒◓"

func NewPandoraEncoding() *PandoraEncoding {
	return &PandoraEncoding{
		Alphabet:  PANDORA_ALPHABET,
		Base:      64,
		Precision: 8,
	}
}

func (pe *PandoraEncoding) Encode(vector []float32) string {
	var encoded strings.Builder
	
	for i, v := range vector {
		normalized := (v + 1.0) / 2.0
		index := int(normalized * float32(len(pe.Alphabet)-1))
		if index < 0 { index = 0 }
		if index >= len(pe.Alphabet) { index = len(pe.Alphabet) - 1 }
		encoded.WriteByte(pe.Alphabet[index])
		if (i+1)%8 == 0 && i < len(vector)-1 {
			encoded.WriteString("│")
		}
	}
	return encoded.String()
}

func (pe *PandoraEncoding) Decode(encoded string) []float32 {
	encoded = strings.ReplaceAll(encoded, "│", "")
	values := make([]float32, 0, len(encoded))
	
	for _, ch := range encoded {
		idx := strings.Index(pe.Alphabet, string(ch))
		if idx >= 0 {
			normalized := float32(idx) / float32(len(pe.Alphabet)-1)
			values = append(values, normalized*2.0-1.0)
		}
	}
	return values
}

// === COMPRESSÃO NEURAL PROPRIETÁRIA ===

type NeuralCompressor struct {
	Method   string
	Ratio    float32
	Quality  float32
}

func NewNeuralCompressor() *NeuralCompressor {
	return &NeuralCompressor{
		Method:  "pandora-mix",
		Ratio:   0.0,
		Quality: 0.95,
	}
}

func (nc *NeuralCompressor) Compress(weights []float32) []byte {
	n := len(weights)
	centroids := nc.kMeansClustering(weights, 8)
	
	compressed := make([]byte, 0, n)
	compressed = append(compressed, byte(len(centroids)))
	
	// Escrever centroides
	for _, c := range centroids {
		bits := math.Float32bits(c)
		compressed = append(compressed, 
			byte(bits>>24), byte(bits>>16), byte(bits>>8), byte(bits))
	}
	
	// Codificar pesos com clusters
	for _, w := range weights {
		minDist := float32(1e10)
		bestIdx := 0
		for i, c := range centroids {
			dist := abs32(w - c)
			if dist < minDist {
				minDist = dist
				bestIdx = i
			}
		}
		compressed = append(compressed, byte(bestIdx))
	}
	
	nc.Ratio = float32(len(compressed)) / float32(n*4)
	return compressed
}

func (nc *NeuralCompressor) kMeansClustering(data []float32, k int) []float32 {
	rand.Seed(time.Now().UnixNano())
	centroids := make([]float32, k)
	for i := range centroids {
		centroids[i] = data[rand.Intn(len(data))]
	}
	
	for iter := 0; iter < 5; iter++ {
		clusters := make([][]float32, k)
		for _, v := range data {
			minDist := float32(1e10)
			best := 0
			for i, c := range centroids {
				dist := abs32(v - c)
				if dist < minDist {
					minDist = dist
					best = i
				}
			}
			clusters[best] = append(clusters[best], v)
		}
		
		for i, cluster := range clusters {
			if len(cluster) > 0 {
				var sum float32
				for _, v := range cluster {
					sum += v
				}
				centroids[i] = sum / float32(len(cluster))
			}
		}
	}
	return centroids
}

func (nc *NeuralCompressor) Decompress(compressed []byte) []float32 {
	if len(compressed) < 2 {
		return nil
	}
	
	k := int(compressed[0])
	centroids := make([]float32, k)
	
	offset := 1
	for i := 0; i < k; i++ {
		if offset+4 > len(compressed) {
			break
		}
		bits := uint32(compressed[offset])<<24 | uint32(compressed[offset+1])<<16 |
			uint32(compressed[offset+2])<<8 | uint32(compressed[offset+3])
		centroids[i] = math.Float32frombits(bits)
		offset += 4
	}
	
	weights := make([]float32, 0, len(compressed)-offset)
	for i := offset; i < len(compressed); i++ {
		idx := int(compressed[i])
		if idx < len(centroids) {
			weights = append(weights, centroids[idx])
		}
	}
	
	return weights
}

// === EVOLUÇÃO GENÉTICA ===

type EvolutionEngine struct {
	Population   []Genome
	Generation   int
	BestFitness  float32
	MutationRate float32
}

type Genome struct {
	ID         string
	Weights    []float32
	Fitness    float32
	Size       int
	Compressed bool
}

func (ee *EvolutionEngine) Evolve(populationSize, generations int) *Genome {
	ee.Population = make([]Genome, populationSize)
	for i := range ee.Population {
		ee.Population[i] = ee.createRandomGenome(500)
	}
	
	ee.Generation = 0
	ee.MutationRate = 0.1
	
	for g := 0; g < generations; g++ {
		ee.Generation = g
		
		for i := range ee.Population {
			ee.Population[i].Fitness = ee.evaluateGenome(&ee.Population[i])
		}
		
		ee.sortByFitness()
		
		if ee.Generation == 0 || ee.Population[0].Fitness > ee.BestFitness {
			ee.BestFitness = ee.Population[0].Fitness
		}
		
		ee.selection()
		ee.mutation()
		
		fmt.Printf("  Geração %d: fitness=%.4f, tamanho=%d bytes\n", 
			g, ee.BestFitness, ee.Population[0].Size)
	}
	
	return &ee.Population[0]
}

func (ee *EvolutionEngine) createRandomGenome(size int) Genome {
	weights := make([]float32, size)
	for i := range weights {
		weights[i] = float32(rand.Float64()*2 - 1)
	}
	
	nc := NewNeuralCompressor()
	compressed := nc.Compress(weights)
	
	return Genome{
		ID:         fmt.Sprintf("genome-%d", rand.Intn(10000)),
		Weights:    weights,
		Fitness:    0,
		Size:       len(compressed),
		Compressed: true,
	}
}

func (ee *EvolutionEngine) evaluateGenome(g *Genome) float32 {
	var sum float32
	for _, w := range g.Weights {
		sum += abs32(w)
	}
	avgWeight := sum / float32(len(g.Weights))
	fitness := 1.0 - (abs32(avgWeight-0.5) / 0.5)
	fitness = fitness*0.8 + 0.2
	
	sizeFactor := 1.0 - (float32(g.Size) / float32(len(g.Weights)*4))
	fitness = fitness*0.7 + sizeFactor*0.3
	
	return fitness
}

func (ee *EvolutionEngine) sortByFitness() {
	for i := 0; i < len(ee.Population)-1; i++ {
		for j := i + 1; j < len(ee.Population); j++ {
			if ee.Population[j].Fitness > ee.Population[i].Fitness {
				ee.Population[i], ee.Population[j] = ee.Population[j], ee.Population[i]
			}
		}
	}
}

func (ee *EvolutionEngine) selection() {
	keep := len(ee.Population) / 5
	newPop := make([]Genome, len(ee.Population))
	
	for i := 0; i < keep; i++ {
		newPop[i] = ee.Population[i]
	}
	
	for i := keep; i < len(ee.Population); i++ {
		parent1 := ee.Population[rand.Intn(keep)]
		parent2 := ee.Population[rand.Intn(keep)]
		newPop[i] = ee.crossover(parent1, parent2)
	}
	
	ee.Population = newPop
}

func (ee *EvolutionEngine) crossover(p1, p2 Genome) Genome {
	child := Genome{
		ID:         fmt.Sprintf("child-%d", rand.Intn(10000)),
		Weights:    make([]float32, len(p1.Weights)),
		Compressed: true,
	}
	
	crossoverRate := float32(0.7)
	for i := range child.Weights {
		if rand.Float32() < crossoverRate {
			child.Weights[i] = p1.Weights[i]
		} else {
			child.Weights[i] = p2.Weights[i]
		}
	}
	
	nc := NewNeuralCompressor()
	compressed := nc.Compress(child.Weights)
	child.Size = len(compressed)
	
	return child
}

func (ee *EvolutionEngine) mutation() {
	for i := 1; i < len(ee.Population); i++ {
		for j := range ee.Population[i].Weights {
			if rand.Float32() < ee.MutationRate {
				ee.Population[i].Weights[j] += float32(rand.NormFloat64() * 0.1)
			}
		}
	}
}

// === REDE NEURAL LEVE ===

type LightNeuron struct {
	ID        int
	Weights   []float32
	Threshold float32
	Output    float32
}

type LightLayer struct {
	Name   string
	Type   string
	Neurons []LightNeuron
}

type NeuralCore struct {
	Layers  []LightLayer
	Version string
}

func NewLightNeuralNet(input, hidden, output int) *NeuralCore {
	nc := &NeuralCore{
		Version: "1.0-pandora",
		Layers: []LightLayer{
			{Name: "input", Type: "input", Neurons: make([]LightNeuron, input)},
			{Name: "hidden", Type: "hidden", Neurons: make([]LightNeuron, hidden)},
			{Name: "output", Type: "output", Neurons: make([]LightNeuron, output)},
		},
	}
	
	// Inicializar neurônios ocultos e saída com pesos
	rand.Seed(time.Now().UnixNano())
	for i := range nc.Layers[1].Neurons {
		weights := make([]float32, input)
		for j := range weights {
			weights[j] = float32(rand.Float64()*2 - 1)
		}
		nc.Layers[1].Neurons[i] = LightNeuron{ID: i, Weights: weights, Threshold: 0.5}
	}
	
	for i := range nc.Layers[2].Neurons {
		weights := make([]float32, hidden)
		for j := range weights {
			weights[j] = float32(rand.Float64()*2 - 1)
		}
		nc.Layers[2].Neurons[i] = LightNeuron{ID: i, Weights: weights, Threshold: 0.5}
	}
	
	return nc
}

func (nc *NeuralCore) Forward(input []float32) []float32 {
	// Input -> Hidden
	hidden := make([]float32, len(nc.Layers[1].Neurons))
	for i, neuron := range nc.Layers[1].Neurons {
		var sum float32
		for j, w := range neuron.Weights {
			if j < len(input) {
				sum += input[j] * w
			}
		}
		sum += neuron.Threshold
		hidden[i] = sigmoid(sum)
		nc.Layers[1].Neurons[i].Output = hidden[i]
	}
	
	// Hidden -> Output
	output := make([]float32, len(nc.Layers[2].Neurons))
	for i, neuron := range nc.Layers[2].Neurons {
		var sum float32
		for j, w := range neuron.Weights {
			if j < len(hidden) {
				sum += hidden[j] * w
			}
		}
		sum += neuron.Threshold
		output[i] = sigmoid(sum)
		nc.Layers[2].Neurons[i].Output = output[i]
	}
	
	return output
}

func sigmoid(x float32) float32 {
	if x > 10 { return 1 }
	if x < -10 { return 0 }
	return 1.0 / (1.0 + float32(math.Exp(-float64(x))))
}

// === MÁQUINA UNIVERSAL ===

type UniversalMachine struct {
	Opcodes   map[string]Operation
	Registers map[string]float32
	Memory    []float32
	PC        int
}

type Operation struct {
	Name     string
	Execute  func(vm *UniversalMachine, params []string)
}

func NewUniversalMachine() *UniversalMachine {
	vm := &UniversalMachine{
		Opcodes:   make(map[string]Operation),
		Registers: make(map[string]float32),
		Memory:    make([]float32, 256),
		PC:        0,
	}
	
	vm.Opcodes = map[string]Operation{
		"LOAD": {Name: "LOAD", Execute: func(vm *UniversalMachine, p []string) {
			if len(p) >= 2 {
				vm.Registers[p[0]] = vm.Memory[strToAddr(p[1])]
			}
		}},
		"SAVE": {Name: "SAVE", Execute: func(vm *UniversalMachine, p []string) {
			if len(p) >= 2 {
				vm.Memory[strToAddr(p[1])] = vm.Registers[p[0]]
			}
		}},
		"ADD": {Name: "ADD", Execute: func(vm *UniversalMachine, p []string) {
			if len(p) >= 2 {
				vm.Registers[p[0]] += vm.Registers[p[1]]
			}
		}},
		"MUL": {Name: "MUL", Execute: func(vm *UniversalMachine, p []string) {
			if len(p) >= 2 {
				vm.Registers[p[0]] *= vm.Registers[p[1]]
			}
		}},
		"SIGMOID": {Name: "SIGMOID", Execute: func(vm *UniversalMachine, p []string) {
			if len(p) >= 1 {
				vm.Registers[p[0]] = sigmoid(vm.Registers[p[0]])
			}
		}},
		"VADD": {Name: "VADD", Execute: func(vm *UniversalMachine, p []string) {
			// Vector add - somar vetores na memória
			if len(p) >= 3 {
				start1 := strToAddr(p[0])
				start2 := strToAddr(p[1])
				count := strToAddr(p[2])
				for i := 0; i < count; i++ {
					vm.Memory[start1+i] += vm.Memory[start2+i]
				}
			}
		}},
		"EMBD": {Name: "EMBD", Execute: func(vm *UniversalMachine, p []string) {
			// Gerar embedding e salvar
			if len(p) >= 2 {
				text := p[0]
				addr := strToAddr(p[1])
				emb := generateEmbedding(text)
				for i, v := range emb {
					if addr+i < len(vm.Memory) {
						vm.Memory[addr+i] = v
					}
				}
			}
		}},
	}
	
	return vm
}

func strToAddr(s string) int {
	var addr int
	fmt.Sscanf(s, "%d", &addr)
	return addr % 256
}

func (vm *UniversalMachine) Run(program []string) {
	for vm.PC < len(program) {
		parts := strings.Fields(program[vm.PC])
		if len(parts) == 0 {
			vm.PC++
			continue
		}
		
		op := parts[0]
		params := parts[1:]
		
		if opCode, ok := vm.Opcodes[op]; ok {
			opCode.Execute(vm, params)
		}
		
		vm.PC++
	}
}

// === COMUNICAÇÃO UNIVERSAL ===

type UniversalComm struct {
	Protocols map[string]Encoder
}

type Encoder interface {
	Encode(interface{}) []byte
	Decode([]byte) interface{}
}

func NewUniversalComm() *UniversalComm {
	uc := &UniversalComm{Protocols: make(map[string]Encoder)}
	uc.Protocols["pandora"] = &PandoraEncoder{}
	uc.Protocols["json"] = &JSONEncoder{}
	uc.Protocols["binary"] = &BinaryEncoder{}
	return uc
}

type PandoraEncoder struct{}
func (pe *PandoraEncoder) Encode(v interface{}) []byte {
	switch val := v.(type) {
	case string:
		enc := NewPandoraEncoding()
		emb := generateEmbedding(val)
		return []byte(enc.Encode(emb))
	case []float32:
		enc := NewPandoraEncoding()
		return []byte(enc.Encode(val))
	default:
		return []byte(fmt.Sprintf("%v", val))
	}
}
func (pe *PandoraEncoder) Decode(data []byte) interface{} { return string(data) }

type JSONEncoder struct{}
func (je *JSONEncoder) Encode(v interface{}) []byte {
	data, _ := json.Marshal(v)
	return data
}
func (je *JSONEncoder) Decode(data []byte) interface{} {
	var v interface{}
	json.Unmarshal(data, &v)
	return v
}

type BinaryEncoder struct{}
func (be *BinaryEncoder) Encode(v interface{}) []byte {
	h := sha256.New()
	h.Write([]byte(fmt.Sprintf("%v", v)))
	return h.Sum(nil)
}
func (be *BinaryEncoder) Decode(data []byte) interface{} {
	if len(data) >= 8 {
		return binary.BigEndian.Uint64(data[:8])
	}
	return nil
}

// === UTILITÁRIOS ===

func abs32(x float32) float32 {
	if x < 0 { return -x }
	return x
}

func generateEmbedding(text string) []float32 {
	h := sha256.Sum256([]byte(text))
	r := binary.BigEndian.Uint64(h[:8])
	
	values := make([]float32, 64)
	for i := range values {
		values[i] = float32(r%1000) / 1000.0
		r = r*1103515245 + 12345
	}
	
	var sum float32
	for _, v := range values {
		sum += v * v
	}
	sum = float32(math.Sqrt(float64(sum)))
	if sum > 0 {
		for i := range values {
			values[i] /= sum
		}
	}
	
	return values
}

// === MAIN ===

func main() {
	rand.Seed(time.Now().UnixNano())
	
	fmt.Println("╔══════════════════════════════════════════════════════════════╗")
	fmt.Println("║  🧬 PANDORA NEURAL CORE - Evolução & Compressão v1.0        ║")
	fmt.Println("╠══════════════════════════════════════════════════════════════╣")
	
	// 1. Encoding Proprietário
	enc := NewPandoraEncoding()
	testVec := []float32{0.1, -0.5, 0.8, -0.3, 0.5, 0.2, -0.7, 0.4}
	encoded := enc.Encode(testVec)
	decoded := enc.Decode(encoded)
	
	fmt.Printf("║ 1. Encoding Pandora:\n")
	fmt.Printf("║    Original:    %v\n", testVec[:4])
	fmt.Printf("║    Codificado:  %s...\n", encoded[:min(20, len(encoded))])
	fmt.Printf("║    Decodificado: %v\n", decoded[:4])
	
	// 2. Compressão Neural
	nc := NewNeuralCompressor()
	weights := make([]float32, 1000)
	for i := range weights {
		weights[i] = float32(rand.Float64()*2 - 1)
	}
	compressed := nc.Compress(weights)
	decompressed := nc.Decompress(compressed)
	
	fmt.Printf("║ 2. Compressão Neural:\n")
	fmt.Printf("║    Original:    %d bytes\n", len(weights)*4)
	fmt.Printf("║    Comprimido:  %d bytes\n", len(compressed))
	fmt.Printf("║    Ratio:       %.2fx\n", nc.Ratio)
	fmt.Printf("║    Decompress:  %d valores\n", len(decompressed))
	
	// 3. Evolução Genética
	fmt.Printf("║ 3. Evolução Genética (10 gerações):\n")
	ee := &EvolutionEngine{}
	bestGenome := ee.Evolve(15, 10)
	fmt.Printf("║    Melhor genome: fitness=%.4f, tamanho=%d bytes\n", 
		bestGenome.Fitness, bestGenome.Size)
	
	// 4. Rede Neural Leve
	nnet := NewLightNeuralNet(8, 4, 2)
	input := []float32{0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8}
	output := nnet.Forward(input)
	fmt.Printf("║ 4. Rede Neural Leve:\n")
	fmt.Printf("║    Camadas: %d\n", len(nnet.Layers))
	fmt.Printf("║    Input: %v\n", input[:4])
	fmt.Printf("║    Output: %v\n", output)
	
	// 5. Máquina Universal
	vm := NewUniversalMachine()
	vm.Memory[0] = 10.0
	vm.Memory[1] = 5.0
	program := []string{
		"LOAD R0 0",
		"LOAD R1 1",
		"ADD R0 R1",
		"MUL R0 R0",
	}
	vm.Run(program)
	fmt.Printf("║ 5. Máquina Universal:\n")
	fmt.Printf("║    Opcodes: %d\n", len(vm.Opcodes))
	fmt.Printf("║    10 + 5 = %v -> quadrado = %.2f\n", vm.Registers["R0"], vm.Registers["R0"])
	
	// 6. Comunicação Universal
	uc := NewUniversalComm()
	testData := map[string]interface{}{"type": "neural", "action": "compute", "version": 1.0}
	pandoraEncoded := uc.Protocols["pandora"].Encode(testData)
	jsonEncoded := uc.Protocols["json"].Encode(testData)
	binaryEncoded := uc.Protocols["binary"].Encode(testData)
	
	fmt.Printf("║ 6. Comunicação Universal:\n")
	fmt.Printf("║    Pandora:  %d bytes (%s...)\n", len(pandoraEncoded), string(pandoraEncoded[:min(15, len(pandoraEncoded))]))
	fmt.Printf("║    JSON:     %d bytes\n", len(jsonEncoded))
	fmt.Printf("║    Binary:   %d bytes\n", len(binaryEncoded))
	
	fmt.Println("╚══════════════════════════════════════════════════════════════╝")
}

func min(a, b int) int {
	if a < b { return a }
	return b
}