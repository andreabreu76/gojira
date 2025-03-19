package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"gojira/functions"
	"gojira/utils/git"
)

// commitCmd representa o comando para gerar mensagens de commit
var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Gera uma mensagem de commit com base nas alterações no Git",
	RunE: func(cmd *cobra.Command, args []string) error {
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
	},
}

func init() {
	RootCmd.AddCommand(commitCmd)
}