package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"gojira/services/ai"
	"gojira/utils/commons"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

var (
	// Flags para o comando de pull request
	prTitle       string
	prDescription string
	prBranch      string
	prRemote      string
	prBaseBranch  string
	prDraft       bool
)

// prCmd representa o comando para criar um pull request
var prCmd = &cobra.Command{
	Use:   "pr",
	Short: "Cria um Pull Request (PR) para a branch atual",
	Long:  `Cria um Pull Request (PR) utilizando a configuração do Git e integração com plataformas como GitHub ou GitLab. Gera automaticamente título e descrição detalhada com IA se não forem fornecidos.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Verifica se estamos em um repositório git
		gitCmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
		if err := gitCmd.Run(); err != nil {
			return fmt.Errorf("este comando deve ser executado dentro de um repositório Git")
		}

		// Obtém a branch atual se não foi especificada
		if prBranch == "" {
			branchCmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
			output, err := branchCmd.Output()
			if err != nil {
				return fmt.Errorf("erro ao obter o nome da branch atual: %w", err)
			}
			prBranch = strings.TrimSpace(string(output))
		}

		// Se não foi fornecido um título, gera um baseado nas alterações
		if prTitle == "" {
			fmt.Println("Gerando título para o PR...")
			var err error
			prTitle, err = generatePRTitle(prBranch)
			if err != nil {
				return fmt.Errorf("erro ao gerar título do PR: %w", err)
			}
			fmt.Printf("Título gerado: %s\n", prTitle)
		}

		// Se não foi fornecida uma descrição, gera uma baseada nas alterações
		if prDescription == "" {
			fmt.Println("Gerando descrição para o PR...")
			var err error
			prDescription, err = generatePRDescription(prBranch, prBaseBranch)
			if err != nil {
				return fmt.Errorf("erro ao gerar descrição do PR: %w", err)
			}
			fmt.Println("Descrição gerada.")
		}

		// Salva a descrição em um arquivo temporário
		descFile := "pr-description.md"
		if err := os.WriteFile(descFile, []byte(prDescription), 0644); err != nil {
			return fmt.Errorf("erro ao salvar a descrição do PR: %w", err)
		}
		defer os.Remove(descFile)

		// Constrói o comando para criar o PR
		var cmdArgs []string
		cmdExists, _ := exec.LookPath("gh")
		if cmdExists != "" {
			// Usa GitHub CLI se disponível
			cmdArgs = []string{"gh", "pr", "create", "--title", prTitle, "--body-file", descFile}
			if prRemote != "" {
				cmdArgs = append(cmdArgs, "--repo", prRemote)
			}
			if prBaseBranch != "" {
				cmdArgs = append(cmdArgs, "--base", prBaseBranch)
			}
			if prBranch != "" {
				cmdArgs = append(cmdArgs, "--head", prBranch)
			}
			if prDraft {
				cmdArgs = append(cmdArgs, "--draft")
			}
		} else {
			cmdExists, _ = exec.LookPath("glab")
			if cmdExists != "" {
				// Usa GitLab CLI se disponível
				cmdArgs = []string{"glab", "mr", "create", "--title", prTitle, "--description", "@" + descFile}
				if prBaseBranch != "" {
					cmdArgs = append(cmdArgs, "--base", prBaseBranch)
				}
				if prDraft {
					cmdArgs = append(cmdArgs, "--draft")
				}
			} else {
				// Fallback para git push se nenhuma CLI estiver disponível
				return fmt.Errorf("não foi encontrado 'gh' (GitHub CLI) ou 'glab' (GitLab CLI). Instale uma dessas ferramentas para criar PRs")
			}
		}

		// Executa o comando
		fmt.Printf("Criando PR usando %s...\n", cmdArgs[0])
		prCmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
		prCmd.Stdout = os.Stdout
		prCmd.Stderr = os.Stderr
		prCmd.Stdin = os.Stdin
		return prCmd.Run()
	},
}

// generatePRTitle gera um título para o PR baseado nas alterações
func generatePRTitle(branch string) (string, error) {
	// Obtém o tipo da branch (feature, bugfix, etc.)
	branchType := "feature"
	if strings.HasPrefix(branch, "fix/") || strings.HasPrefix(branch, "bugfix/") || strings.HasPrefix(branch, "hotfix/") {
		branchType = "fix"
	} else if strings.HasPrefix(branch, "docs/") {
		branchType = "docs"
	} else if strings.HasPrefix(branch, "chore/") {
		branchType = "chore"
	}

	// Extrai o ID do ticket se existir
	ticketID := ""
	parts := strings.Split(branch, "/")
	if len(parts) > 1 {
		branchName := parts[1]
		ticketParts := strings.Split(branchName, "-")
		if len(ticketParts) >= 2 {
			// Verifica se a primeira parte parece ser um ID de ticket (ex: ABC-123)
			firstPart := strings.ToUpper(ticketParts[0])
			if len(firstPart) >= 2 && len(firstPart) <= 5 {
				if _, err := strconv.Atoi(ticketParts[1]); err == nil {
					ticketID = firstPart + "-" + ticketParts[1]
				}
			}
		}
	}

	// Obtém as alterações desde a branch principal
	gitCmd := exec.Command("git", "log", "--oneline", "--no-merges", "origin/main.."+branch)
	output, err := gitCmd.Output()
	if err != nil {
		// Se falhar, tenta sem o origin/
		gitCmd = exec.Command("git", "log", "--oneline", "--no-merges", "main.."+branch)
		output, err = gitCmd.Output()
		if err != nil {
			return "", fmt.Errorf("erro ao obter log de commits: %w", err)
		}
	}

	// Constrói o prompt para a IA
	prompt := fmt.Sprintf(
		"Baseado nas seguintes alterações de commit, gere um título conciso e descritivo para um Pull Request. "+
			"O título deve começar com o tipo '%s' seguido de dois pontos. Se '%s' for um ID de ticket, inclua-o entre colchetes. "+
			"O título deve ter no máximo 72 caracteres.\n\nCommits:\n%s",
		branchType, ticketID, string(output),
	)

	// Carrega configuração
	config, err := commons.LoadConfig()
	if err != nil {
		return "", fmt.Errorf("erro ao carregar configuração: %w", err)
	}

	// Obtém o provedor de IA configurado
	provider, exists := ai.GetProvider(config.AIProvider)
	if !exists {
		provider = ai.GetDefaultProvider()
	}

	// Gera o título
	title, err := provider.GetCompletions(prompt, config.AIModel)
	if err != nil {
		return "", fmt.Errorf("erro ao gerar título com IA: %w", err)
	}

	// Limpa e formata o título
	title = strings.TrimSpace(title)
	if len(title) > 72 {
		title = title[:72]
	}

	return title, nil
}

// generatePRDescription gera uma descrição detalhada para o PR baseada nas alterações
func generatePRDescription(branch, baseBranch string) (string, error) {
	if baseBranch == "" {
		baseBranch = "main"
	}

	// Obtém a diferença entre as branches
	gitCmd := exec.Command("git", "diff", "--stat", "origin/"+baseBranch+".."+branch)
	diffStat, err := gitCmd.Output()
	if err != nil {
		// Se falhar, tenta sem o origin/
		gitCmd = exec.Command("git", "diff", "--stat", baseBranch+".."+branch)
		diffStat, err = gitCmd.Output()
		if err != nil {
			return "", fmt.Errorf("erro ao obter diff stat: %w", err)
		}
	}

	// Obtém a lista de commits
	gitCmd = exec.Command("git", "log", "--pretty=format:%h - %s (%an)", "--no-merges", "origin/"+baseBranch+".."+branch)
	commits, err := gitCmd.Output()
	if err != nil {
		// Se falhar, tenta sem o origin/
		gitCmd = exec.Command("git", "log", "--pretty=format:%h - %s (%an)", "--no-merges", baseBranch+".."+branch)
		commits, err = gitCmd.Output()
		if err != nil {
			return "", fmt.Errorf("erro ao obter commits: %w", err)
		}
	}

	// Constrói o prompt para a IA
	prompt := fmt.Sprintf(
		"Crie uma descrição detalhada para um Pull Request baseado nas seguintes alterações. "+
			"A descrição deve incluir:\n"+
			"1. Um resumo do que este PR implementa ou corrige\n"+
			"2. Contexto sobre por que essas alterações são necessárias\n"+
			"3. Quaisquer decisões técnicas importantes que foram tomadas\n"+
			"4. Como testar as alterações\n"+
			"5. Uma lista de verificação (checklist) do que foi implementado\n\n"+
			"Diferenças de arquivos:\n%s\n\n"+
			"Commits incluídos:\n%s\n\n"+
			"Formate a resposta em Markdown. Inclua títulos (##) para cada seção.",
		string(diffStat), string(commits),
	)

	// Carrega configuração
	config, err := commons.LoadConfig()
	if err != nil {
		return "", fmt.Errorf("erro ao carregar configuração: %w", err)
	}

	// Obtém o provedor de IA configurado
	provider, exists := ai.GetProvider(config.AIProvider)
	if !exists {
		provider = ai.GetDefaultProvider()
	}

	// Gera a descrição
	description, err := provider.GetCompletions(prompt, config.AIModel)
	if err != nil {
		return "", fmt.Errorf("erro ao gerar descrição com IA: %w", err)
	}

	return description, nil
}

func init() {
	RootCmd.AddCommand(prCmd)

	// Flags para o comando PR
	prCmd.Flags().StringVarP(&prTitle, "title", "t", "", "Título do PR")
	prCmd.Flags().StringVarP(&prDescription, "description", "d", "", "Descrição do PR")
	prCmd.Flags().StringVarP(&prBranch, "branch", "b", "", "Nome da branch de origem (padrão: branch atual)")
	prCmd.Flags().StringVarP(&prRemote, "remote", "r", "", "Repositório remoto (formato: owner/repo)")
	prCmd.Flags().StringVarP(&prBaseBranch, "base", "B", "main", "Branch base para o PR")
	prCmd.Flags().BoolVarP(&prDraft, "draft", "D", false, "Criar o PR como rascunho")
}