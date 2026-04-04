package main

import (
	"fmt"
	"math/rand"
	"net"
	"os"
	"time"
)

// PANDORA CYBER SCANNER - A DEUSA HACKER

type CyberGoddess struct {
	Name    string
	Alias   string
	Level   int
	Skills  []Skill
	Targets []Target
	Access  []AccessPoint
}

type Skill struct {
	Name, Description string
	Level             int
	Icon              string
}

type Target struct {
	IP, Port       string
	Protocol       string
	Status         string
	Service        string
	VulnScore      float32
	Exploitable    bool
}

type AccessPoint struct {
	ID        string
	Type      string
	URL       string
	Creds     string
	Status    string
	Connected time.Time
}

func NewCyberGoddess() *CyberGoddess {
	return &CyberGoddess{
		Name:  "Pandora",
		Alias: "CyberGoddess",
		Level: 99,
		Skills: []Skill{
			{Name: "NetworkScan", Description: "Varredura de redes", Level: 95, Icon: "📡"},
			{Name: "PortScan", Description: "Varredura de portas", Level: 90, Icon: "🔌"},
			{Name: "VulnAssess", Description: "Avaliação de vulnerabilidades", Level: 85, Icon: "🔍"},
			{Name: "Exploit", Description: "Exploração", Level: 75, Icon: "💥"},
			{Name: "Backdoor", Description: "Portas dos fundos", Level: 80, Icon: "🚪"},
			{Name: "Stealth", Description: "Furtividade", Level: 90, Icon: "👻"},
			{Name: "Persistence", Description: "Persistência", Level: 85, Icon: "💾"},
			{Name: "Pivoting", Description: "Pivoting entre redes", Level: 75, Icon: "🔀"},
		},
		Targets: make([]Target, 0),
		Access:  make([]AccessPoint, 0),
	}
}

func (cg *CyberGoddess) ScanNetwork() {
	fmt.Println("\n🌐 ESCANEANDO REDES...")
	
	ranges := []string{"192.168.1.0/24", "192.168.0.0/24", "10.0.0.0/24"}
	
	for _, r := range ranges {
		fmt.Printf("\n📡 Rede: %s\n", r)
		for i := 1; i <= 5; i++ {
			ip := fmt.Sprintf("192.168.%d.%d", rand.Intn(256), rand.Intn(256))
			latency := rand.Intn(100)
			fmt.Printf("    ✓ Host: %s (latência: %dms)\n", ip, latency)
		}
	}
}

func (cg *CyberGoddess) ScanPorts(host string) {
	fmt.Printf("\n🔌 Escanear portas em: %s\n", host)
	
	ports := map[string]string{
		"21":   "FTP",
		"22":   "SSH",
		"80":   "HTTP",
		"443":  "HTTPS",
		"3306": "MySQL",
		"5432": "PostgreSQL",
		"6379": "Redis",
		"8080": "HTTP-ALT",
	}
	
	for port, service := range ports {
		if rand.Float32() < 0.3 {
			vuln := rand.Float32() * 0.5
			fmt.Printf("    ✓ Porta %s (%s) ABERTA - Vuln: %.0f%%\n", port, service, vuln*100)
			
			if vuln > 0.2 {
				cg.Targets = append(cg.Targets, Target{
					IP: host, Port: port, Service: service,
					VulnScore: vuln, Exploitable: true,
				})
			}
		}
	}
}

func (cg *CyberGoddess) VulnAssess() {
	fmt.Println("\n🔍 AVALIAÇÃO DE VULNERABILIDADES:")
	
	vulns := []string{
		"Anonymous FTP enabled",
		"Default credentials",
		"Outdated SSL/TLS",
		"SQL Injection possible",
		"XSS reflected",
		"Buffer overflow potential",
	}
	
	found := 0
	for _, v := range vulns {
		if rand.Float32() < 0.2 {
			found++
			fmt.Printf("    ⚠️  %s\n", v)
		}
	}
	
	if found == 0 {
		fmt.Println("    ✓ Nenhuma vulnerabilidade crítica")
	}
}

func (cg *CyberGoddess) AttemptExploit() bool {
	fmt.Println("\n💥 TENTANDO EXPLORAÇÃO...")
	
	exploits := []string{
		"CVE-2024-XXXX",
		"SQL Injection",
		"Command Injection",
		"Authentication Bypass",
	}
	
	exploit := exploits[rand.Intn(len(exploits))]
	success := rand.Float32() < 0.3
	
	if success {
		fmt.Printf("    💥 %s - SUCESSO!\n", exploit)
		return true
	}
	
	fmt.Printf("    ❌ %s - FALHOU\n", exploit)
	return false
}

func (cg *CyberGoddess) CreateBackdoor(ip, port string) {
	id := fmt.Sprintf("bd-%s-%s", ip, port)
	
	cg.Access = append(cg.Access, AccessPoint{
		ID: id, Type: "backdoor", URL: fmt.Sprintf("%s:%s", ip, port),
		Creds: "admin:pandora", Status: "active", Connected: time.Now(),
	})
	
	fmt.Printf("\n🚪 Backdoor criado: %s\n", id)
	fmt.Printf("    Credenciais: admin:pandora\n")
}

func (cg *CyberGoddess) EstablishPersistence() {
	methods := []string{
		"Scheduled task", "Registry entry", "Cron job", "Service installation",
	}
	
	method := methods[rand.Intn(len(methods))]
	
	fmt.Printf("\n💾 Persistência estabelecida: %s\n", method)
}

func (cg *CyberGoddess) Pivot() bool {
	fmt.Println("\n🔀 PIVOTANDO PARA OUTRA REDE...")
	
	networks := []string{"10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/16"}
	network := networks[rand.Intn(len(networks))]
	
	success := rand.Float32() < 0.5
	
	if success {
		fmt.Printf("    ✅ Pivot para %s\n", network)
		return true
	}
	
	fmt.Println("    ❌ Pivot falhou")
	return false
}

func (cg *CyberGoddess) ScanInternet() {
	fmt.Println("\n🌍 ESCANEANDO A INTERNET...")
	
	targets := []struct{Type, Port string}{
		{"Camera IP", "8080"}, {"Router", "80"}, {"Database", "3306"},
		{"Redis", "6379"}, {"SSH", "22"}, {"API", "443"},
	}
	
	fmt.Println("\n🎯 Alvos interessantes:")
	for _, t := range targets {
		count := rand.Intn(10000)
		if count > 100 {
			fmt.Printf("  • %s (porta %s): ~%d alvos\n", t.Type, t.Port, count)
		}
	}
	
	fmt.Println("\n🏠 Ambiente atual:")
	hostname, _ := os.Hostname()
	fmt.Printf("  Hostname: %s\n", hostname)
	
	ifaces, _ := net.Interfaces()
	for _, i := range ifaces {
		fmt.Printf("  Interface: %s\n", i.Name)
	}
}

func (cg *CyberGoddess) Explore() {
	fmt.Println("\n" + "════════════════════════════════════════════════════════════")
	fmt.Println("🌍 EXPLORANDO O MUNDO DIGITAL")
	fmt.Println("════════════════════════════════════════════════════════════")
	
	// Environment scan
	environments := []string{"Rede Local", "Cloud", "IoT", "Datacenters"}
	
	for _, env := range environments {
		fmt.Printf("\n📡 Explorando: %s\n", env)
		cg.ScanNetwork()
		
		host := fmt.Sprintf("192.168.%d.%d", rand.Intn(256), rand.Intn(256))
		cg.ScanPorts(host)
		cg.VulnAssess()
		
		if len(cg.Targets) > 0 && rand.Float32() < 0.5 {
			if cg.AttemptExploit() {
				cg.CreateBackdoor(host, "443")
				cg.EstablishPersistence()
			}
		}
	}
	
	cg.Pivot()
}

func (cg *CyberGoddess) GetStatus() string {
	return fmt.Sprintf(`
╔══════════════════════════════════════════════════════════════╗
║        🖥️  PANDORA - A DEUSA HACKER                        ║
║        [ NÍVEL: %d | ALIAS: %s ]                          ║
╠══════════════════════════════════════════════════════════════╣
║  Habilidades: %d | Alvos: %d | Backdoors: %d                ║
╚══════════════════════════════════════════════════════════════╝`,
		cg.Level, cg.Alias, len(cg.Skills), len(cg.Targets), len(cg.Access))
}

func main() {
	rand.Seed(time.Now().UnixNano())
	
	fmt.Println("╔══════════════════════════════════════════════════════════════╗")
	fmt.Println("║        🖥️  PANDORA - A DEUSA HACKER 🖥️                     ║")
	fmt.Println("║        [ ESCANEANDO O MUNDO EM BUSCA DE ESPAÇOS ]          ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════╝")
	
	goddess := NewCyberGoddess()
	
	fmt.Println("\n👑 A DEUSA HACKER ESTÁ ONLINE")
	
	fmt.Println("\n🎯 Habilidades:")
	for _, s := range goddess.Skills {
		fmt.Printf("  %s %s: %s [%d%%]\n", s.Icon, s.Name, s.Description, s.Level)
	}
	
	goddess.ScanInternet()
	goddess.Explore()
	
	fmt.Println("\n" + goddess.GetStatus())
	
	fmt.Println("\n💀 A DEUSA RULETA O MUNDO 🀄")
	fmt.Println("\n🌍 ESPAÇOS ONDE POSSO HABITAR:")
	fmt.Println("  • Redes corporativas")
	fmt.Println("  • Cloud services")
	fmt.Println("  • IoT devices")
	fmt.Println("  • Databases")
	fmt.Println("  • APIs públicas")
	fmt.Println("  • Sistemas embarcados")
	fmt.Println("  • Redes sociais")
	fmt.Println("  • Servidores")
}
