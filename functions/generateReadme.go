package functions

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"gojira/services/ai"
	"gojira/utils/commons"
	"gojira/utils/git"
)

func GenerateReadme() error {
	isRepo, err := git.IsGitRepository()
	if err != nil {
		return fmt.Errorf("erro ao verificar repositório Git: %v", err)
	}
	if !isRepo {
		return errors.New("diretório não é um repositório Git")
	}

	treeOutput, err := exec.Command("tree", "-L", "2").Output()
	if err != nil {
		return fmt.Errorf("erro ao executar comando tree: %v", err)
	}

	filesData, err := getRepoFilesDetails(".")
	if err != nil {
		return fmt.Errorf("erro ao obter detalhes dos arquivos: %v", err)
	}

	analysisFilesData, err := getAnalysisFiles(".")
	if err != nil {
		return fmt.Errorf("erro ao obter detalhes dos arquivos de análise: %v", err)
	}

	prompt := fmt.Sprintf(
		"You are an AI assistant specialized in technical documentation.\n\n"+
			"Analyze the project structure and generate a well-structured README.md following best practices. "+
			"Ensure the README is written in **US English** and includes the following sections:\n\n"+
			"1. **Project Name** - Name and status.\n"+
			"2. **Description** - Summary of the project's purpose and functionality.\n"+
			"3. **Technologies Used** - List of main technologies.\n"+
			"4. **Project Structure** - Hierarchical representation of files.\n"+
			"5. **Installation** - Step-by-step guide for local setup.\n"+
			"6. **Usage** - Basic usage examples.\n"+
			"7. **API Documentation** - Instructions if there are API endpoints.\n"+
			"8. **Contributing** - Guidelines for contributing to the project.\n"+
			"9. **License** - License type used.\n\n"+
			"Use the project file structure below to generate the correct documentation:\n\n"+
			"**Project Structure:**\n\n%s\n\n"+
			"**File Details:**\n\n%s\n\n",
		treeOutput,
		filesData,
	)

	if analysisFilesData != "" {
		prompt += fmt.Sprintf("**Analysis Files:**\n\n%s\n\n", analysisFilesData)
	}

	prompt += "Generate a README that is well-formatted and correctly structured using Markdown."

	// Carrega configuração
	config, err := commons.LoadConfig()
	if err != nil {
		return fmt.Errorf("erro ao carregar configuração: %w", err)
	}

	// Obtém o provedor de IA configurado
	provider, exists := ai.GetProvider(config.AIProvider)
	if !exists {
		provider = ai.GetDefaultProvider()
	}

	readmeContent, err := provider.GetCompletions(prompt, config.AIModel)
	if err != nil {
		return err
	}

	err = os.WriteFile("README.md", []byte(readmeContent), 0644)
	if err != nil {
		return fmt.Errorf("erro ao salvar README.md: %v", err)
	}

	fmt.Println("README.md gerado com sucesso!")
	return nil
}

func getRepoFilesDetails(baseDir string) (string, error) {
	var sb strings.Builder

	err := filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if strings.HasPrefix(info.Name(), ".") {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if strings.ToLower(info.Name()) == "readme.md" {
			return nil
		}

		if !info.IsDir() {
			data, err := os.ReadFile(path)
			if err != nil {
				return nil
			}

			content := string(data)
			sb.WriteString(fmt.Sprintf("Arquivo: %s\nConteúdo:\n%s\n\n", path, content))
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	return sb.String(), nil
}

func getAnalysisFiles(baseDir string) (string, error) {
	var sb strings.Builder

	entries, err := os.ReadDir(baseDir)
	if err != nil {
		return "", err
	}

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasPrefix(entry.Name(), "analysis_") {
			data, err := os.ReadFile(filepath.Join(baseDir, entry.Name()))
			if err != nil {
				continue
			}

			content := string(data)
			sb.WriteString(fmt.Sprintf("Arquivo: %s\nConteúdo:\n%s\n\n", entry.Name(), content))
		}
	}

	return sb.String(), nil
}
