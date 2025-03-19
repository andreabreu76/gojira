package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"gojira/services/ai"
	"gojira/utils/commons"
	"strings"
)

var (
	// Flags para o comando de configuração
	providerName string
	modelName    string
	jiraUrl      string
	jiraToken    string
	jiraProject  string
)

// configCmd representa o comando para configurar o aplicativo
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configura o Gojira",
	Long:  `Configura o Gojira, incluindo provedores de IA e integração com Jira.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Carrega a configuração atual
		config, err := commons.LoadConfig()
		if err != nil {
			return fmt.Errorf("erro ao carregar configuração: %w", err)
		}

		// Atualiza a configuração com os valores das flags
		if providerName != "" {
			config.AIProvider = strings.ToLower(providerName)
		}

		if modelName != "" {
			config.AIModel = modelName
		}

		if jiraUrl != "" {
			config.JiraURL = jiraUrl
		}

		if jiraToken != "" {
			config.JiraToken = jiraToken
		}

		if jiraProject != "" {
			config.DefaultJira = jiraProject
		}

		// Salva a configuração
		if err := commons.SaveConfig(config); err != nil {
			return fmt.Errorf("erro ao salvar configuração: %w", err)
		}

		fmt.Println("Configuração atualizada com sucesso!")
		return nil
	},
}

// configShowCmd representa o comando para mostrar a configuração atual
var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Mostra a configuração atual",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Carrega a configuração atual
		config, err := commons.LoadConfig()
		if err != nil {
			return fmt.Errorf("erro ao carregar configuração: %w", err)
		}

		fmt.Println("Configuração atual:")
		fmt.Printf("- Provedor de IA: %s\n", config.AIProvider)
		
		// Se o provedor existir, mostra o modelo atual e os disponíveis
		if provider, exists := ai.GetProvider(config.AIProvider); exists {
			if config.AIModel == "" {
				fmt.Printf("- Modelo de IA: %s (padrão)\n", provider.GetDefaultModel())
			} else {
				fmt.Printf("- Modelo de IA: %s\n", config.AIModel)
			}
			
			fmt.Println("- Modelos disponíveis:")
			for _, model := range provider.GetAvailableModels() {
				fmt.Printf("  * %s\n", model)
			}
		}
		
		// Mostra configuração do Jira
		if config.JiraURL != "" {
			fmt.Printf("- URL do Jira: %s\n", config.JiraURL)
		} else {
			fmt.Println("- URL do Jira: Não configurado")
		}
		
		if config.DefaultJira != "" {
			fmt.Printf("- Projeto Jira padrão: %s\n", config.DefaultJira)
		} else {
			fmt.Println("- Projeto Jira padrão: Não configurado")
		}
		
		if config.JiraToken != "" {
			fmt.Println("- Token do Jira: Configurado")
		} else {
			fmt.Println("- Token do Jira: Não configurado")
		}

		return nil
	},
}

// configProvidersCmd representa o comando para listar os provedores de IA disponíveis
var configProvidersCmd = &cobra.Command{
	Use:   "providers",
	Short: "Lista os provedores de IA disponíveis",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Provedores de IA disponíveis:")
		for name := range ai.ProviderFactory {
			provider, _ := ai.GetProvider(name)
			fmt.Printf("- %s: %s\n", name, provider.GetName())
			fmt.Println("  Modelos disponíveis:")
			for _, model := range provider.GetAvailableModels() {
				if model == provider.GetDefaultModel() {
					fmt.Printf("  * %s (padrão)\n", model)
				} else {
					fmt.Printf("  * %s\n", model)
				}
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(configCmd)
	
	// Adiciona os subcomandos
	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configProvidersCmd)
	
	// Adiciona as flags
	configCmd.Flags().StringVarP(&providerName, "provider", "p", "", "Nome do provedor de IA (openai, anthropic)")
	configCmd.Flags().StringVarP(&modelName, "model", "m", "", "Nome do modelo de IA")
	configCmd.Flags().StringVarP(&jiraUrl, "jira-url", "j", "", "URL da instância do Jira")
	configCmd.Flags().StringVarP(&jiraToken, "jira-token", "t", "", "Token de autenticação do Jira")
	configCmd.Flags().StringVarP(&jiraProject, "jira-project", "r", "", "ID do projeto Jira padrão")
}