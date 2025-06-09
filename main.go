package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	openai "github.com/sashabaranov/go-openai"
	"github.com/joho/godotenv"
)

func runGitAdd() {
	cmd := exec.Command("git", "add", ".")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatalf("Erro ao executar git add .: %v", err)
	}
	fmt.Println("✅ git add . executado com sucesso.")
}

func getCurrentBranch() string {
	out, err := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD").Output()
	if err != nil {
		log.Fatalf("Erro ao obter nome da branch atual: %v", err)
	}
	return strings.TrimSpace(string(out))
}

func runGitPush(branch string) {
	cmd := exec.Command("git", "push", "origin", branch)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatalf("Erro ao fazer push para a branch %s: %v", branch, err)
	}
	fmt.Printf("🚀 Push realizado para 'origin/%s'\n", branch)
}

func loadAPIKey() string {
	_ = godotenv.Load() // Tenta carregar .env, mas não interrompe se falhar

	apiKey := os.Getenv("api_key")
	if apiKey == "" {
		log.Fatal("❌ API key não encontrada. Defina no arquivo .env ou como variável de ambiente (api_key).")
	}
	return apiKey
}

func main() {
	pushAfterCommit := len(os.Args) > 1 && os.Args[1] == "-a"
	apiKey := loadAPIKey()

	runGitAdd()

	files, err := exec.Command("git", "diff", "--cached", "--name-only").Output()
	if err != nil {
		log.Fatalf("Erro ao obter arquivos staged: %v", err)
	}
	fileList := strings.Fields(string(files))
	if len(fileList) == 0 {
		fmt.Println("Nenhum arquivo staged encontrado com `git add`.")
		os.Exit(0)
	}

	diff, err := exec.Command("git", "diff", "--cached", "--unified=1").Output()
	if err != nil {
		log.Fatalf("Erro ao obter diff: %v", err)
	}

	summary := string(diff)
	if len(summary) > 1000 {
		summary = summary[:1000]
	}

	client := openai.NewClient(apiKey)
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:     "gpt-3.5-turbo",
			MaxTokens: 40,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "Você é um assistente que gera mensagens de commit curtas e descritivas.",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: "Gere uma mensagem de commit curta baseada nessas mudanças:\n\n" + summary,
				},
			},
		},
	)
	if err != nil {
		log.Fatalf("Erro ao gerar commit message: %v", err)
	}

	commitMessage := strings.TrimSpace(resp.Choices[0].Message.Content)
	fmt.Printf("\n📦 Commit message gerada:\n%s\n", commitMessage)

	cmd := exec.Command("git", "commit", "-m", commitMessage)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatalf("Erro ao fazer commit: %v", err)
	}

	fmt.Println("✅ Commit realizado com sucesso!")

	if pushAfterCommit {
		branch := getCurrentBranch()
		runGitPush(branch)
	}
}
