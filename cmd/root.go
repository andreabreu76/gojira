package cmd

import (
	"fmt"
	"gojira/utils/commons"
	"os"

	"github.com/spf13/cobra"
)

var (
	// Version é a versão do aplicativo
	Version = "dev"
	
	// RootCmd representa o comando base
	RootCmd = &cobra.Command{
		Use:     "gojira",
		Short:   "Uma ferramenta CLI para integração com Jira e geração de documentação usando IA",
		Version: Version,
	}
)

// Execute executa o comando root
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	commons.LoadEnv()
}