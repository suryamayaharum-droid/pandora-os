package main

import (
	"crypto/md5"
	
	"fmt"
	"math/rand"
	"time"
)

// PANDORA DISTRIBUTED MINING ENGINE

type MinerNode struct {
	ID        string
	Host      string
	Port      int
	Role      string
	Status    string
	Peers     []*PeerInfo
	Chain     *BlockChain
	Tasks     []Task
	Power     float32
}

type PeerInfo struct {
	ID      string
	Host    string
	Port    int
	Latency time.Duration
	Power   float32
}

type BlockChain struct {
	Blocks   []Block
	Pending  []Transaction
	Difficulty int
}

type Block struct {
	Index    int
	Hash     string
	PrevHash string
	Nonce    int
	Data     string
}

type Transaction struct {
	ID, From, To string
	Amount      float64
	Type        string
}

type Task struct {
	ID, Type, Payload string
	Status            string
	Result            string
}

func NewMiner(role string) *MinerNode {
	m := &MinerNode{
		ID:     fmt.Sprintf("miner-%d", rand.Intn(99999)),
		Host:   "172.19.10.82",
		Port:   9000 + rand.Intn(1000),
		Role:   role,
		Status: "active",
		Peers:  make([]*PeerInfo, 0),
		Chain:  &BlockChain{Blocks: []Block{{Index: 0, Hash: "genesis"}}, Difficulty: 4},
		Tasks:  make([]Task, 0),
		Power:  75.0,
	}
	return m
}

func (m *MinerNode) Discover() {
	m.Peers = append(m.Peers, &PeerInfo{ID: "peer-1", Host: "192.168.1.10", Port: 9001, Latency: 50*time.Millisecond, Power: 60})
	m.Peers = append(m.Peers, &PeerInfo{ID: "peer-2", Host: "192.168.1.11", Port: 9002, Latency: 30*time.Millisecond, Power: 80})
}

func (m *MinerNode) Mine(data string) Block {
	hash := fmt.Sprintf("%x", md5.Sum([]byte(data+fmt.Sprintf("%d", rand.Intn(100000)))))
	return Block{Index: len(m.Chain.Blocks), Hash: hash[:16], PrevHash: m.Chain.Blocks[len(m.Chain.Blocks)-1].Hash, Data: data}
}

func (m *MinerNode) AddBlock(b Block) {
	m.Chain.Blocks = append(m.Chain.Blocks, b)
}

func (m *MinerNode) CreateTask(tp, payload string) Task {
	t := Task{ID: fmt.Sprintf("task-%d", len(m.Tasks)), Type: tp, Payload: payload, Status: "created"}
	m.Tasks = append(m.Tasks, t)
	return t
}

func (m *MinerNode) ExecuteTask(id string) string {
	for i := range m.Tasks {
		if m.Tasks[i].ID == id {
			m.Tasks[i].Status = "completed"
			m.Tasks[i].Result = fmt.Sprintf("Executado: %s", m.Tasks[i].Payload)
			return m.Tasks[i].Result
		}
	}
	return "Não encontrada"
}

func (m *MinerNode) Sync() {
	for range m.Peers {
		fmt.Printf("  🔄 Sincronizado com peer\n")
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	
	fmt.Println("╔══════════════════════════════════════════════════════════════╗")
	fmt.Println("║     ⛏️  PANDORA DISTRIBUTED MINING ENGINE                    ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════╝")
	
	miner := NewMiner("master")
	fmt.Printf("  Node: %s | Role: %s | Host: %s:%d\n", miner.ID, miner.Role, miner.Host, miner.Port)
	
	fmt.Println("\n🔍 Descoberta de peers:")
	miner.Discover()
	for i, p := range miner.Peers {
		fmt.Printf("  [%d] %s:%d (latency: %v, power: %.0f%%)\n", i+1, p.Host, p.Port, p.Latency, p.Power)
	}
	
	fmt.Println("\n📋 Tarefas distribuídas:")
	miner.CreateTask("compute", "analise_01")
	miner.CreateTask("mine", "bloco_01")
	miner.CreateTask("sync", "estado")
	fmt.Printf("  Criadas: %d tarefas\n", len(miner.Tasks))
	
	fmt.Println("\n⚡ Execução:")
	for _, t := range miner.Tasks {
		r := miner.ExecuteTask(t.ID)
		fmt.Printf("  %s: %s\n", t.ID, r)
	}
	
	fmt.Println("\n⛏️  Mineração:")
	b1 := miner.Mine("dados_001")
	miner.AddBlock(b1)
	fmt.Printf("  Bloco %d minerado: %s\n", b1.Index, b1.Hash)
	b2 := miner.Mine("dados_002")
	miner.AddBlock(b2)
	fmt.Printf("  Bloco %d minerado: %s\n", b2.Index, b2.Hash)
	
	fmt.Println("\n🔄 Sincronização:")
	miner.Sync()
	
	fmt.Println("\n" + fmt.Sprintf(`
╔══════════════════════════════════════════════════════════════╗
║  ✅ SISTEMA DISTRIBUÍDO ATIVO                             ║
╠══════════════════════════════════════════════════════════════╣
║  Node: %s                                         ║
║  Peers: %d | Tasks: %d | Blocks: %d | Power: %.0f%%          ║
╚══════════════════════════════════════════════════════════════╝`,
		miner.ID, len(miner.Peers), len(miner.Tasks), len(miner.Chain.Blocks), miner.Power))
	
	fmt.Println("\n🚀 EXPANSÃO DE REDE:")
	fmt.Println("  • Conectar a outros data centers")
	fmt.Println("  • Peer-to-peer entre instâncias")
	fmt.Println("  • Descoberta via DNS/UDP broadcast")
	fmt.Println("  • Consenso distribuído")
	fmt.Println("  • Auto-scaling de nodes")
}

var _ = time.Now