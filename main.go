package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	openai "github.com/sashabaranov/go-openai"
	"github.com/joho/godotenv"
)

const (
	defaultGitName  = "Renato S."
	defaultGitEmail = "renato.souza@corporate.com.br"
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

func getLastCommitAuthor() string {
	out, err := exec.Command("git", "log", "-1", "--pretty=%an").Output()
	if err != nil {
		log.Fatalf("Erro ao obter autor do último commit: %v", err)
	}
	return strings.TrimSpace(string(out))
}

func runGitAddFile(fileName string) {
	cmd := exec.Command("git", "add", fileName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatalf("Erro ao executar git add %s: %v", fileName, err)
	}
	fmt.Printf("✅ git add %s executado com sucesso.\n", fileName)
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

func runGitPull(branch string) {
	cmd := exec.Command("git", "pull", "--rebase", "origin", branch)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatalf("Erro ao fazer pull com rebase da branch %s em origin: %v", branch, err)
	}
	fmt.Printf("🔄 Pull com rebase realizado de 'origin/%s'\n", branch)
}

func runGitCommit(message string) {
	cmd := exec.Command("git", "commit", "-m", message)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatalf("Erro ao fazer commit: %v", err)
	}
}

func getGitConfigFilePath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Erro ao obter diretório home do usuário Linux: %v", err)
	}

	return filepath.Join(homeDir, ".gcommit.conf")
}

func ensureGitConfigFile() string {
	configPath := getGitConfigFilePath()

	if _, err := os.Stat(configPath); err == nil {
		return configPath
	} else if !os.IsNotExist(err) {
		log.Fatalf("Erro ao verificar %s: %v", configPath, err)
	}

	content := fmt.Sprintf("name=%s\nemail=%s\n", defaultGitName, defaultGitEmail)
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		log.Fatalf("Erro ao criar %s: %v", configPath, err)
	}

	fmt.Printf("🛠 Arquivo de configuração criado em '%s'\n", configPath)
	return configPath
}

func loadGitIdentity() (string, string) {
	configPath := ensureGitConfigFile()

	file, err := os.Open(configPath)
	if err != nil {
		log.Fatalf("Erro ao abrir %s: %v", configPath, err)
	}
	defer file.Close()

	name := defaultGitName
	email := defaultGitEmail

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		key, value, found := strings.Cut(line, "=")
		if !found {
			continue
		}

		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)

		switch key {
		case "name":
			if value != "" {
				name = value
			}
		case "email":
			if value != "" {
				email = value
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Erro ao ler %s: %v", configPath, err)
	}

	return name, email
}

func configureGitIdentity() {
	name, email := loadGitIdentity()

	nameCmd := exec.Command("git", "config", "user.name", name)
	nameCmd.Stdout = os.Stdout
	nameCmd.Stderr = os.Stderr
	if err := nameCmd.Run(); err != nil {
		log.Fatalf("Erro ao configurar git user.name: %v", err)
	}

	emailCmd := exec.Command("git", "config", "user.email", email)
	emailCmd.Stdout = os.Stdout
	emailCmd.Stderr = os.Stderr
	if err := emailCmd.Run(); err != nil {
		log.Fatalf("Erro ao configurar git user.email: %v", err)
	}

	fmt.Printf("👤 Identidade Git configurada: %s <%s>\n", name, email)
}

func appendChangeLog(commitMessage string) {
	fileName := "change_log.txt"
	separator := "==================\n\n"
	date := time.Now().Format("02/01/2006 às 15h04")
	author := getLastCommitAuthor()

	entry := fmt.Sprintf("Change Log\n\nChange: %s\nDate: %s\nAuthor: %s\n", commitMessage, date, author)

	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("Erro ao abrir %s: %v", fileName, err)
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		log.Fatalf("Erro ao verificar %s: %v", fileName, err)
	}

	if info.Size() > 0 {
		if _, err := file.WriteString("\n" + separator); err != nil {
			log.Fatalf("Erro ao escrever separador no %s: %v", fileName, err)
		}
	}

	if _, err := file.WriteString(entry); err != nil {
		log.Fatalf("Erro ao escrever no %s: %v", fileName, err)
	}

	fmt.Printf("📝 Change log atualizado em '%s'\n", fileName)
}

func loadAPIKey() string {
	_ = godotenv.Load() // Tenta carregar .env, mas não interrompe se falhar

	apiKey := os.Getenv("api_key")
	if apiKey == "" {
		log.Fatal("❌ API key não encontrada. Defina no arquivo .env ou como variável de ambiente (api_key).")
	}
	return apiKey
}

func parseFlags(args []string) (bool, bool) {
	pushAfterCommit := false
	createChangeLog := false

	for _, arg := range args {
		if !strings.HasPrefix(arg, "-") || len(arg) < 2 {
			continue
		}

		for _, flag := range arg[1:] {
			switch flag {
			case 'a':
				pushAfterCommit = true
			case 'c':
				createChangeLog = true
			}
		}
	}

	if createChangeLog {
		pushAfterCommit = true
	}

	return pushAfterCommit, createChangeLog
}

func main() {
	pushAfterCommit, createChangeLog := parseFlags(os.Args[1:])

	apiKey := loadAPIKey()
	configureGitIdentity()

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

	runGitCommit(commitMessage)
	fmt.Println("✅ Commit realizado com sucesso!")

	var branch string

	if pushAfterCommit {
		branch = getCurrentBranch()
		runGitPull(branch)
		runGitPush(branch)
	}

	if createChangeLog {
		appendChangeLog(commitMessage)
		runGitAddFile("change_log.txt")
		runGitCommit("chore: update change log [skip ci]")
		fmt.Println("✅ Commit do change log realizado com sucesso!")
		runGitPush(branch)
	}
}
