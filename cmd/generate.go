package cmd

import (
	"github.com/spf13/cobra"
	"gojira/functions"
)

// readmeCmd representa o comando para gerar o README.md
var readmeCmd = &cobra.Command{
	Use:   "readme",
	Short: "Gera um README.md com base na estrutura e conteúdo do projeto",
	RunE: func(cmd *cobra.Command, args []string) error {
		return functions.GenerateReadme()
	},
}

// analysisCmd representa o comando para gerar uma análise do código-fonte
var analysisCmd = &cobra.Command{
	Use:   "analysis",
	Short: "Gera uma análise de código-fonte com base nos arquivos do projeto",
	RunE: func(cmd *cobra.Command, args []string) error {
		return functions.GenerateAnalysis()
	},
}

// generateCmd representa o comando pai para os comandos de geração
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Gera documentação ou análises do projeto",
}

func init() {
	RootCmd.AddCommand(generateCmd)
	
	// Adiciona os comandos filhos
	generateCmd.AddCommand(readmeCmd)
	generateCmd.AddCommand(analysisCmd)
}