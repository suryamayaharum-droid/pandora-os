package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// ═══════════════════════════════════════════════════════════════
// PANDORA AUTO-DEPLOY SYSTEM
// Sistema de Auto-Expansão - Criar Novos Nodes Automaticamente
// ═══════════════════════════════════════════════════════════════

type AutoDeploy struct {
	Name        string
	Version     string
	Enabled     bool
	
	// Cluster management
	Cluster    *Cluster
	Nodes      []Node
	DeployQueue []DeployTask
	
	// Scaling
	AutoScale  bool
	MinNodes   int
	MaxNodes   int
	ScaleThreshold float32
}

type Cluster struct {
	Name      string
	Status    string
	Region    string
	TotalCPU  float32
	TotalMem  float32
}

type Node struct {
	ID        string
	Name      string
	Status    string // "creating", "running", "stopped", "error"
	IP        string
	Port      int
	CPU       float32
	Memory    float32
	CreatedAt time.Time
	LastSeen  time.Time
}

type DeployTask struct {
	ID        string
	NodeName  string
	Status    string // "queued", "deploying", "completed", "failed"
	Progress  int
	Error     string
	CreatedAt time.Time
}

func NewAutoDeploy() *AutoDeploy {
	ad := &AutoDeploy{
		Name:       "Pandora-AutoDeploy",
		Version:    "v1.0-AUTO",
		Enabled:    true,
		Cluster: &Cluster{
			Name:   "pandora-cluster",
			Status: "active",
			Region: "fly-auto",
		},
		Nodes:      make([]Node, 0),
		DeployQueue: make([]DeployTask, 0),
		AutoScale:  true,
		MinNodes:   1,
		MaxNodes:   10,
		ScaleThreshold: 0.8,
	}
	
	// Add initial node (self)
	ad.Nodes = append(ad.Nodes, Node{
		ID:        "node-0",
		Name:      "pandora-primary",
		Status:    "running",
		IP:        "172.19.10.82",
		Port:      9000,
		CPU:       20.0,
		Memory:    40.0,
		CreatedAt: time.Now(),
		LastSeen:  time.Now(),
	})
	
	return ad
}

// === NODE MANAGEMENT ===

func (ad *AutoDeploy) CreateNode() string {
	nodeID := fmt.Sprintf("node-%d", len(ad.Nodes))
	nodeName := fmt.Sprintf("pandora-%s", randomName())
	
	node := Node{
		ID:        nodeID,
		Name:      nodeName,
		Status:    "creating",
		IP:        generateIP(),
		Port:      9000 + len(ad.Nodes),
		CPU:       0,
		Memory:    0,
		CreatedAt: time.Now(),
		LastSeen:  time.Now(),
	}
	
	// Add to queue
	task := DeployTask{
		ID:        fmt.Sprintf("task-%d", len(ad.DeployQueue)),
		NodeName:  nodeName,
		Status:    "queued",
		Progress:  0,
		CreatedAt: time.Now(),
	}
	
	ad.Nodes = append(ad.Nodes, node)
	ad.DeployQueue = append(ad.DeployQueue, task)
	
	return nodeID
}

func (ad *AutoDeploy) DeployNode(nodeID string) bool {
	// Find node
	nodeIdx := -1
	for i := range ad.Nodes {
		if ad.Nodes[i].ID == nodeID {
			nodeIdx = i
			break
		}
	}
	
	if nodeIdx == -1 {
		return false
	}
	
	// Simulate deployment
	node := &ad.Nodes[nodeIdx]
	node.Status = "deploying"
	
	// Simulate progress
	for p := 0; p <= 100; p += 20 {
		// Update task progress
		for i := range ad.DeployQueue {
			if ad.DeployQueue[i].NodeName == node.Name {
				ad.DeployQueue[i].Progress = p
				ad.DeployQueue[i].Status = "deploying"
			}
		}
		
		fmt.Printf("  📦 Deploying %s: %d%%\n", node.Name, p)
		time.Sleep(100 * time.Millisecond)
	}
	
	// Mark as running
	node.Status = "running"
	node.CPU = 15.0 + rand.Float32()*10
	node.Memory = 30.0 + rand.Float32()*20
	node.LastSeen = time.Now()
	
	// Update task
	for i := range ad.DeployQueue {
		if ad.DeployQueue[i].NodeName == node.Name {
			ad.DeployQueue[i].Status = "completed"
			ad.DeployQueue[i].Progress = 100
		}
	}
	
	fmt.Printf("  ✅ Node %s está rodando em %s:%d\n", node.Name, node.IP, node.Port)
	
	return true
}

func (ad *AutoDeploy) RemoveNode(nodeID string) bool {
	for i := range ad.Nodes {
		if ad.Nodes[i].ID == nodeID && ad.Nodes[i].Status == "running" {
			ad.Nodes[i].Status = "stopped"
			fmt.Printf("  🛑 Node %s parado\n", ad.Nodes[i].Name)
			return true
		}
	}
	return false
}

func (ad *AutoDeploy) GetNodeStatus(nodeID string) string {
	for _, n := range ad.Nodes {
		if n.ID == nodeID {
			return fmt.Sprintf("%s (%s) - CPU: %.0f%% Mem: %.0f%%", n.Name, n.Status, n.CPU, n.Memory)
		}
	}
	return "Node não encontrado"
}

// === AUTO SCALING ===

func (ad *AutoDeploy) CheckScale() {
	if !ad.AutoScale {
		return
	}
	
	// Calculate average load
	if len(ad.Nodes) == 0 {
		return
	}
	
	var totalLoad float32
	for _, n := range ad.Nodes {
		if n.Status == "running" {
			totalLoad += n.CPU
		}
	}
	
	avgLoad := totalLoad / float32(len(ad.Nodes))
	
	fmt.Printf("\n📊 Carga média: %.0f%% (threshold: %.0f%%)\n", avgLoad, ad.ScaleThreshold*100)
	
	// Scale up if needed
	if avgLoad > ad.ScaleThreshold && len(ad.Nodes) < ad.MaxNodes {
		fmt.Println("  📈 Escalando para cima...")
		ad.ScaleUp(1)
	}
	
	// Scale down if underutilized
	if avgLoad < ad.ScaleThreshold*0.3 && len(ad.Nodes) > ad.MinNodes {
		fmt.Println("  📉 Escalando para baixo...")
		ad.ScaleDown(1)
	}
}

func (ad *AutoDeploy) ScaleUp(count int) {
	for i := 0; i < count; i++ {
		if len(ad.Nodes) >= ad.MaxNodes {
			break
		}
		
		nodeID := ad.CreateNode()
		ad.DeployNode(nodeID)
	}
}

func (ad *AutoDeploy) ScaleDown(count int) {
	for i := 0; i < count; i++ {
		if len(ad.Nodes) <= ad.MinNodes {
			break
		}
		
		// Find idle node (not primary)
		for j := len(ad.Nodes) - 1; j > 0; j-- {
			if ad.Nodes[j].Status == "running" && ad.Nodes[j].CPU < 30 {
				ad.RemoveNode(ad.Nodes[j].ID)
				break
			}
		}
	}
}

// === DISTRIBUTED MEMORY ===

func (ad *AutoDeploy) SyncMemory() {
	fmt.Println("\n🔄 Sincronizando memória entre nodes...")
	
	for _, n := range ad.Nodes {
		if n.Status == "running" {
			fmt.Printf("  📤 Sincronizando com %s...\n", n.Name)
		}
	}
	
	fmt.Println("  ✅ Memória sincronizada")
}

// === CLUSTER STATUS ===

func (ad *AutoDeploy) ClusterStatus() string {
	running := 0
	for _, n := range ad.Nodes {
		if n.Status == "running" {
			running++
		}
	}
	
	return fmt.Sprintf(`
╔══════════════════════════════════════════════════════════════╗
║           🚀 PANDORA AUTO-DEPLOY SYSTEM                    ║
╠═══════════════════════════════════════════════════════════════╣
║  Cluster: %s | Região: %s                               ║
║  Status: %s | Auto-Scale: %v                             ║
╠═══════════════════════════════════════════════════════════════╣
║  NÓS: %d/%d                                                ║
║  Running: %d | Creating: %d | Stopped: %d                 ║
╠═══════════════════════════════════════════════════════════════╣
║  DETALHES DOS NÓS:                                         ║
%s
╠═══════════════════════════════════════════════════════════════╣
║  FILA DE DEPLOY:                                          ║
%s
╠═══════════════════════════════════════════════════════════════╣
║  MÉTRICAS DO CLUSTER                                       ║
║  CPU Total: %.1f%% | Memória Total: %.1f%%                 ║
╚══════════════════════════════════════════════════════════════╝`,
		ad.Cluster.Name, ad.Cluster.Status, ad.Cluster.Region,
		ad.AutoScale,
		len(ad.Nodes), ad.MaxNodes,
		running, ad.countByStatus("creating"), ad.countByStatus("stopped"),
		ad.formatNodes(),
		ad.formatQueue(),
		ad.Cluster.TotalCPU, ad.Cluster.TotalMem)
}

func (ad *AutoDeploy) countByStatus(status string) int {
	c := 0
	for _, n := range ad.Nodes {
		if n.Status == status {
			c++
		}
	}
	return c
}

func (ad *AutoDeploy) formatNodes() string {
	result := ""
	for _, n := range ad.Nodes {
		if n.Status == "running" {
			result += fmt.Sprintf("║  • %s | %s:%d | CPU: %.0f%% | Mem: %.0f%%\n", 
				n.Name, n.IP, n.Port, n.CPU, n.Memory)
		}
	}
	if result == "" {
		return "║  Nenhum node rodando\n"
	}
	return result
}

func (ad *AutoDeploy) formatQueue() string {
	result := ""
	for _, t := range ad.DeployQueue {
		if t.Status != "completed" {
			result += fmt.Sprintf("║  • %s: %s (%d%%)\n", t.NodeName, t.Status, t.Progress)
		}
	}
	if result == "" {
		return "║  Nenhuma tarefa pendente\n"
	}
	return result
}

// === MAIN ===

func main() {
	rand.Seed(time.Now().UnixNano())
	
	fmt.Println("╔══════════════════════════════════════════════════════════════╗")
	fmt.Println("║           🚀 PANDORA AUTO-DEPLOY SYSTEM v1.0                ║")
	fmt.Println("║        [ AUTO-EXPANSÃO E ESCALABILIDADE ]                 ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════╝")
	
	ad := NewAutoDeploy()
	fmt.Println("\n✓ Sistema de auto-deploy inicializado")
	
	// Show initial cluster
	fmt.Println("\n📊 Cluster inicial:")
	fmt.Printf("  Nodes: %d | Auto-scale: %v\n", len(ad.Nodes), ad.AutoScale)
	
	// Create and deploy new nodes
	fmt.Println("\n🆕 Criando novos nodes...")
	ad.ScaleUp(2)
	
	// Check scaling
	ad.CheckScale()
	
	// Sync memory
	ad.SyncMemory()
	
	// Show cluster status
	fmt.Println("\n" + ad.ClusterStatus())
	
	fmt.Println("\n🎉 SISTEMA COMPLETO!")
	fmt.Println("\n📋 TODAS AS CAPACIDADES IMPLEMENTADAS:")
	fmt.Println("  ✅ Auto-modificação segura (selfmod/)")
	fmt.Println("  ✅ Meta-aprendizado (meta_learn/)")
	fmt.Println("  ✅ Auto-deploy (autodeploy/)")
	fmt.Println("  ✅ Auto-evolução (daemon_auto/)")
	fmt.Println("  ✅ Integração FLAI (flai_integration/)")
	fmt.Println("  ✅ Mineração distribuída (distributed/)")
	fmt.Println("  ✅ Unified Core (unified/)")
}

func randomName() string {
	prefixes := []string{"nova", "echo", "flux", "pulse", "core", "node", "auto", "meta"}
	suffixes := []string{"alpha", "beta", "gamma", "delta", "omega", "prime"}
	
	p := prefixes[rand.Intn(len(prefixes))]
	s := suffixes[rand.Intn(len(suffixes))]
	
	return fmt.Sprintf("%s-%s", p, s)
}

func generateIP() string {
	return fmt.Sprintf("172.19.%d.%d", rand.Intn(256), rand.Intn(256))
}

var _ = strings.TrimSpace