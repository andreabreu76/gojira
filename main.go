package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

var version = "dev"

type ReadmeData struct {
	ProjectName       string
	BriefDescription  string
	Description       string
	Technologies      string
	ProjectType       string
	UseCases          string
	Prerequisites     string
	RepositoryURL     string
	InstallationSteps string
	Usage             string
	ProjectStructure  string
	ContactInfo       string
}

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

			if title == "" {
				return errors.New("o título da tarefa não pode estar vazio")
			}
			if taskType != "EPICO" && taskType != "BUG" && taskType != "TASK" {
				return errors.New("tipo de tarefa inválido. Use EPICO, BUG ou TASK")
			}

			if taskType == "README" {
				if !isGitProject() {
					return errors.New("o diretório atual não parece um projeto válido")
				}

				InitializeRedis()
				fmt.Println("[INFO] Iniciando varredura do diretório...")

				files := scanDirectory(".")
				for _, file := range files {
					fmt.Printf("[INFO] Analisando %s...\n", file)
					result := analyzeFileWithOpenAI(file)
					err := StoreResult("gojira:"+file, result)
					if err != nil {
						return errors.New("falha ao armazenar análise no Redis")
					}
				}

				results := FetchAllResults("gojira:")
				var combinedResults string
				for file, result := range results {
					combinedResults += fmt.Sprintf("### Arquivo: %s\n%s\n\n", file, result)
				}

				prompt := fmt.Sprintf("Baseado nos arquivos e descrições abaixo, crie um README.md:\n%s", combinedResults)
				response, err := CallOpenAI(prompt)
				if err != nil {
					return errors.New("falha ao gerar README com OpenAI")
				}

				generateReadmeFile(response)
				fmt.Println("[INFO] README_gojira.MD gerado com sucesso!")
			}

			model := getModel(taskType)

			prompt := fmt.Sprintf("Crie uma descrição detalhada de uma tarefa do tipo %s com o título '%s'. %s Baseando-se no modelo: %s os testes e infroações para o time de infra são opcionais",
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

	rootCmd.Flags().StringP("title", "t", "", "Título da tarefa (obrigatório)")
	rootCmd.Flags().StringP("type", "y", "", "Tipo da tarefa: EPICO, BUG ou TASK (obrigatório)")
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
	case "README":
		return `
# Nome do Projeto
> Breve descrição do projeto, o que ele faz e seu propósito.
---
## Sumário
- [Descrição](#descrição)
- [Pré-requisitos](#pré-requisitos)
- [Instalação](#instalação)
- [Uso](#uso)
- [Estrutura do Projeto](#estrutura-do-projeto)
---
## Descrição
Inclua aqui uma explicação clara do que o projeto faz, quem o utiliza e por que ele é importante.
- Tecnologias principais: [Go](https://golang.org), [Node.js](https://nodejs.org/), [MongoDB](https://www.mongodb.com/), etc.
- Tipo de projeto: Biblioteca, API, Serviço, Ferramenta CLI, etc.
- Casos de uso: Explique como o projeto resolve problemas ou automatiza processos.
---
## Pré-requisitos
Liste os requisitos mínimos para rodar ou usar o projeto.
- Linguagem: Go (>= 1.18) ou Node.js (>= 16.x).
- Dependências: Docker, Redis, PostgreSQL, etc.
- Sistemas operacionais suportados: Windows, macOS, Linux.
Exemplo:

---
## Instalação
Inclua instruções claras para configurar o ambiente e instalar as dependências.
### Clone o repositório

### Instalar Dependências
Para projetos em Node.js:
Para projetos em Go:

---
## Uso
Inclua exemplos ou instruções para rodar o projeto.
### Para rodar localmente:
## Estrutura do Projeto
Forneça uma visão geral dos diretórios principais e seus propósitos:

---
## Contato de Apoio e Dúvidas
Fulano da Silva [fulano.silva@yooga.com.br](mailto:fulano.silva@yooga.com.br).
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

	value = os.Getenv(key)
	if value == "" {
		fmt.Printf("Aviso: A variável %s não está definida.\n", key)
	}

	return value
}

func isGitProject() bool {
	cmd := exec.Command("git", "config", "--list")
	output, err := cmd.Output()
	if err != nil {
		return false
	}

	return strings.Contains(string(output), "remote.origin.url")
}

func scanDirectory(root string) []string {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		fmt.Println("Erro ao listar arquivos:", err)
		return nil
	}
	return files
}

func analyzeFileWithOpenAI(filePath string) string {
	content, _ := os.ReadFile(filePath)
	prompt := fmt.Sprintf("Analise o seguinte arquivo de código e descreva suas classes, funções e propósito:\n\n%s", content)

	response, _ := CallOpenAI(prompt)
	return response
}

func summarizeResults(results map[string]string) string {
	var combined string
	for file, result := range results {
		combined += fmt.Sprintf("### Arquivo: %s\n%s\n\n", file, result)
	}

	prompt := fmt.Sprintf("Crie um resumo consolidado com base nas seguintes análises:\n\n%s", combined)
	finalSummary, _ := CallOpenAI(prompt) // Utiliza OpenAI para gerar resumo
	return finalSummary
}

func generateReadmeFile(summary string) {
	data := ReadmeData{
		ProjectName:      "NomeDoProjeto",
		BriefDescription: "Projeto automatizado para análise e documentação.",
		Description:      summary,
		ProjectStructure: "Gerado automaticamente pelo GOJIRA.",
		ContactInfo:      "Fulano da Silva - fulano.silva@yooga.com.br",
	}

	generateReadme(data) // Usa a função existente para gerar README com o template
}

func generateReadme(data ReadmeData) {
	// Template do README
	const readmeTemplate = `# {{.ProjectName}}
> {{.BriefDescription}}

## Descrição
{{.Description}}

## Estrutura do Projeto
{{.ProjectStructure}}

## Contato
{{.ContactInfo}}
`

	tmpl, err := template.New("readme").Parse(readmeTemplate)
	if err != nil {
		log.Fatal("Erro ao processar template:", err)
	}

	file, err := os.Create("README_gojira.MD")
	if err != nil {
		log.Fatal("Erro ao criar arquivo:", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println("Erro ao fechar o arquivo:", err)
		}
	}(file)

	err = tmpl.Execute(file, data)
	if err != nil {
		log.Fatal("Erro ao gerar README:", err)
	}

	fmt.Println("README_gojira.MD gerado com sucesso!")
}
