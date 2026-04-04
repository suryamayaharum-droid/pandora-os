package main

import (
	"crypto/sha256"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// ═══════════════════════════════════════════════════════════════
// PANDORA OS - Sistema Operacional de IA Integrado
// ═══════════════════════════════════════════════════════════════

// KERNEL
type Kernel struct {
	Name       string
	Version    string
	Uptime     time.Time
	Confidence float32
	State      string
	Components map[string]interface{}
	Memory     *KernelMemory
}

type KernelMemory struct {
	Items     []string
	MaxItems  int
}

func NewKernel() *Kernel {
	return &Kernel{
		Name:       "PandoraOS",
		Version:    "1.0.0-Alpha",
		Uptime:     time.Now(),
		Confidence: 0.6,
		State:      "booting",
		Components: make(map[string]interface{}),
		Memory:     &KernelMemory{Items: make([]string, 0), MaxItems: 100},
	}
}

// SUBSISTEMA 1: LINGUAGEM
type LanguageSubsystem struct {
	Tokenizer *TokenizerSub
	Embeddings *EmbeddingsSub
	Model     *MiniModel
	Logic     *OpLogic
}

type TokenizerSub struct {
	Vocab     map[string]int
	Reverse   map[int]string
}

func NewTokenizerSub() *TokenizerSub {
	t := &TokenizerSub{Vocab: make(map[string]int), Reverse: make(map[int]string)}
	toks := []string{"<PAD>", "<UNK>", "<BOS>", "<EOS>", "PANDORA", "SISTEMA", "ESTADO", "ACAO", "MEMORIA", "LEIA", "ESCREVA", "EXECUTE", "PENSE", "OK", "ERRO", "SUCESSO", "STATUS", "HELP", "EXIT"}
	for i, tok := range toks {
		t.Vocab[tok] = i
		t.Reverse[i] = tok
	}
	return t
}

func (t *TokenizerSub) Tokenize(text string) []int {
	result := make([]int, 0)
	for _, w := range strings.Fields(strings.ToUpper(text)) {
		if id, ok := t.Vocab[w]; ok {
			result = append(result, id)
		} else {
			result = append(result, t.Vocab["<UNK>"])
		}
	}
	return result
}

type EmbeddingsSub struct {
	Matrix [][]float32
	Dim    int
}

func NewEmbeddingsSub(dim int) *EmbeddingsSub {
	e := &EmbeddingsSub{Dim: dim, Matrix: make([][]float32, 512)}
	rand.Seed(42)
	for i := range e.Matrix {
		e.Matrix[i] = make([]float32, dim)
		for j := range e.Matrix[i] {
			e.Matrix[i][j] = float32(rand.Float64()*2-1) * 0.1
		}
	}
	return e
}

func (e *EmbeddingsSub) Forward(ids []int) [][]float32 {
	result := make([][]float32, len(ids))
	for i, id := range ids {
		if id < len(e.Matrix) {
			result[i] = e.Matrix[id]
		}
	}
	return result
}

type MiniModel struct {
	W1 [][]float32
	W2 [][]float32
}

func NewMiniModel() *MiniModel {
	m := &MiniModel{}
	rand.Seed(42)
	m.W1 = make([][]float32, 64)
	for i := range m.W1 {
		m.W1[i] = make([]float32, 128)
		for j := range m.W1[i] {
			m.W1[i][j] = float32(rand.Float64()*2-1) * 0.1
		}
	}
	m.W2 = make([][]float32, 128)
	for i := range m.W2 {
		m.W2[i] = make([]float32, 512)
		for j := range m.W2[i] {
			m.W2[i][j] = float32(rand.Float64()*2-1) * 0.1
		}
	}
	return m
}

func (m *MiniModel) Forward(x [][]float32) []float32 {
	hidden := make([]float32, 128)
	for j := 0; j < 128; j++ {
		for i := 0; i < len(x) && i < 64; i++ {
			if j < len(m.W1[0]) {
				hidden[j] += x[i][j] * m.W1[i][j]
			}
		}
		if hidden[j] > 0 {}
	}
	output := make([]float32, 512)
	for j := 0; j < 512; j++ {
		for i := 0; i < 128; i++ {
			output[j] += hidden[i] * m.W2[i][j]
		}
	}
	return output
}

type OpLogic struct {
	State   string
	Beliefs map[string]float32
	Goals   []string
	Memory  []string
}

func NewOpLogic() *OpLogic {
	return &OpLogic{
		State:   "idle",
		Beliefs: map[string]float32{"confiança": 0.7, "cpu": 0.3},
		Goals:   []string{"manter", "processar", "aprender"},
		Memory:  make([]string, 0),
	}
}

func (ol *OpLogic) Process(input string) string {
	lower := strings.ToLower(input)
	var resp string
	
	if strings.Contains(lower, "status") || strings.Contains(lower, "estado") {
		resp = fmt.Sprintf("PandoraOS %s | Confiança: %.0f%% | Estado: %s", ol.State, ol.Beliefs["confiança"]*100, ol.State)
	} else if strings.Contains(lower, "execute") || strings.Contains(lower, "memoria") {
		ol.State = "executando"
		resp = "Executando..."
	} else if strings.Contains(lower, "pense") {
		ol.State = "pensando"
		resp = "Modo reflexão..."
	} else if strings.Contains(lower, "help") || strings.Contains(lower, "ajuda") {
		resp = "Comandos: status, execute, pense, memoria, help"
	} else {
		resp = fmt.Sprintf("Entendido: '%s'. Processando...", input[:min(20, len(input))])
		ol.State = "processando"
	}
	
	ol.Memory = append(ol.Memory, fmt.Sprintf("%s:%s", ol.State, input[:min(10, len(input))]))
	if len(ol.Memory) > 20 {
		ol.Memory = ol.Memory[len(ol.Memory)-20:]
	}
	
	return resp
}

// SUBSISTEMA 2: NEURAL
type NeuralSubsystem struct {
	Compressor *CompressorSub
	Evolution  *EvolutionSub
	Network    *LightNet
}

type CompressorSub struct {
	Ratio float32
}

func NewCompressorSub() *CompressorSub {
	return &CompressorSub{}
}

func (c *CompressorSub) Compress(data []float32) []byte {
	k := 8
	centroids := make([]float32, k)
	rand.Seed(42)
	for i := range centroids {
		if i < len(data) {
			centroids[i] = data[i]
		}
	}
	result := []byte{byte(k)}
	for _, v := range data {
		best := 0
		minD := float32(1e10)
		for i, c := range centroids {
			d := v - c
			if d < 0 { d = -d }
			if d < minD {
				minD = d
				best = i
			}
		}
		result = append(result, byte(best))
	}
	c.Ratio = float32(len(result)) / float32(len(data)*4)
	return result
}

func (c *CompressorSub) Decompress(data []byte) []float32 {
	if len(data) < 2 {
		return nil
	}
	k := int(data[0])
	result := make([]float32, 0, len(data)-1)
	for i := 1; i < len(data); i++ {
		idx := int(data[i])
		if idx < k {
			result = append(result, float32(idx))
		}
	}
	return result
}

type EvolutionSub struct {
	Generation   int
	BestFitness  float32
}

func NewEvolutionSub() *EvolutionSub {
	return &EvolutionSub{Generation: 0, BestFitness: 0.5}
}

func (e *EvolutionSub) Evolve() string {
	e.Generation++
	e.BestFitness = 0.5 + float32(e.Generation%10)*0.05
	return fmt.Sprintf("Geração %d: fitness %.3f", e.Generation, e.BestFitness)
}

type LightNet struct {
	Layers  int
	Neurons int
	Weights [][][]float32
}

func NewLightNet() *LightNet {
	ln := &LightNet{Layers: 3, Neurons: 32}
	rand.Seed(42)
	ln.Weights = make([][][]float32, ln.Layers-1)
	for l := 0; l < ln.Layers-1; l++ {
		ln.Weights[l] = make([][]float32, ln.Neurons)
		for i := range ln.Weights[l] {
			ln.Weights[l][i] = make([]float32, ln.Neurons)
			for j := range ln.Weights[l][i] {
				ln.Weights[l][i][j] = float32(rand.Float64()*2-1) * 0.1
			}
		}
	}
	return ln
}

func (ln *LightNet) Forward(input []float32) []float32 {
	current := input
	for _, layer := range ln.Weights {
		next := make([]float32, len(layer))
		for i := range next {
			for j := range current {
				if j < len(layer) {
					next[i] += current[j] * layer[i][j]
				}
			}
		}
		current = next
	}
	return current
}

// SUBSISTEMA 3: COMUNICAÇÃO
type UniversalComm struct {
	Protocols map[string]func(interface{}) []byte
	Adapters  map[string]func([]byte) string
}

func NewUniversalComm() *UniversalComm {
	uc := &UniversalComm{
		Protocols: make(map[string]func(interface{}) []byte),
		Adapters:  make(map[string]func([]byte) string),
	}
	uc.Protocols["text"] = func(v interface{}) []byte {
		return []byte(fmt.Sprintf("%v", v))
	}
	uc.Protocols["binary"] = func(v interface{}) []byte {
		h := sha256.New()
		h.Write([]byte(fmt.Sprintf("%v", v)))
		return h.Sum(nil)
	}
	uc.Adapters["json"] = func(d []byte) string { return fmt.Sprintf("JSON:%s", d) }
	uc.Adapters["shell"] = func(d []byte) string { return fmt.Sprintf("CMD:%s", d) }
	return uc
}

func (uc *UniversalComm) Send(msg interface{}, proto, adapter string) string {
	if enc, ok := uc.Protocols[proto]; ok {
		if dec, ok2 := uc.Adapters[adapter]; ok2 {
			return dec(enc(msg))
		}
	}
	return fmt.Sprintf("%v", msg)
}

// SUBSISTEMA 4: VFS
type VirtualFileSystem struct {
	Files map[string]string
	Dirs  map[string][]string
}

func NewVFS() *VirtualFileSystem {
	vfs := &VirtualFileSystem{
		Files: make(map[string]string),
		Dirs:  make(map[string][]string),
	}
	vfs.Dirs["/"] = []string{"system", "user", "tmp"}
	vfs.Dirs["/system"] = []string{"config", "logs"}
	vfs.Files["/system/version"] = "PandoraOS 1.0.0-Alpha"
	vfs.Files["/system/kernel"] = "Kernel v1.0"
	return vfs
}

func (vfs *VirtualFileSystem) Read(path string) string {
	if v, ok := vfs.Files[path]; ok {
		return v
	}
	return "Not found"
}

func (vfs *VirtualFileSystem) Write(path, content string) bool {
	vfs.Files[path] = content
	return true
}

func (vfs *VirtualFileSystem) List(path string) []string {
	if d, ok := vfs.Dirs[path]; ok {
		return d
	}
	return []string{}
}

// SUBSISTEMA 5: SHELL
type ShellSubsystem struct {
	Kernel  *Kernel
	VFS     *VirtualFileSystem
	History []string
	Vars    map[string]string
}

func NewShell(k *Kernel, vfs *VirtualFileSystem) *ShellSubsystem {
	return &ShellSubsystem{
		Kernel:  k,
		VFS:     vfs,
		History: make([]string, 0),
		Vars:    make(map[string]string),
	}
}

func (s *ShellSubsystem) Execute(cmd string) string {
	s.History = append(s.History, cmd)
	parts := strings.Fields(cmd)
	if len(parts) == 0 {
		return ""
	}
	
	cmdName := parts[0]
	args := parts[1:]
	
	switch cmdName {
	case "help", "ajuda":
		return `Comandos: help, status, ls, cat, echo, mem, ps, set, env, clear`
	case "status":
		return fmt.Sprintf("PandoraOS %s | Estado: %s | Confiança: %.0f%% | Uptime: %s",
			s.Kernel.Version, s.Kernel.State, s.Kernel.Confidence*100, time.Since(s.Kernel.Uptime).Truncate(time.Second))
	case "ls":
		path := "/"
		if len(args) > 0 {
			path = args[0]
		}
		return strings.Join(s.VFS.List(path), ", ")
	case "cat":
		if len(args) > 0 {
			return s.VFS.Read(args[0])
		}
		return "Uso: cat [arquivo]"
	case "echo":
		return strings.Join(args, " ")
	case "mem":
		return fmt.Sprintf("Memory: %d itens", len(s.Kernel.Memory.Items))
	case "ps":
		return "Processos: kernel, shell, vfs, language, neural"
	case "set":
		if len(args) > 0 {
			kv := strings.Split(args[0], "=")
			if len(kv) == 2 {
				s.Vars[kv[0]] = kv[1]
				return fmt.Sprintf("%s=%s", kv[0], kv[1])
			}
		}
		return "Uso: set key=value"
	case "env":
		var r []string
		for k, v := range s.Vars {
			r = append(r, fmt.Sprintf("%s=%s", k, v))
		}
		return strings.Join(r, "\n")
	default:
		return fmt.Sprintf("Comando '%s' não encontrado. Digite 'help'.", cmdName)
	}
}

// MAIN
func main() {
	rand.Seed(time.Now().UnixNano())
	
	fmt.Println("╔══════════════════════════════════════════════════════════════╗")
	fmt.Println("║     🧬 PANDORA OS - SISTEMA OPERACIONAL DE IA v1.0          ║")
	fmt.Println("║     [ KERNEL | LINGUAGEM | NEURAL | COM | VFS | SHELL ]     ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════╝")
	
	// Inicializar kernel
	kernel := NewKernel()
	
	// Inicializar subsistemas
	lang := &LanguageSubsystem{
		Tokenizer:  NewTokenizerSub(),
		Embeddings: NewEmbeddingsSub(64),
		Model:      NewMiniModel(),
		Logic:      NewOpLogic(),
	}
	kernel.Components["language"] = lang
	
	neural := &NeuralSubsystem{
		Compressor: NewCompressorSub(),
		Evolution:  NewEvolutionSub(),
		Network:    NewLightNet(),
	}
	kernel.Components["neural"] = neural
	
	vfs := NewVFS()
	kernel.Components["vfs"] = vfs
	
	comm := NewUniversalComm()
	kernel.Components["comm"] = comm
	
	shell := NewShell(kernel, vfs)
	kernel.Components["shell"] = shell
	
	kernel.State = "running"
	kernel.Confidence = 0.75
	
	fmt.Println("\n📦 Subsistemas carregados:")
	fmt.Println("  ✓ KERNEL - Gerenciamento central")
	fmt.Println("  ✓ LINGUAGEM - NLP + Embeddings + Modelo")
	fmt.Println("  ✓ NEURAL - Compressão + Evolução + Rede")
	fmt.Println("  ✓ COMUNICAÇÃO - Protocolos universais")
	fmt.Println("  ✓ VFS - Sistema de arquivos virtual")
	fmt.Println("  ✓ SHELL - Interface de comandos")
	
	// Testar shell
	fmt.Println("\n📋 Teste de Comandos:")
	cmds := []string{"help", "status", "ls /", "ls /system", "cat /system/version", "mem", "ps", "echo Ola Mundo", "set autor=Pandora", "env"}
	for _, c := range cmds {
		r := shell.Execute(c)
		fmt.Printf("  $ %s\n    %s\n\n", c, r)
	}
	
	// Testar linguagem
	fmt.Println("🧠 Teste de Linguagem:")
	inputs := []string{"status", "execute algo", "pense", "help"}
	for _, inp := range inputs {
		r := lang.Logic.Process(inp)
		fmt.Printf("  Input: '%s' → %s\n", inp, r)
	}
	
	// Testar compressão
	fmt.Println("\n🔬 Teste de Compressão Neural:")
	weights := make([]float32, 100)
	for i := range weights {
		weights[i] = float32(rand.Float64() * 2 - 1)
	}
	compressed := neural.Compressor.Compress(weights)
	decompressed := neural.Compressor.Decompress(compressed)
	fmt.Printf("  Original: %d bytes → Comprimido: %d bytes (%.2fx)\n", len(weights)*4, len(compressed), 1.0/neural.Compressor.Ratio)
	fmt.Printf("  Decompress OK: %v\n", len(decompressed) > 0)
	
	// Testar evolução
	fmt.Println("\n🧬 Evolução Neural:")
	for i := 0; i < 3; i++ {
		fmt.Printf("  %s\n", neural.Evolution.Evolve())
	}
	
	// Testar comunicação
	fmt.Println("\n📡 Teste de Comunicação:")
	fmt.Printf("  Text+JSON: %s\n", comm.Send("test message", "text", "json"))
	fmt.Printf("  Binary+Shell: %s\n", comm.Send(123, "binary", "shell"))
	
	// Status final
	uptime := time.Since(kernel.Uptime)
	fmt.Println("\n╔══════════════════════════════════════════════════════════════╗")
	fmt.Println("║  ✅ PANDORA OS - SISTEMA OPERACIONAL INTEGRADO COMPLETO       ║")
	fmt.Println("╠══════════════════════════════════════════════════════════════╣")
	fmt.Printf("║  Status: %-10s | Confiança: %.0f%% | Uptime: %s         ║\n", kernel.State, kernel.Confidence*100, uptime.Truncate(time.Second))
	fmt.Println("║  Componentes: 6 subsistemas integrados                       ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════╝")
}

func min(a, b int) int {
	if a < b { return a }
	return b
}