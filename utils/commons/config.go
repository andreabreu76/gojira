package commons

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config representa a configuração do aplicativo
type Config struct {
	AIProvider  string `json:"ai_provider"`  // Nome do provedor de IA (openai, anthropic)
	AIModel     string `json:"ai_model"`     // ID do modelo de IA a ser usado
	DefaultJira string `json:"default_jira"` // ID do projeto Jira padrão
	JiraURL     string `json:"jira_url"`     // URL da instância do Jira
	JiraToken   string `json:"jira_token"`   // Token de autenticação do Jira
}

// GetConfigFilePath retorna o caminho para o arquivo de configuração
func GetConfigFilePath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Erro ao obter diretório home: %v\n", err)
		return ".gojira.json"
	}
	return filepath.Join(homeDir, ".gojira.json")
}

// LoadConfig carrega a configuração do arquivo
func LoadConfig() (*Config, error) {
	configPath := GetConfigFilePath()
	
	// Verifica se o arquivo existe
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Retorna configuração padrão se o arquivo não existir
		return &Config{
			AIProvider: "openai",
			AIModel:    "",
		}, nil
	}
	
	// Lê o arquivo de configuração
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler arquivo de configuração: %w", err)
	}
	
	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("erro ao processar arquivo de configuração: %w", err)
	}
	
	return &config, nil
}

// SaveConfig salva a configuração no arquivo
func SaveConfig(config *Config) error {
	configPath := GetConfigFilePath()
	
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("erro ao serializar configuração: %w", err)
	}
	
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("erro ao salvar arquivo de configuração: %w", err)
	}
	
	return nil
}