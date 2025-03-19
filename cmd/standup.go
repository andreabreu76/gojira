package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"gojira/services/ai"
	"gojira/utils/commons"
	"os"
	"os/exec"
	"strings"
	"time"
)

var (
	// Flags para o comando standup
	days       int
	userEmail  string
	teamOnly   bool
	issuesOnly bool
	exportFile string
)

// standupCmd representa o comando para gerar um relatório de standup
var standupCmd = &cobra.Command{
	Use:   "standup",
	Short: "Gera um relatório de standup com base nas atividades recentes",
	Long:  `Analisa as atividades recentes (commits, issues, etc.) e gera um relatório formatado para reuniões de standup diárias, detalhando o que foi feito, o que está planejado e quaisquer bloqueios.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Verifica se estamos em um repositório git
		gitCmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
		if err := gitCmd.Run(); err != nil {
			return fmt.Errorf("este comando deve ser executado dentro de um repositório Git")
		}

		// Se o email do usuário não foi especificado, tenta obter o email do git config
		if userEmail == "" {
			emailCmd := exec.Command("git", "config", "user.email")
			emailBytes, err := emailCmd.Output()
			if err == nil {
				userEmail = strings.TrimSpace(string(emailBytes))
			}
		}

		// Coleta informações para o standup
		activities, err := collectActivities(days, userEmail, teamOnly, issuesOnly)
		if err != nil {
			return fmt.Errorf("erro ao coletar atividades: %w", err)
		}

		// Constrói o prompt para a IA
		prompt := buildStandupPrompt(activities, days)

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

		// Gera o relatório
		fmt.Println("Gerando relatório de standup...")
		standupReport, err := provider.GetCompletions(prompt, config.AIModel)
		if err != nil {
			return fmt.Errorf("erro ao gerar relatório: %w", err)
		}

		// Salva o relatório em um arquivo se solicitado
		if exportFile != "" {
			if err := os.WriteFile(exportFile, []byte(standupReport), 0644); err != nil {
				return fmt.Errorf("erro ao salvar o relatório: %w", err)
			}
			fmt.Printf("Relatório salvo em %s\n", exportFile)
		}

		// Exibe o relatório
		fmt.Println(standupReport)
		return nil
	},
}

// Estrutura para armazenar atividades coletadas
type Activities struct {
	Commits   []string
	Issues    []string
	PullReqs  []string
	UserName  string
	RepoName  string
	WorksInJira   bool // Indica se existem integrações com Jira
	HasIssues     bool // Indica se existem issues no projeto
}

// collectActivities coleta as atividades recentes (commits, issues, etc.)
func collectActivities(days int, userEmail string, teamOnly, issuesOnly bool) (*Activities, error) {
	activities := &Activities{
		Commits:  []string{},
		Issues:   []string{},
		PullReqs: []string{},
	}

	// Obtém o nome do repositório
	remoteCmd := exec.Command("git", "remote", "get-url", "origin")
	remoteBytes, err := remoteCmd.Output()
	if err == nil {
		remote := strings.TrimSpace(string(remoteBytes))
		parts := strings.Split(remote, "/")
		if len(parts) > 1 {
			repoWithGit := parts[len(parts)-1]
			activities.RepoName = strings.TrimSuffix(repoWithGit, ".git")
		}
	}

	// Obtém o nome do usuário
	if userEmail != "" {
		nameCmd := exec.Command("git", "log", "-1", "--author="+userEmail, "--format=%an")
		nameBytes, err := nameCmd.Output()
		if err == nil {
			activities.UserName = strings.TrimSpace(string(nameBytes))
		}
	}

	// Calcula a data de início (X dias atrás)
	since := time.Now().AddDate(0, 0, -days).Format("2006-01-02")

	// Coleta commits
	if !issuesOnly {
		var commitCmd *exec.Cmd
		if userEmail != "" && !teamOnly {
			// Commits do usuário específico
			commitCmd = exec.Command("git", "log", "--since="+since, "--author="+userEmail, "--format=%h | %s")
		} else {
			// Todos os commits
			commitCmd = exec.Command("git", "log", "--since="+since, "--format=%h | %an | %s")
		}
		
		commitBytes, err := commitCmd.Output()
		if err == nil {
			commitLines := strings.Split(strings.TrimSpace(string(commitBytes)), "\n")
			for _, line := range commitLines {
				if line != "" {
					activities.Commits = append(activities.Commits, line)
				}
			}
		}
	}

	// Tenta obter issues e PRs do GitHub se o gh cli estiver disponível
	ghPath, _ := exec.LookPath("gh")
	if ghPath != "" {
		// Verifica se tem issues
		issueListCmd := exec.Command("gh", "issue", "list", "--limit", "1")
		issueListCmd.Stderr = nil
		if err := issueListCmd.Run(); err == nil {
			activities.HasIssues = true
			
			// Coleta issues atribuídas ao usuário ou criadas recentemente
			issueCmd := exec.Command("gh", "issue", "list", "--limit", "10", "--state", "open")
			issueBytes, err := issueCmd.Output()
			if err == nil {
				issueLines := strings.Split(strings.TrimSpace(string(issueBytes)), "\n")
				for _, line := range issueLines {
					if line != "" {
						activities.Issues = append(activities.Issues, line)
					}
				}
			}
			
			// Coleta PRs abertos
			prCmd := exec.Command("gh", "pr", "list", "--limit", "5", "--state", "open")
			prBytes, err := prCmd.Output()
			if err == nil {
				prLines := strings.Split(strings.TrimSpace(string(prBytes)), "\n")
				for _, line := range prLines {
					if line != "" {
						activities.PullReqs = append(activities.PullReqs, line)
					}
				}
			}
		}
	}

	// Verifica integração com Jira (simplificado - apenas verifica a configuração)
	config, err := commons.LoadConfig()
	if err == nil && config.JiraURL != "" && config.JiraToken != "" {
		activities.WorksInJira = true
	}

	return activities, nil
}

// buildStandupPrompt cria o prompt para a IA gerar o relatório de standup
func buildStandupPrompt(activities *Activities, days int) string {
	var sb strings.Builder
	
	sb.WriteString("Gere um relatório para uma reunião de standup diária com base nas atividades a seguir. ")
	sb.WriteString("O relatório deve seguir o formato padrão de standup:\n\n")
	sb.WriteString("1. O que foi feito (últimos " + fmt.Sprintf("%d", days) + " dias)\n")
	sb.WriteString("2. O que será feito hoje\n")
	sb.WriteString("3. Existe algum bloqueador?\n\n")
	
	// Adiciona contexto sobre o usuário e projeto
	if activities.UserName != "" {
		sb.WriteString("Usuário: " + activities.UserName + "\n")
	}
	if activities.RepoName != "" {
		sb.WriteString("Projeto: " + activities.RepoName + "\n\n")
	}
	
	// Adiciona commits
	if len(activities.Commits) > 0 {
		sb.WriteString("## Commits recentes:\n\n")
		for _, commit := range activities.Commits {
			sb.WriteString("- " + commit + "\n")
		}
		sb.WriteString("\n")
	}
	
	// Adiciona issues
	if len(activities.Issues) > 0 {
		sb.WriteString("## Issues abertas:\n\n")
		for _, issue := range activities.Issues {
			sb.WriteString("- " + issue + "\n")
		}
		sb.WriteString("\n")
	}
	
	// Adiciona PRs
	if len(activities.PullReqs) > 0 {
		sb.WriteString("## Pull Requests abertos:\n\n")
		for _, pr := range activities.PullReqs {
			sb.WriteString("- " + pr + "\n")
		}
		sb.WriteString("\n")
	}
	
	// Dicas para o modelo baseadas no contexto
	sb.WriteString("Com base nessas informações, gere um relatório conciso e informativo para um standup. ")
	sb.WriteString("Infira as tarefas atuais e planejadas dos commits e issues. ")
	
	if activities.WorksInJira {
		sb.WriteString("O usuário trabalha com Jira, então inclua referências a tickets do Jira se identificados nos commits. ")
	}
	
	if !activities.HasIssues && len(activities.Commits) == 0 {
		sb.WriteString("Não há muitas informações disponíveis, então faça suposições razoáveis sobre o trabalho baseado no nome do projeto. ")
	}
	
	sb.WriteString("\nFormate o relatório de forma limpa e profissional. Use listas com marcadores para facilitar a leitura.")
	
	return sb.String()
}

func init() {
	RootCmd.AddCommand(standupCmd)
	
	// Flags para o comando standup
	standupCmd.Flags().IntVarP(&days, "days", "d", 1, "Número de dias para incluir no relatório")
	standupCmd.Flags().StringVarP(&userEmail, "email", "e", "", "Email do usuário para filtrar as atividades (padrão: email do git config)")
	standupCmd.Flags().BoolVarP(&teamOnly, "team", "t", false, "Incluir atividades de toda a equipe, não apenas do usuário")
	standupCmd.Flags().BoolVarP(&issuesOnly, "issues", "i", false, "Focar apenas em issues, ignorando commits")
	standupCmd.Flags().StringVarP(&exportFile, "output", "o", "", "Arquivo para salvar o relatório (opcional)")
}