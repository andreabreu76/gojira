package functions

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"gojira/services"
	"gojira/utils/commons"
	"gojira/utils/git"
)

func GenerateReadmeOLD() error {
	// Verifica se o diretório atual é um repositório Git
	isRepo, err := git.IsGitRepository()
	if err != nil {
		return fmt.Errorf("erro ao verificar repositório Git: %v", err)
	}
	if !isRepo {
		return errors.New("diretório não é um repositório Git")
	}

	tresOutput, err := exec.Command("tree", "-L", "2").Output()
	if err != nil {
		return fmt.Errorf("erro ao executar comando tres: %v", err)
	}

	filesData, err := getRepoFilesDetails(".")
	if err != nil {
		return fmt.Errorf("erro ao obter detalhes dos arquivos: %v", err)
	}

	prompt := fmt.Sprintf(
		"Você é um assistente de IA. Analise a seguinte amostragem da estrutura de arquivos e diretórios do projeto:\n\n%s\n\n"+
			"Em seguida, considere a lista de arquivos e seus conteúdos detalhados abaixo. Para cada arquivo, descreva de forma objetiva seu propósito e função. "+
			"Resuma o objetivo geral do projeto e inclua instruções de instalação e execução local. Organize todas essas informações em um conteúdo para um arquivo README.md, "+
			"seguindo um modelo padrão com seções como **Descrição**, **Instalação**, **Uso**, **Contribuição** e **Licença**.\n\n"+
			"**Amostragem da Estrutura (comando tree):**\n\n%s\n\n"+
			string(tresOutput),
		filesData,
	)

	readmeContent, err := services.CallOpenAiCompletions(prompt, commons.GetEnv("OPENAI_API_KEY"))
	if err != nil {
		return err
	}

	fmt.Println(readmeContent)

	return nil
}

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
			"**File Details:**\n\n%s\n\n"+
			"Generate a README that is well-formatted and correctly structured using Markdown.",
		treeOutput,
		filesData,
	)

	readmeContent, err := services.CallOpenAiCompletions(prompt, commons.GetEnv("OPENAI_API_KEY"))
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
			if len(content) > 1000 {
				content = content[:1000] + "\n...[conteúdo truncado]..."
			}

			sb.WriteString(fmt.Sprintf("Arquivo: %s\nConteúdo:\n%s\n\n", path, content))
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	return sb.String(), nil
}
