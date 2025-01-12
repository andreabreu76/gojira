package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

var version = "dev"

func main() {
	LoadEnv()

	var rootCmd = &cobra.Command{
		Use:     "gojira",
		Short:   "Uma ferramenta CLI de uso interno para criar descrições de tarefas para o Jira",
		Version: version,
		RunE: func(cmd *cobra.Command, args []string) error {
			title, _ := cmd.Flags().GetString("title")
			taskType, _ := cmd.Flags().GetString("type")
			briefDesc, _ := cmd.Flags().GetString("description")

			if taskType == "Commit" {
				isRepo, err := isGitRepository()
				if err != nil {
					return err
				}
				if !isRepo {
					return errors.New("o diretório atual não é um repositório Git")
				}

				fmt.Println("Repositório Git detectado. Preparando diffs...")

				branch, err := getBranchName()
				if err != nil {
					return err
				}

				diffs, err := getGitDiff()
				if err != nil {
					return err
				}

				if diffs != nil {
					commitMessage, err := generateCommitMessage(diffs, branch)
					if err != nil {
						return err
					}
					fmt.Println("\nMensagem de commit sugerida:\n")
					fmt.Println(commitMessage)
				}
				return nil
			}

			if taskType != "Commit" && title == "" {
				return errors.New("o título da tarefa não pode estar vazio")
			}
			if taskType != "EPICO" && taskType != "BUG" && taskType != "TASK" {
				return errors.New("tipo de tarefa inválido. Use EPICO, BUG ou TASK")
			}

			model := getModel(taskType)

			prompt := fmt.Sprintf("Crie uma descrição detalhada de uma tarefa do tipo %s com o título '%s'. %s "+
				"Baseando-se no modelo: %s os testes e informações para o time de infra são opcionais",
				strings.ToUpper(taskType), title, briefDesc, model)

			response, err := CallOpenAI(prompt)
			if err != nil {
				return err
			}

			fmt.Println("\nDescrição gerada:")
			fmt.Println(response)

			err = clipboard.WriteAll(response)
			if err != nil {
				return errors.New("não foi possível copiar para o clipboard")
			}
			fmt.Println("\nA descrição foi copiada para o clipboard!")
			return nil
		},
	}

	rootCmd.Flags().StringP("title", "t", "", "Título da tarefa (opcional somente quando o tipo for Commit)")
	rootCmd.Flags().StringP("type", "y", "", "Tipo da tarefa: EPICO, BUG, TASK ou Commit (obrigatório)")
	rootCmd.Flags().StringP("description", "d", "", "Descrição breve da tarefa (opcional)")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func getModel(taskType string) string {
	switch strings.ToUpper(taskType) {
	case "EPICO", "TASK":
		return `
Objetivo:
Como:
Critérios de Aceite:
Testes:
Informações para o time de Infra:
Outras observações:
`
	case "BUG":
		return `
Resumo

Título do bug:
ID do bug:
Data de identificação:
Cliente:
Autor do ticket:

Descrição:
Passos para reproduzir:
Comportamento esperado:
Comportamento real:
Capturas de tela:

Ambiente:
Sistema operacional:
Navegador:
Versão do App/Dashboard/Emissor:
Dispositivo:

Informações para o time de Infra:

Outras observações:
`
	default:
		return ""
	}
}

func CallOpenAI(prompt string) (string, error) {
	apiKey := GetEnv("OPENAI_API_KEY")
	if apiKey == "" {
		return "", errors.New("OPENAI_API_KEY não definido no arquivo .env")
	}

	url := "https://api.openai.com/v1/chat/completions"
	body := map[string]interface{}{
		"model":    "gpt-4",
		"messages": []map[string]string{{"role": "user", "content": prompt}},
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 120 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("Erro ao fechar o corpo da resposta")
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("falha na chamada à API: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if choices, ok := result["choices"].([]interface{}); ok && len(choices) > 0 {
		firstChoice := choices[0].(map[string]interface{})
		return firstChoice["message"].(map[string]interface{})["content"].(string), nil
	}

	return "", errors.New("resposta inesperada da API")
}

func LoadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Aviso: Arquivo .env não encontrado. Verificando variáveis de ambiente do sistema...")
	}
}

func GetEnv(key string) string {
	value := os.Getenv(key)
	if value != "" {
		return value
	}

	fmt.Printf("Aviso: A variável %s não está definida.\n", key)
	return ""
}

func isGitRepository() (bool, error) {
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	cmd.Stderr = nil
	err := cmd.Run()
	if err != nil {
		return false, errors.New("o diretório atual não foi identificado como um repositório Git")
	}
	return true, nil
}

func getBranchName() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return "", errors.New("erro ao obter o nome da branch")
	}
	return strings.TrimSpace(string(output)), nil
}

func getGitDiff() (map[string]string, error) {
	ignoredFiles := getIgnoredFiles()

	cmd := exec.Command("git", "status", "--porcelain")
	output, err := cmd.Output()
	if err != nil {
		return nil, errors.New("erro ao executar git status")
	}

	modifiedFiles := []string{}
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, " M") || strings.HasPrefix(line, "A ") || strings.HasPrefix(line, "?? ") {
			file := strings.TrimSpace(line[3:])
			if !isIgnored(file, ignoredFiles) {
				modifiedFiles = append(modifiedFiles, file)
			}
		}
	}

	if len(modifiedFiles) == 0 {
		return nil, errors.New("nenhum arquivo modificado ou não rastreado encontrado")
	}

	fmt.Printf("Arquivos detectados: %v\n\n", modifiedFiles)

	diffs := make(map[string]string)
	for _, file := range modifiedFiles {
		cmd = exec.Command("git", "diff", file)
		diffOutput, err := cmd.Output()
		if err != nil {
			return nil, fmt.Errorf("erro ao obter diff para o arquivo %s", file)
		}
		diffs[file] = string(diffOutput)
	}

	if len(diffs) == 0 {
		return nil, errors.New("nenhuma diferença encontrada nos arquivos selecionados")
	}

	return diffs, nil
}

func generateCommitMessage(diffs map[string]string, branch string) (string, error) {
	prompt := "Você é um assistente de IA. Analise os seguintes diffs do Git e, considerando o contexto da branch atual " +
		"e as regras do Git Flow, forneça uma mensagem de commit objetiva, clara e sucinta em inglês dos EUA, utilizando " +
		"gitemoji que reflita a tarefa somente no titulo. Certifique-se de que a mensagem:\n- Reflita as alterações feitas.\n" +
		"- Alinhe-se com o propósito da branch (por exemplo, feature, bugfix, hotfix, etc.).\n- Use convenções padrão do " +
		"Git Flow para mensagens de commit.:\n\n"
	for file, diff := range diffs {
		prompt += fmt.Sprintf("Branch: %s\nFile: %s\nChanges:\n%s\n\n", branch, file, diff)
	}

	return CallOpenAI(prompt)
}

func getIgnoredFiles() []string {
	filePath := ".gitignore"
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return []string{}
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Erro ao ler o .gitignore: %v\n", err)
		return []string{}
	}

	lines := strings.Split(string(content), "\n")
	var ignoredPatterns []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "#") {
			ignoredPatterns = append(ignoredPatterns, line)
		}
	}
	return ignoredPatterns
}

func isIgnored(file string, ignoredPatterns []string) bool {
	for _, pattern := range ignoredPatterns {
		matched, err := filepath.Match(pattern, file)
		if err != nil {
			fmt.Printf("Erro ao processar o padrão %s: %v\n", pattern, err)
			continue
		}
		if matched {
			return true
		}
	}
	return false
}
