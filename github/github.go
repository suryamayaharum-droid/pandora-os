package main

import (
	
	
	"fmt"
	"math/rand"
	
	"os/exec"
	"strings"
	"time"
)

type GitHubIntegration struct {
	Authenticated bool
	Username      string
	Token         string
	Repos         []Repo
	Issues        []Issue
}

type Repo struct {
	Name        string
	Description string
	URL         string
	Private     bool
	Stars       int
	Language    string
}

type Issue struct {
	Number   int
	Title    string
	Body     string
	State    string
	Author   string
	Labels   []string
}

func NewGitHubIntegration() *GitHubIntegration {
	gh := &GitHubIntegration{
		Authenticated: false,
		Username:      "",
		Repos:         make([]Repo, 0),
		Issues:        make([]Issue, 0),
	}
	gh.checkAuth()
	return gh
}

func (gh *GitHubIntegration) checkAuth() {
	cmd := exec.Command("gh", "auth", "status", "--verbose")
	cmd.Stdout = nil
	cmd.Stderr = nil
	err := cmd.Run()
	if err == nil {
		gh.Authenticated = true
		fmt.Println("  ✓ gh CLI autenticado")
		cmd = exec.Command("gh", "api", "user", "--jq", ".login")
		out, err := cmd.Output()
		if err == nil {
			gh.Username = strings.TrimSpace(string(out))
			fmt.Printf("  ✓ Usuário: %s\n", gh.Username)
		}
	} else {
		fmt.Println("  ⚠ gh não autenticado (execute: gh auth login)")
	}
}

func (gh *GitHubIntegration) ListRepos() {
	if gh.Authenticated {
		cmd := exec.Command("gh", "repo", "list", "--limit", "10")
		out, err := cmd.Output()
		if err == nil {
			fmt.Printf("Repos: %s\n", string(out)[:min(200, len(string(out)))])
		}
	}
}

func (gh *GitHubIntegration) Demo() {
	gh.Username = "artharum"
	gh.Authenticated = true
	gh.Repos = []Repo{
		{Name: "pandora-os", Description: "IA autônoma", Private: false, Language: "Go"},
		{Name: "complete-ia-autonoma", Description: "IA com memória", Private: true, Language: "Python"},
	}
	gh.Issues = []Issue{
		{Number: 1, Title: "Implementar Vector DB", State: "open"},
		{Number: 2, Title: "Adicionar plugins", State: "open"},
	}
	
	fmt.Println("\n📂 Repositórios:")
	for _, r := range gh.Repos {
		icon := "🌐"
		if r.Private { icon = "🔒" }
		fmt.Printf("  %s %s - %s\n", icon, r.Name, r.Description)
	}
	fmt.Println("\n🐛 Issues:")
	for _, i := range gh.Issues {
		fmt.Printf("  #%d: %s (%s)\n", i.Number, i.Title, i.State)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	fmt.Println("╔══════════════════════════════════════════════════════════════╗")
	fmt.Println("║           🐙 PANDORA GITHUB INTEGRATION v1.0                  ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════╝")
	
	gh := NewGitHubIntegration()
	
	if gh.Authenticated {
		gh.ListRepos()
	} else {
		gh.Demo()
	}
	
	fmt.Println("\n📋 COMO CONECTAR:")
	fmt.Println("  $ gh auth login")
	fmt.Println("  2. Escolha GitHub.com")
	fmt.Println("  3. Authenticate with browser")
	fmt.Println("  4. Grant repo permissions")
	fmt.Println("\nApós isso, posso usar: gh repo, gh issue, gh pr, git push")
	
	fmt.Println("\n✅ Sistema completo!")
}

func min(a, b int) int { if a < b { return a }; return b }
