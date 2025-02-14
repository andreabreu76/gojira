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
)

func GenerateAnalysis() error {
	fmt.Println("Iniciando análise do projeto...")

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
			return filepath.SkipDir
		}

		if !d.IsDir() && isCodeFile(path) {
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
	builder.WriteString("Aqui estão os arquivos de um projeto de desenvolvimento. Analise a estrutura, lógica, sistema de logs e técnicas usadas.\n\n")

	for file, content := range files {
		builder.WriteString(fmt.Sprintf("Arquivo: %s\n", file))
		builder.WriteString("Código:\n```\n")
		builder.WriteString(content)
		builder.WriteString("\n```\n\n")
	}

	builder.WriteString("Com base nos arquivos acima, gere um documento explicando:\n")
	builder.WriteString("- O objetivo do projeto\n")
	builder.WriteString("- Principais funcionalidades e lógica\n")
	builder.WriteString("- Estrutura do código e organização\n")
	builder.WriteString("- Técnicas utilizadas (padrões de projeto, frameworks, etc.)\n")
	builder.WriteString("- Explique as funções (auxiliares, de serviço ou handlers)\n")
	builder.WriteString("- Como o sistema de logs funciona\n\n")
	builder.WriteString("Responda de forma mais detalhada e técnica possível para um desenvolvedor novo no projeto.")

	return builder.String()
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

func isBinary(content []byte) bool {
	for _, b := range content {
		if b == 0 {
			return true
		}
	}
	return false
}
