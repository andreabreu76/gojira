package functions

import (
	"errors"
	"fmt"
	"gojira/services"
	"gojira/utils/commons"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func GenerateAnalysis() error {
	projectName := getProjectName()
	files, err := getProjectFiles(".")
	if err != nil {
		return fmt.Errorf("erro ao obter arquivos do projeto: %w", err)
	}

	if len(files) == 0 {
		return errors.New("nenhum arquivo relevante encontrado no projeto")
	}

	fileContents, err := readProjectFiles(files)
	if err != nil {
		return fmt.Errorf("erro ao ler arquivos do projeto: %w", err)
	}

	prompt := buildAnalysisPrompt(fileContents, projectName)
	prompt = minifyPrompt(prompt)

	if err := logPrompt(prompt); err != nil {
		return fmt.Errorf("erro ao gravar log da análise: %w", err)
	}

	response, err := services.CallOpenAiCompletions(prompt, commons.GetEnv("OPENAI_API_KEY"))
	if err != nil {
		return fmt.Errorf("erro ao obter resposta da OpenAI: %w", err)
	}

	fmt.Println("\n=== Análise do Projeto ===")
	fmt.Println(response)

	return nil
}

func getProjectFiles(root string) ([]string, error) {
	var files []string

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() && (strings.Contains(path, ".git") ||
			strings.Contains(path, "node_modules") ||
			strings.Contains(path, "vendor")) {
			return filepath.SkipDir
		}

		if !d.IsDir() {
			files = append(files, path)
		}

		return nil
	})

	return files, err
}

func readProjectFiles(files []string) (map[string]string, error) {
	fileContents := make(map[string]string)

	for _, file := range files {
		if strings.HasSuffix(file, ".csproj") {
			log.Printf("ignorando arquivo .csproj: %s", file)
			continue
		}

		if strings.Contains(file, "/bin/") || strings.Contains(file, "/obj/") ||
			strings.HasPrefix(file, "bin/") || strings.HasPrefix(file, "obj/") {
			log.Printf("ignorando diretório bin/ ou obj/: %s", file)
			continue
		}

		if strings.Contains(file, "/.idea/") || strings.Contains(file, "/.vscode/") ||
			strings.HasPrefix(file, ".idea/") || strings.HasPrefix(file, ".vscode/") {
			log.Printf("ignorando diretório .idea/ ou .vscode/: %s", file)
			continue
		}

		content, err := os.ReadFile(file)
		if err != nil {
			log.Printf("erro ao ler arquivo %s: %v", file, err)
			continue
		}

		if isBinary(content) || len(content) > 500*1024 {
			log.Printf("ignorando arquivo binário ou muito grande: %s", file)
			continue
		}

		fileContents[file] = string(content)
	}

	return fileContents, nil
}

func buildAnalysisPrompt(files map[string]string, projectName string) string {
	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("# Análise do Projeto: %s\n\n", projectName))
	builder.WriteString("Você é um assistente especializado em análise de código-fonte.\n")
	builder.WriteString("Aqui estão os arquivos de um projeto de desenvolvimento. Analise detalhadamente a estrutura, lógica, " +
		"sistema de logs, técnicas utilizadas e a configuração de CI/CD.\n\n")

	ciCdFiles := []string{}

	for file, content := range files {
		builder.WriteString(fmt.Sprintf("## Arquivo: %s\n", file))
		builder.WriteString("```yaml\n")
		builder.WriteString(content)
		builder.WriteString("\n```\n\n")

		if strings.HasPrefix(file, ".github/workflows/") {
			ciCdFiles = append(ciCdFiles, file)
		}
	}

	builder.WriteString("## Relatório de Análise\n")
	builder.WriteString("Com base nos arquivos acima, gere um documento explicando:\n")
	builder.WriteString("- O objetivo do projeto\n")
	builder.WriteString("- Principais funcionalidades e lógica\n")
	builder.WriteString("- Estrutura do código e organização\n")
	builder.WriteString("- Técnicas utilizadas (padrões de projeto, frameworks, etc.)\n")
	builder.WriteString("- Explicação detalhada das funções (auxiliares, de serviço ou handlers)\n")
	builder.WriteString("- Como o sistema de logs funciona\n\n")

	if len(ciCdFiles) > 0 {
		builder.WriteString("## Análise de CI/CD\n")
		builder.WriteString("O projeto contém arquivos de configuração de CI/CD. Avalie como o pipeline está estruturado e identifique:\n")
		builder.WriteString("- Ferramentas utilizadas (GitHub Actions, CircleCI, etc.)\n")
		builder.WriteString("- Passos do pipeline (build, test, deploy)\n")
		builder.WriteString("- Melhorias sugeridas\n\n")
	}

	builder.WriteString("Responda de forma mais detalhada e técnica possível para um desenvolvedor novo no projeto.")

	return builder.String()
}

func logPrompt(prompt string) error {
	logDir := filepath.Join(os.Getenv("HOME"), ".log")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("erro ao criar diretório de logs: %w", err)
	}

	timestamp := time.Now().Format("20060102-150405")
	logFile := filepath.Join(logDir, fmt.Sprintf("%s-gojira-analysis.log", timestamp))

	file, err := os.Create(logFile)
	if err != nil {
		return fmt.Errorf("erro ao criar arquivo de log: %w", err)
	}
	defer file.Close()

	_, err = file.WriteString(prompt)
	return err
}

func isBinary(content []byte) bool {
	for _, b := range content {
		if b == 0 {
			return true
		}
	}
	return false
}

func getProjectName() string {
	wd, err := os.Getwd()
	if err != nil {
		log.Printf("Erro ao obter diretório atual: %v", err)
		return "Projeto Desconhecido"
	}

	return filepath.Base(wd)
}

func minifyPrompt(prompt string) string {
	prompt = strings.ReplaceAll(prompt, "\n\n", "\n") // Remove linhas vazias extras
	prompt = strings.ReplaceAll(prompt, "\t", " ")    // Substitui tabulações por espaços únicos
	prompt = strings.TrimSpace(prompt)                // Remove espaços extras do começo e fim
	return prompt
}
