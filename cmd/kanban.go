package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"gojira/services"
	"gojira/utils/commons"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

var (
	// Flags para o comando kanban
	projectKey    string
	userFilter    string
	statusFilter  string
	limitIssues   int
	outputFormat  string
)

// kanbanCmd representa o comando para visualizar tarefas do Jira em formato kanban
var kanbanCmd = &cobra.Command{
	Use:   "kanban",
	Short: "Mostra as tarefas do Jira em formato kanban no terminal",
	Long:  `Exibe as tarefas do Jira em um formato de quadro kanban diretamente no terminal, permitindo visualizar o progresso das tarefas sem sair da linha de comando.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Carrega configuração
		config, err := commons.LoadConfig()
		if err != nil {
			return fmt.Errorf("erro ao carregar configuração: %w", err)
		}

		// Valida a configuração do Jira
		if config.JiraURL == "" || config.JiraToken == "" {
			return fmt.Errorf("configuração do Jira incompleta. Use 'gojira config' para configurar")
		}

		// Se o projeto não for especificado, usa o padrão da configuração
		if projectKey == "" {
			projectKey = config.DefaultJira
			if projectKey == "" {
				return fmt.Errorf("projeto não especificado. Use --project ou configure um projeto padrão")
			}
		}

		// Busca as tarefas do projeto
		fmt.Printf("Buscando tarefas do projeto %s...\n", projectKey)
		issues, err := fetchJiraIssues(projectKey, userFilter, statusFilter, limitIssues)
		if err != nil {
			return fmt.Errorf("erro ao buscar tarefas: %w", err)
		}

		// Organiza as tarefas por status
		issuesByStatus := organizeIssuesByStatus(issues)
		
		// Exibe as tarefas
		if outputFormat == "plain" {
			displayPlainKanban(issuesByStatus)
		} else {
			displayKanban(issuesByStatus)
		}

		return nil
	},
}

// fetchJiraIssues busca as tarefas do Jira
func fetchJiraIssues(project, user, status string, limit int) ([]*services.JiraIssue, error) {
	// Este é um stub - em uma implementação real, faria uma chamada à API do Jira
	// com os filtros apropriados. Como simplificação, vamos retornar alguns dados mockados.
	
	// Mock de tarefas
	mockIssues := []*services.JiraIssue{
		{Key: "PROJ-123", Summary: "Implementar autenticação OAuth", Type: services.JiraTask, ProjectKey: project},
		{Key: "PROJ-124", Summary: "Corrigir bug na página de login", Type: services.JiraBug, ProjectKey: project},
		{Key: "PROJ-125", Summary: "Melhorar performance da API", Type: services.JiraTask, ProjectKey: project},
		{Key: "PROJ-126", Summary: "Adicionar documentação", Type: services.JiraTask, ProjectKey: project},
		{Key: "PROJ-127", Summary: "Refatorar módulo de pagamentos", Type: services.JiraTask, ProjectKey: project},
	}
	
	// Para uma implementação real, este método faria:
	// 1. Construir uma query JQL apropriada
	// 2. Chamar a API do Jira usando a configuração do usuário
	// 3. Parsear a resposta em objetos JiraIssue
	// 4. Aplicar os filtros de usuário e status
	// 5. Limitar o número de resultados
	
	return mockIssues, nil
}

// organizeIssuesByStatus organiza as tarefas por status
func organizeIssuesByStatus(issues []*services.JiraIssue) map[string][]*services.JiraIssue {
	// Em uma implementação real, usaríamos o status real da tarefa
	// Como simplificação, vamos atribuir status aleatórios
	statuses := []string{"To Do", "In Progress", "Review", "Done"}
	result := make(map[string][]*services.JiraIssue)
	
	// Inicializa o mapa com todas as colunas
	for _, status := range statuses {
		result[status] = []*services.JiraIssue{}
	}
	
	// Distribui as tarefas pelas colunas
	for i, issue := range issues {
		status := statuses[i%len(statuses)]
		result[status] = append(result[status], issue)
	}
	
	return result
}

// displayKanban exibe as tarefas em formato kanban com cores e formatação
func displayKanban(issuesByStatus map[string][]*services.JiraIssue) {
	// Cores ANSI
	reset := "\033[0m"
	bold := "\033[1m"
	blue := "\033[34m"
	green := "\033[32m"
	yellow := "\033[33m"
	red := "\033[31m"
	
	// Determina a largura de cada coluna
	width := 25 // Largura default
	if len(issuesByStatus) > 0 {
		termWidth := getTerminalWidth()
		width = termWidth / len(issuesByStatus)
		if width < 20 {
			width = 20
		} else if width > 40 {
			width = 40
		}
	}
	
	// Cabeçalhos
	for status := range issuesByStatus {
		statusColor := blue
		switch status {
		case "To Do":
			statusColor = yellow
		case "In Progress":
			statusColor = blue
		case "Review":
			statusColor = green
		case "Done":
			statusColor = green
		}
		fmt.Printf("%s%s%s%s%s", bold, statusColor, centerText(status, width), reset, strings.Repeat(" ", 4))
	}
	fmt.Println()
	
	// Separador
	for range issuesByStatus {
		fmt.Printf("%s%s", strings.Repeat("-", width), strings.Repeat(" ", 4))
	}
	fmt.Println()
	
	// Encontra o número máximo de tarefas em uma coluna
	maxIssues := 0
	for _, issues := range issuesByStatus {
		if len(issues) > maxIssues {
			maxIssues = len(issues)
		}
	}
	
	// Imprime as tarefas
	for i := 0; i < maxIssues; i++ {
		for status, issues := range issuesByStatus {
			if i < len(issues) {
				issue := issues[i]
				issueColor := ""
				switch issue.Type {
				case services.JiraBug:
					issueColor = red
				case services.JiraTask:
					issueColor = blue
				case services.JiraEpic:
					issueColor = green
				}
				
				// Trunca o título se for muito longo
				summary := issue.Summary
				if len(summary) > width-10 {
					summary = summary[:width-13] + "..."
				}
				
				// Exibe a tarefa
				fmt.Printf("%s%s %-"+fmt.Sprintf("%d", width-7)+"s%s%s", 
					issueColor, issue.Key, summary, reset, strings.Repeat(" ", 4))
			} else {
				fmt.Printf("%s%s", strings.Repeat(" ", width), strings.Repeat(" ", 4))
			}
		}
		fmt.Println()
	}
}

// displayPlainKanban exibe as tarefas em formato texto simples
func displayPlainKanban(issuesByStatus map[string][]*services.JiraIssue) {
	for status, issues := range issuesByStatus {
		fmt.Printf("\n=== %s ===\n\n", status)
		
		if len(issues) == 0 {
			fmt.Println("Nenhuma tarefa")
			continue
		}
		
		for _, issue := range issues {
			fmt.Printf("%s: %s (%s)\n", issue.Key, issue.Summary, string(issue.Type))
		}
	}
}

// centerText centraliza o texto em um espaço de largura definida
func centerText(text string, width int) string {
	if len(text) >= width {
		return text[:width]
	}
	
	spaces := width - len(text)
	leftPad := spaces / 2
	rightPad := spaces - leftPad
	
	return strings.Repeat(" ", leftPad) + text + strings.Repeat(" ", rightPad)
}

// getTerminalWidth tenta obter a largura do terminal
func getTerminalWidth() int {
	// Valor padrão caso não consiga determinar
	defaultWidth := 80
	
	// Esta é uma implementação simples - em um caso real, usaríamos 
	// alguma biblioteca para determinar a largura do terminal
	cmd := "tput cols"
	out, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		return defaultWidth
	}
	
	width, err := strconv.Atoi(strings.TrimSpace(string(out)))
	if err != nil {
		return defaultWidth
	}
	
	return width
}

func init() {
	RootCmd.AddCommand(kanbanCmd)
	
	// Flags para o comando kanban
	kanbanCmd.Flags().StringVarP(&projectKey, "project", "p", "", "Chave do projeto Jira")
	kanbanCmd.Flags().StringVarP(&userFilter, "user", "u", "", "Filtrar por usuário")
	kanbanCmd.Flags().StringVarP(&statusFilter, "status", "s", "", "Filtrar por status")
	kanbanCmd.Flags().IntVarP(&limitIssues, "limit", "l", 10, "Número máximo de tarefas por status")
	kanbanCmd.Flags().StringVarP(&outputFormat, "format", "f", "color", "Formato de saída (color, plain)")
}