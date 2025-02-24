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

	builder.WriteString("Você é um assistente sênior especializado em análise de código-fonte e práticas de desenvolvimento.\n")
	builder.WriteString("Seu objetivo é analisar com profundidade cada trecho de código apresentado.\n")
	builder.WriteString("Explique cada função, cada parâmetro, a lógica, a arquitetura e as razões por trás das escolhas feitas.\n")
	builder.WriteString("Inclua explicações sobre como tudo se conecta, apontando pontos fortes e oportunidades de melhoria.\n\n")

	builder.WriteString("A seguir estão os arquivos de um projeto de desenvolvimento que passará por migração para uma versão mais recente.\n")
	builder.WriteString("Analise detalhadamente:\n")
	builder.WriteString("- Estrutura do código e organização de pastas\n")
	builder.WriteString("- Cada função e seu propósito\n")
	builder.WriteString("- Módulos/libraries/frameworks utilizados e por quê\n")
	builder.WriteString("- Sistema de logs (como funciona e se pode ser melhorado)\n")
	builder.WriteString("- Configuração de CI/CD (arquivos, pipelines, etapas, melhorias possíveis)\n\n")

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
	builder.WriteString("Crie um relatório completo para um desenvolvedor novo no time, abrangendo:\n")
	builder.WriteString("1. Objetivo do projeto e seu contexto\n")
	builder.WriteString("2. Principais funcionalidades e como estão implementadas\n")
	builder.WriteString("3. Detalhes de cada arquivo relevante (classes, funções, parâmetros, objetos trocados)\n")
	builder.WriteString("4. Técnicas, padrões de projeto e frameworks utilizados\n")
	builder.WriteString("5. Sistema de logs (como está configurado, pontos de melhoria)\n")
	builder.WriteString("6. Possíveis pontos de refatoração para a migração da versão\n\n")

	if len(ciCdFiles) > 0 {
		builder.WriteString("## Análise de CI/CD\n")
		builder.WriteString("Explique como o pipeline está estruturado e identifique:\n")
		builder.WriteString("- Quais ferramentas de CI/CD são usadas\n")
		builder.WriteString("- As etapas do pipeline (build, testes, deploy)\n")
		builder.WriteString("- Configurações específicas de ambiente\n")
		builder.WriteString("- Possíveis melhorias e otimizações\n\n")
	}

	builder.WriteString("Estruture a resposta de forma clara e técnica, usando exemplos do código sempre que necessário. ")
	builder.WriteString("Se algo não estiver claro no código, proponha soluções ou hipóteses prováveis. ")
	builder.WriteString("Finalize com um resumo das recomendações para a migração.\n")

	prompt := builder.String()
	fmt.Println(prompt)

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
