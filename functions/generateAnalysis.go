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

	prompt := buildAnalysisPrompt(fileContents)

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

		if d.IsDir() && (strings.Contains(path, ".git") || strings.Contains(path, "node_modules") || strings.Contains(path, "vendor")) {
			if !strings.HasPrefix(path, ".github/workflows") {
				return filepath.SkipDir
			}
		}

		if !d.IsDir() && (isCodeFile(path) || isYamlFile(path)) {
			files = append(files, path)
		}

		return nil
	})

	return files, err
}

func readProjectFiles(files []string) (map[string]string, error) {
	fileContents := make(map[string]string)

	for _, file := range files {
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

func buildAnalysisPrompt(files map[string]string) string {
	var builder strings.Builder

	builder.WriteString("Você é um assistente especializado em análise de código-fonte.\n")
	builder.WriteString("Aqui estão os arquivos de um projeto de desenvolvimento. Analise detalhadamente a estrutura, lógica, " +
		"sistema de logs, técnicas utilizadas e, se presente, a configuração de CI/CD.\n\n")

	ciCdFiles := []string{}

	for file, content := range files {
		builder.WriteString(fmt.Sprintf("Arquivo: %s\n", file))
		builder.WriteString("Código:\n```\n")
		builder.WriteString(content)
		builder.WriteString("\n```\n\n")

		if strings.HasPrefix(file, ".github/workflows/") {
			ciCdFiles = append(ciCdFiles, file)
		}
	}

	builder.WriteString("Com base nos arquivos acima, gere um documento explicando:\n")
	builder.WriteString("- O objetivo do projeto\n")
	builder.WriteString("- Principais funcionalidades e lógica\n")
	builder.WriteString("- Estrutura do código e organização\n")
	builder.WriteString("- Técnicas utilizadas (padrões de projeto, frameworks, etc.)\n")
	builder.WriteString("- Explicação detalhada das funções (auxiliares, de serviço ou handlers)\n")
	builder.WriteString("- Como o sistema de logs funciona\n\n")

	if len(ciCdFiles) > 0 {
		builder.WriteString("### Análise de CI/CD\n")
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
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Printf("erro ao fechar arquivo de log: %v", err)
		}
	}(file)

	_, err = file.WriteString(prompt)
	return err
}

func isCodeFile(path string) bool {
	ext := filepath.Ext(path)
	codeExtensions := []string{".go", ".js", ".ts", ".py", ".java", ".cpp", ".h", ".cs", ".rb", ".php", ".rs"}
	for _, e := range codeExtensions {
		if ext == e {
			return true
		}
	}
	return false
}

func isYamlFile(path string) bool {
	ext := filepath.Ext(path)
	return ext == ".yaml" || ext == ".yml"
}

func isBinary(content []byte) bool {
	for _, b := range content {
		if b == 0 {
			return true
		}
	}
	return false
}
