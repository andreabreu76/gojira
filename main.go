package main

import (
	"errors"
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/spf13/cobra"
	"gojira/functions"
	"gojira/services"
	"gojira/utils/commons"
	"gojira/utils/git"
	"os"
	"strings"
)

var version = "dev"

func main() {
	commons.LoadEnv()

	var rootCmd = &cobra.Command{
		Use: "gojira",
		Short: "Uma ferramenta CLI de uso interno para criar descrições de tarefas para o Jira ou mensagens de commit " +
			"se estiver em um repositório Git e houver arquivos modificados",
		Version: version,
		RunE: func(cmd *cobra.Command, args []string) error {
			title, _ := cmd.Flags().GetString("title")
			taskType, _ := cmd.Flags().GetString("type")
			briefDesc, _ := cmd.Flags().GetString("description")

			if taskType == "Commit" {
				isRepo, err := git.IsGitRepository()
				if err != nil {
					return err
				}
				if !isRepo {
					return errors.New("o diretório atual não é um repositório Git")
				}

				fmt.Println("Repositório Git detectado. Preparando diffs...")

				branch, err := git.GetBranchName()
				if err != nil {
					return err
				}

				diffs, err := git.GetGitDiff()
				if err != nil {
					return err
				}

				if diffs != nil {
					commitMessage, err := functions.GenerateCommitMessage(diffs, branch)
					if err != nil {
						return err
					}
					fmt.Println("\nMensagem de commit sugerida:")
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

			model := commons.GetModel(taskType)

			prompt := fmt.Sprintf("Crie uma descrição detalhada de uma tarefa do tipo %s com o título '%s'. %s "+
				"Baseando-se no modelo: %s os testes e informações para o time de infra são opcionais",
				strings.ToUpper(taskType), title, briefDesc, model)

			response, err := services.CallOpenAiCompletions(prompt, commons.GetEnv("OPENAI_API_KEY"))
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
