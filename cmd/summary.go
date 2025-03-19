package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"gojira/services/ai"
	"gojira/utils/commons"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var (
	// Flags para o comando summary
	base        string
	format      string
	saveReport  bool
	reportFile  string
	maxChanges  int
	includeCode bool
)

// summaryCmd representa o comando para gerar um resumo das alterações
var summaryCmd = &cobra.Command{
	Use:   "summary",
	Short: "Gera um resumo das alterações desde um ponto específico",
	Long:  `Analisa as alterações no código desde um commit ou branch específica e gera um resumo detalhado do que foi alterado, adicionado ou removido.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Verifica se estamos em um repositório git
		gitCmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
		if err := gitCmd.Run(); err != nil {
			return fmt.Errorf("este comando deve ser executado dentro de um repositório Git")
		}

		// Se a base não foi especificada, usa HEAD~10 (10 commits atrás)
		if base == "" {
			base = "HEAD~10"
		}

		// Obtém a diferença entre a base e o HEAD atual
		diffCmd := exec.Command("git", "diff", "--name-only", base)
		changedFilesBytes, err := diffCmd.Output()
		if err != nil {
			return fmt.Errorf("erro ao obter lista de arquivos alterados: %w", err)
		}

		changedFiles := strings.Split(strings.TrimSpace(string(changedFilesBytes)), "\n")
		if len(changedFiles) == 0 || (len(changedFiles) == 1 && changedFiles[0] == "") {
			return fmt.Errorf("nenhuma alteração encontrada desde %s", base)
		}

		// Limita o número de arquivos se necessário
		if maxChanges > 0 && len(changedFiles) > maxChanges {
			changedFiles = changedFiles[:maxChanges]
		}

		// Obtém o diff detalhado para cada arquivo
		fileChanges := make(map[string]string)
		for _, file := range changedFiles {
			// Ignora arquivos binários, imagens, etc.
			if isIgnorableFile(file) {
				continue
			}

			diffCmd := exec.Command("git", "diff", base, "--", file)
			diffOutput, err := diffCmd.Output()
			if err != nil {
				fmt.Printf("Aviso: Erro ao obter diff para %s: %v\n", file, err)
				continue
			}

			fileChanges[file] = string(diffOutput)
		}

		// Constrói o prompt para a IA
		prompt := buildSummaryPrompt(fileChanges, includeCode)

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

		// Gera o resumo
		fmt.Println("Gerando resumo das alterações...")
		summary, err := provider.GetCompletions(prompt, config.AIModel)
		if err != nil {
			return fmt.Errorf("erro ao gerar resumo: %w", err)
		}

		// Formata o resumo conforme solicitado
		formattedSummary := formatSummary(summary, format)

		// Salva o resumo em um arquivo se solicitado
		if saveReport {
			if reportFile == "" {
				reportFile = "alteracoes-resumo.md"
			}
			if err := os.WriteFile(reportFile, []byte(formattedSummary), 0644); err != nil {
				return fmt.Errorf("erro ao salvar o resumo: %w", err)
			}
			fmt.Printf("Resumo salvo em %s\n", reportFile)
		}

		// Exibe o resumo
		fmt.Println(formattedSummary)
		return nil
	},
}

// buildSummaryPrompt cria o prompt para a IA gerar o resumo
func buildSummaryPrompt(fileChanges map[string]string, includeCode bool) string {
	var sb strings.Builder
	
	sb.WriteString("Gere um resumo detalhado das seguintes alterações em um repositório Git. ")
	sb.WriteString("Agrupe as alterações por funcionalidade ou componente, e descreva: ")
	sb.WriteString("1. As principais funcionalidades adicionadas ou modificadas\n")
	sb.WriteString("2. Correções de bugs realizadas\n")
	sb.WriteString("3. Refatorações e melhorias de código\n")
	sb.WriteString("4. Alterações de dependências ou configurações\n\n")
	
	sb.WriteString("Arquivos alterados:\n")
	for file, diff := range fileChanges {
		sb.WriteString("\n## " + file + "\n")
		if includeCode {
			sb.WriteString("```diff\n" + diff + "\n```\n")
		} else {
			// Extrai apenas os cabeçalhos do diff para reduzir o tamanho
			diffLines := strings.Split(diff, "\n")
			for _, line := range diffLines {
				if strings.HasPrefix(line, "@@") && strings.Contains(line, "@@") {
					sb.WriteString(line + "\n")
				}
			}
		}
	}
	
	sb.WriteString("\nOrganize o resumo de forma clara e concisa, destacando as alterações mais importantes. ")
	sb.WriteString("Formate o resultado usando Markdown, com títulos e listas para melhor legibilidade.")
	
	return sb.String()
}

// formatSummary formata o resumo conforme o formato solicitado
func formatSummary(summary, format string) string {
	switch strings.ToLower(format) {
	case "jira":
		// Converte markdown para formato Jira
		summary = strings.ReplaceAll(summary, "# ", "h1. ")
		summary = strings.ReplaceAll(summary, "## ", "h2. ")
		summary = strings.ReplaceAll(summary, "### ", "h3. ")
		summary = strings.ReplaceAll(summary, "**", "*")
		summary = strings.ReplaceAll(summary, "- ", "* ")
		return summary
	case "texto", "text", "plain":
		// Remove formatação markdown
		summary = strings.ReplaceAll(summary, "# ", "")
		summary = strings.ReplaceAll(summary, "## ", "")
		summary = strings.ReplaceAll(summary, "### ", "")
		summary = strings.ReplaceAll(summary, "**", "")
		summary = strings.ReplaceAll(summary, "__", "")
		summary = strings.ReplaceAll(summary, "```", "")
		return summary
	case "html":
		// Converte markdown para HTML (implementação simplificada)
		summary = strings.ReplaceAll(summary, "# ", "<h1>") + "</h1>"
		summary = strings.ReplaceAll(summary, "## ", "<h2>") + "</h2>"
		summary = strings.ReplaceAll(summary, "### ", "<h3>") + "</h3>"
		summary = strings.ReplaceAll(summary, "**", "<strong>")
		summary = strings.ReplaceAll(summary, "**", "</strong>")
		summary = strings.ReplaceAll(summary, "- ", "<li>") + "</li>"
		summary = "<html><body>" + summary + "</body></html>"
		return summary
	default:
		// Markdown (padrão)
		return summary
	}
}

// isIgnorableFile verifica se o arquivo deve ser ignorado no relatório
func isIgnorableFile(file string) bool {
	// Extensões de arquivos que não são código
	ignorableExts := []string{
		".png", ".jpg", ".jpeg", ".gif", ".bmp", ".ico", ".svg",
		".pdf", ".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx",
		".zip", ".tar", ".gz", ".rar", ".7z",
		".mp3", ".mp4", ".avi", ".mov", ".wav",
		".o", ".so", ".dll", ".exe", ".bin",
		".log", ".cache",
	}
	
	ext := strings.ToLower(filepath.Ext(file))
	for _, ignorable := range ignorableExts {
		if ext == ignorable {
			return true
		}
	}
	
	// Diretórios ou arquivos que não são código
	ignorablePaths := []string{
		"node_modules/", ".git/", "dist/", "build/", "vendor/",
		"package-lock.json", "yarn.lock", ".DS_Store",
	}
	
	for _, ignorable := range ignorablePaths {
		if strings.Contains(file, ignorable) {
			return true
		}
	}
	
	return false
}

func init() {
	RootCmd.AddCommand(summaryCmd)
	
	// Flags para o comando summary
	summaryCmd.Flags().StringVarP(&base, "base", "b", "", "Commit ou branch base para comparação (padrão: HEAD~10)")
	summaryCmd.Flags().StringVarP(&format, "format", "f", "markdown", "Formato do relatório (markdown, jira, text, html)")
	summaryCmd.Flags().BoolVarP(&saveReport, "save", "s", false, "Salvar relatório em um arquivo")
	summaryCmd.Flags().StringVarP(&reportFile, "output", "o", "", "Arquivo para salvar o relatório (padrão: alteracoes-resumo.md)")
	summaryCmd.Flags().IntVarP(&maxChanges, "max", "m", 20, "Número máximo de arquivos a incluir (0 para todos)")
	summaryCmd.Flags().BoolVarP(&includeCode, "code", "c", false, "Incluir código detalhado no prompt (aumenta precisão, mas consome mais tokens)")
}