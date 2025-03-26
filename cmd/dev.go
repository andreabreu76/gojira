package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"gojira/services"
	"gojira/services/ai"
	"gojira/utils/commons"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var (
	// Flags para o comando de desenvolvimento
	branchName    string
	branchPrefix  string
	issueKey      string
	skipChecklist bool
)

// devCmd representa o comando para tarefas de desenvolvimento
var devCmd = &cobra.Command{
	Use:   "dev",
	Short: "Comandos para auxiliar no desenvolvimento",
}

// branchCmd representa o comando para criar uma nova branch
var branchCmd = &cobra.Command{
	Use:   "branch",
	Short: "Cria uma nova branch com o padrão correto",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Se não foi fornecido um nome de branch, pergunta para o usuário
		if branchName == "" {
			fmt.Println("Digite o nome da branch (sem o prefixo):")
			_, err := fmt.Scanln(&branchName)
			if err != nil {
				return fmt.Errorf("erro ao ler o nome da branch: %w", err)
			}
		}

		// Se o prefixo não foi fornecido, usa o padrão
		if branchPrefix == "" {
			branchPrefix = "feature"
		}

		// Se a issue não foi fornecida e o nome não contém uma issue,
		// solicita a issue do usuário
		if issueKey == "" && !strings.Contains(branchName, "-") {
			fmt.Println("Digite a chave da issue (ex: ABC-123):")
			_, err := fmt.Scanln(&issueKey)
			if err != nil {
				return fmt.Errorf("erro ao ler a chave da issue: %w", err)
			}
		}

		// Formata o nome da branch
		formattedName := strings.ToLower(branchName)
		formattedName = strings.ReplaceAll(formattedName, " ", "-")

		// Adiciona a issue ao nome se fornecida
		if issueKey != "" {
			formattedName = fmt.Sprintf("%s/%s-%s", branchPrefix, strings.ToUpper(issueKey), formattedName)
		} else {
			formattedName = fmt.Sprintf("%s/%s", branchPrefix, formattedName)
		}

		// Cria a branch
		gitCmd := exec.Command("git", "checkout", "-b", formattedName)
		gitCmd.Stdout = os.Stdout
		gitCmd.Stderr = os.Stderr
		if err := gitCmd.Run(); err != nil {
			return fmt.Errorf("erro ao criar a branch: %w", err)
		}

		fmt.Printf("Branch %s criada com sucesso!\n", formattedName)
		return nil
	},
}

// startCmd representa o comando para iniciar o trabalho em uma tarefa
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Inicia o trabalho em uma tarefa",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Verifica se uma issue foi fornecida
		if issueKey == "" {
			fmt.Println("Digite a chave da issue (ex: ABC-123):")
			_, err := fmt.Scanln(&issueKey)
			if err != nil {
				return fmt.Errorf("erro ao ler a chave da issue: %w", err)
			}
		}

		// Busca os detalhes da issue no Jira
		issue, err := services.GetJiraIssue(issueKey)
		if err != nil {
			fmt.Printf("Não foi possível obter detalhes da issue %s: %v\n", issueKey, err)
			fmt.Println("Deseja continuar mesmo assim? (s/n)")
			var response string
			_, err := fmt.Scanln(&response)
			if err != nil || strings.ToLower(response) != "s" {
				return fmt.Errorf("operação cancelada pelo usuário")
			}
		} else {
			fmt.Printf("Tarefa: %s - %s\n", issue.Key, issue.Summary)

			// Cria um nome de branch a partir do título da tarefa
			suggestedBranchName := strings.ToLower(issue.Summary)
			suggestedBranchName = strings.ReplaceAll(suggestedBranchName, " ", "-")
			// Remove caracteres especiais
			suggestedBranchName = strings.Map(func(r rune) rune {
				if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
					return r
				}
				return -1
			}, suggestedBranchName)

			// Limita o tamanho do nome da branch
			if len(suggestedBranchName) > 50 {
				suggestedBranchName = suggestedBranchName[:50]
			}

			// Define o tipo de branch com base no tipo de issue
			switch issue.Type {
			case services.JiraEpic:
				branchPrefix = "epic"
			case services.JiraBug:
				branchPrefix = "fix"
			default:
				branchPrefix = "feature"
			}

			// Confirma o nome da branch
			fmt.Printf("Nome sugerido para a branch: %s/%s-%s\n", branchPrefix, issue.Key, suggestedBranchName)
			fmt.Println("Deseja usar este nome? (s/n)")
			var response string
			_, err := fmt.Scanln(&response)
			if err == nil && strings.ToLower(response) == "s" {
				branchName = suggestedBranchName
				issueKey = issue.Key
			} else {
				fmt.Println("Digite o nome desejado para a branch (sem o prefixo e sem a issue):")
				_, err := fmt.Scanln(&branchName)
				if err != nil {
					return fmt.Errorf("erro ao ler o nome da branch: %w", err)
				}
			}
		}

		// Cria a branch
		return branchCmd.RunE(cmd, args)
	},
}

// checklistCmd representa o comando para gerar um checklist de tarefas
var checklistCmd = &cobra.Command{
	Use:   "checklist",
	Short: "Gera um checklist de tarefas para a issue atual",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Obtém a branch atual
		gitCmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
		output, err := gitCmd.Output()
		if err != nil {
			return fmt.Errorf("erro ao obter o nome da branch atual: %w", err)
		}

		branch := strings.TrimSpace(string(output))

		// Extrai a issue da branch
		issuePattern := "[A-Z]+-[0-9]+"
		matches, err := filepath.Match(issuePattern, branch)
		if !matches {
			// Se não encontrou pela correspondência direta, tenta extrair usando separador
			parts := strings.Split(branch, "/")
			if len(parts) > 1 {
				issueParts := strings.Split(parts[1], "-")
				if len(issueParts) > 1 {
					issueKey = strings.Join(issueParts[:2], "-")
				}
			}
		}

		// Se ainda não temos a issue, pergunta ao usuário
		if issueKey == "" {
			fmt.Println("Digite a chave da issue (ex: ABC-123):")
			_, err := fmt.Scanln(&issueKey)
			if err != nil {
				return fmt.Errorf("erro ao ler a chave da issue: %w", err)
			}
		}

		// Busca os detalhes da issue no Jira
		issue, err := services.GetJiraIssue(issueKey)
		if err != nil {
			return fmt.Errorf("não foi possível obter detalhes da issue %s: %w", issueKey, err)
		}

		// Constrói o prompt para gerar o checklist
		prompt := fmt.Sprintf(
			"Crie um checklist detalhado para a issue '%s' com título '%s'.\n\n"+
				"Descrição da issue:\n%s\n\n"+
				"O checklist deve incluir etapas para:\n"+
				"1. Preparação do ambiente\n"+
				"2. Implementação da solução\n"+
				"3. Testes a serem realizados\n"+
				"4. Revisão de código\n"+
				"5. Documentação\n\n"+
				"Formato o checklist como uma lista de tarefas em Markdown com checkboxes, por exemplo:\n"+
				"- [ ] Tarefa 1\n"+
				"- [ ] Tarefa 2\n   - [ ] Subtarefa 2.1\n",
			issue.Key, issue.Summary, issue.Description,
		)

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

		// Gera o checklist
		checklist, err := provider.GetCompletions(prompt, config.AIModel)
		if err != nil {
			return fmt.Errorf("erro ao gerar checklist: %w", err)
		}

		// Salva o checklist em um arquivo
		filename := fmt.Sprintf("checklist-%s.md", issueKey)
		err = os.WriteFile(filename, []byte(checklist), 0644)
		if err != nil {
			return fmt.Errorf("erro ao salvar checklist: %w", err)
		}

		fmt.Printf("Checklist gerado e salvo em %s\n", filename)
		fmt.Println("\nChecklist:")
		fmt.Println(checklist)

		return nil
	},
}

func init() {
	RootCmd.AddCommand(devCmd)

	// Adiciona os subcomandos
	devCmd.AddCommand(branchCmd)
	devCmd.AddCommand(startCmd)
	devCmd.AddCommand(checklistCmd)

	// Flags para o comando branch
	branchCmd.Flags().StringVarP(&branchName, "name", "n", "", "Nome da branch")
	branchCmd.Flags().StringVarP(&branchPrefix, "prefix", "p", "", "Prefixo da branch (feature, fix, chore)")
	branchCmd.Flags().StringVarP(&issueKey, "issue", "i", "", "Chave da issue (ex: ABC-123)")

	// Flags para o comando start
	startCmd.Flags().StringVarP(&issueKey, "issue", "i", "", "Chave da issue (ex: ABC-123)")
	startCmd.Flags().BoolVarP(&skipChecklist, "no-checklist", "c", false, "Não gerar checklist")

	// Flags para o comando checklist
	checklistCmd.Flags().StringVarP(&issueKey, "issue", "i", "", "Chave da issue (ex: ABC-123)")
}
