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
	defaultGitName    = "Renato S."
	defaultGitEmail   = "renato.souza@corporate.com.br"
	changeLogFileName = "change_log.txt"
)

type GCommitConfig struct {
	Name  string
	Email string
}

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

func getStagedFiles() []string {
	files, err := exec.Command("git", "diff", "--cached", "--name-only").Output()
	if err != nil {
		log.Fatalf("Erro ao obter arquivos staged: %v", err)
	}

	return strings.Fields(string(files))
}

func getGitAddCandidates() []string {
	out, err := exec.Command("git", "add", "--dry-run", ".").CombinedOutput()
	if err != nil {
		log.Fatalf("Erro ao simular git add .: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	candidates := make([]string, 0, len(lines))

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "add '") && strings.HasSuffix(line, "'") {
			candidates = append(candidates, strings.TrimSuffix(strings.TrimPrefix(line, "add '"), "'"))
			continue
		}

		fields := strings.Fields(line)
		if len(fields) > 0 {
			candidates = append(candidates, fields[len(fields)-1])
		}
	}

	return candidates
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

func loadGCommitConfig() GCommitConfig {
	configPath := ensureGitConfigFile()

	file, err := os.Open(configPath)
	if err != nil {
		log.Fatalf("Erro ao abrir %s: %v", configPath, err)
	}
	defer file.Close()

	config := GCommitConfig{
		Name:  defaultGitName,
		Email: defaultGitEmail,
	}

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
				config.Name = value
			}
		case "email":
			if value != "" {
				config.Email = value
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Erro ao ler %s: %v", configPath, err)
	}

	return config
}

func configureGitIdentity(config GCommitConfig) {
	name := config.Name
	email := config.Email

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
	fileName := changeLogFileName
	separator := "==================\n\n"
	date := time.Now().Format("02/01/2006 às 15h04")
	author := getLastCommitAuthor()

	entry := fmt.Sprintf("Change Log\n\nChange: %s\nDate: %s\nAuthor: %s\n", commitMessage, date, author)

	existingContent, err := os.ReadFile(fileName)
	if err != nil && !os.IsNotExist(err) {
		log.Fatalf("Erro ao ler %s: %v", fileName, err)
	}

	newContent := entry
	if len(existingContent) > 0 {
		newContent += "\n" + separator + string(existingContent)
	}

	if err := os.WriteFile(fileName, []byte(newContent), 0644); err != nil {
		log.Fatalf("Erro ao escrever no %s: %v", fileName, err)
	}

	fmt.Printf("📝 Change log atualizado em '%s'\n", fileName)
}

func isTerraformFile(name string) bool {
	switch {
	case strings.HasSuffix(name, ".tf"):
		return true
	case strings.HasSuffix(name, ".tfvars"):
		return true
	case strings.HasSuffix(name, ".tf.json"):
		return true
	case strings.HasSuffix(name, ".tfvars.json"):
		return true
	case name == "terragrunt.hcl":
		return true
	default:
		return false
	}
}

func gitignoreHasEntry(content string, candidates []string) bool {
	lines := strings.Split(content, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "!") {
			continue
		}

		for _, candidate := range candidates {
			if strings.Contains(line, candidate) {
				return true
			}
		}
	}

	return false
}

func containsProtectedTerraformPath(path string) bool {
	return strings.Contains(path, "/.terraform/") ||
		strings.HasPrefix(path, ".terraform/") ||
		strings.Contains(path, "/.terragrunt/") ||
		strings.HasPrefix(path, ".terragrunt/") ||
		strings.Contains(path, "/.terragrunt-cache/") ||
		strings.HasPrefix(path, ".terragrunt-cache/")
}

func ensureTerraformGitignoreEntries() {
	candidates := getGitAddCandidates()
	shouldProtect := false

	for _, candidate := range candidates {
		if containsProtectedTerraformPath(candidate) || isTerraformFile(filepath.Base(candidate)) {
			shouldProtect = true
			break
		}
	}

	if !shouldProtect {
		return
	}

	existingContent, err := os.ReadFile(".gitignore")
	if err != nil && !os.IsNotExist(err) {
		log.Fatalf("Erro ao ler .gitignore: %v", err)
	}

	gitignore := string(existingContent)
	entriesToAdd := make([]string, 0, 3)

	if !gitignoreHasEntry(gitignore, []string{".terraform"}) {
		entriesToAdd = append(entriesToAdd, ".terraform/")
	}
	if !gitignoreHasEntry(gitignore, []string{".terragrunt"}) {
		entriesToAdd = append(entriesToAdd, ".terragrunt/")
	}
	if !gitignoreHasEntry(gitignore, []string{".terragrunt-cache"}) {
		entriesToAdd = append(entriesToAdd, ".terragrunt-cache/")
	}

	if len(entriesToAdd) == 0 {
		return
	}

	newContent := gitignore
	if strings.TrimSpace(newContent) != "" && !strings.HasSuffix(newContent, "\n") {
		newContent += "\n"
	}
	if strings.TrimSpace(newContent) != "" {
		newContent += "\n"
	}
	newContent += strings.Join(entriesToAdd, "\n") + "\n"

	if err := os.WriteFile(".gitignore", []byte(newContent), 0644); err != nil {
		log.Fatalf("Erro ao atualizar .gitignore: %v", err)
	}

	fmt.Printf("🛡 .gitignore atualizado automaticamente com proteções de Terraform/Terragrunt: %s\n", strings.Join(entriesToAdd, ", "))
}

func ensureNoForbiddenTerraformPathsStaged() {
	stagedFiles := getStagedFiles()
	for _, file := range stagedFiles {
		if containsProtectedTerraformPath(file) {
			log.Fatalf("Arquivo sensível de Terraform/Terragrunt staged detectado: %s. Revise seu .gitignore e remova esses arquivos do stage antes de continuar.", file)
		}
	}
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
	config := loadGCommitConfig()
	configureGitIdentity(config)
	ensureTerraformGitignoreEntries()

	runGitAdd()

	ensureNoForbiddenTerraformPathsStaged()

	fileList := getStagedFiles()
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
		runGitAddFile(changeLogFileName)
		runGitCommit("chore: update change log [skip ci]")
		fmt.Println("✅ Commit do change log realizado com sucesso!")
		runGitPush(branch)
	}
}
