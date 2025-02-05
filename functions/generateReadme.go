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

func GenerateReadme() error {
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
