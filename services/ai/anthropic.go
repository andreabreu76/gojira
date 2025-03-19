package ai

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"gojira/utils/commons"
	"io"
	"net/http"
	"time"
)

// AnthropicProvider implementa a interface Provider para a Anthropic
type AnthropicProvider struct {
	apiKey string
}

// NewAnthropicProvider cria uma nova instância do provedor Anthropic
func NewAnthropicProvider() Provider {
	return &AnthropicProvider{
		apiKey: commons.GetEnv("ANTHROPIC_API_KEY"),
	}
}

// GetName retorna o nome do provedor
func (p *AnthropicProvider) GetName() string {
	return "Anthropic"
}

// GetAvailableModels retorna a lista de modelos disponíveis
func (p *AnthropicProvider) GetAvailableModels() []string {
	return []string{
		"claude-3-opus-20240229",
		"claude-3-sonnet-20240229",
		"claude-3-haiku-20240307",
		"claude-3-5-sonnet-20240620",
	}
}

// GetDefaultModel retorna o modelo padrão
func (p *AnthropicProvider) GetDefaultModel() string {
	return "claude-3-5-sonnet-20240620"
}

// GetCompletions implementa a interface Provider.GetCompletions
func (p *AnthropicProvider) GetCompletions(prompt string, modelID string) (string, error) {
	if p.apiKey == "" {
		return "", errors.New("ANTHROPIC_API_KEY não fornecido")
	}

	if modelID == "" {
		modelID = p.GetDefaultModel()
	}

	url := "https://api.anthropic.com/v1/messages"
	body := map[string]interface{}{
		"model":      modelID,
		"messages":   []map[string]string{{"role": "user", "content": prompt}},
		"max_tokens": 4096,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", err
	}

	req.Header.Set("x-api-key", p.apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 120 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		formattedBody, _ := json.MarshalIndent(body, "", "  ")
		return "", fmt.Errorf("falha na chamada à API Anthropic (%d): %s", resp.StatusCode, string(formattedBody))
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	// A estrutura de resposta da Anthropic é diferente da OpenAI
	if content, ok := result["content"].([]interface{}); ok && len(content) > 0 {
		firstBlock := content[0].(map[string]interface{})
		if text, ok := firstBlock["text"].(string); ok {
			return text, nil
		}
	}

	return "", errors.New("resposta inesperada da API Anthropic")
}
