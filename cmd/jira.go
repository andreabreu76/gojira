package cmd

import (
	"errors"
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/spf13/cobra"
	"gojira/services"
	"gojira/services/ai"
	"gojira/utils/commons"
	"strings"
)

var (
	title      string
	taskType   string
	briefDesc  string
	projectKey string
)

// jiraCmd representa o comando para interagir com o Jira
var jiraCmd = &cobra.Command{
	Use:   "jira",
	Short: "Gera descrições para tarefas do Jira",
	RunE: func(cmd *cobra.Command, args []string) error {
		if taskType != "EPICO" && taskType != "BUG" && taskType != "TASK" {
			return errors.New("tipo de tarefa inválido. Use EPICO, BUG ou TASK")
		}

		if title == "" {
			return errors.New("o título da tarefa não pode estar vazio")
		}

		model := commons.GetModel(taskType)

		prompt := fmt.Sprintf("Crie uma descrição detalhada de uma tarefa do tipo %s com o título '%s'. %s "+
			"Baseando-se no modelo: %s os testes e informações para o time de infra são opcionais",
			strings.ToUpper(taskType), title, briefDesc, model)

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

		response, err := provider.GetCompletions(prompt, config.AIModel)
		if err != nil {
			return err
		}

		fmt.Println("\nDescrição gerada:")
		fmt.Println(response)

		// Copia para o clipboard
		err = clipboard.WriteAll(response)
		if err != nil {
			return errors.New("não foi possível copiar para o clipboard")
		}
		fmt.Println("\nA descrição foi copiada para o clipboard!")

		// Se o projeto estiver especificado, cria a tarefa no Jira
		if projectKey != "" {
			var issueType services.JiraIssueType
			switch strings.ToUpper(taskType) {
			case "EPICO":
				issueType = services.JiraEpic
			case "BUG":
				issueType = services.JiraBug
			default:
				issueType = services.JiraTask
			}

			issue := &services.JiraIssue{
				Summary:     title,
				Description: response,
				Type:        issueType,
				ProjectKey:  projectKey,
			}

			issueKey, err := services.CreateJiraIssue(issue)
			if err != nil {
				fmt.Printf("\nAtenção: Não foi possível criar a tarefa no Jira: %v\n", err)
			} else {
				fmt.Printf("\nTarefa criada no Jira com sucesso: %s\n", issueKey)
			}
		}

		return nil
	},
}

func init() {
	RootCmd.AddCommand(jiraCmd)

	jiraCmd.Flags().StringVarP(&title, "title", "t", "", "Título da tarefa (obrigatório)")
	jiraCmd.Flags().StringVarP(&taskType, "type", "y", "TASK", "Tipo da tarefa: EPICO, BUG ou TASK")
	jiraCmd.Flags().StringVarP(&briefDesc, "description", "d", "", "Descrição breve da tarefa (opcional)")
	jiraCmd.Flags().StringVarP(&projectKey, "project", "p", "", "Chave do projeto no Jira (opcional)")

	_ = jiraCmd.MarkFlagRequired("title")
}