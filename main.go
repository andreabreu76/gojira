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

			if title == "" {
				return errors.New("o título da tarefa não pode estar vazio")
			}
			if taskType != "EPICO" && taskType != "BUG" && taskType != "TASK" {
				return errors.New("tipo de tarefa inválido. Use EPICO, BUG ou TASK")
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
