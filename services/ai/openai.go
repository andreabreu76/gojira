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

// OpenAIProvider implementa a interface Provider para a OpenAI
type OpenAIProvider struct {
	apiKey string
}

// NewOpenAIProvider cria uma nova instância do provedor OpenAI
func NewOpenAIProvider() Provider {
	return &OpenAIProvider{
		apiKey: commons.GetEnv("OPENAI_API_KEY"),
	}
}

// GetName retorna o nome do provedor
func (p *OpenAIProvider) GetName() string {
	return "OpenAI"
}

// GetAvailableModels retorna a lista de modelos disponíveis
func (p *OpenAIProvider) GetAvailableModels() []string {
	return []string{
		"gpt-4o",
		"gpt-4-turbo",
		"gpt-4",
		"gpt-3.5-turbo",
	}
}

// GetDefaultModel retorna o modelo padrão
func (p *OpenAIProvider) GetDefaultModel() string {
	return "gpt-4o"
}

// GetCompletions implementa a interface Provider.GetCompletions
func (p *OpenAIProvider) GetCompletions(prompt string, modelID string) (string, error) {
	if p.apiKey == "" {
		return "", errors.New("OPENAI_API_KEY não fornecido")
	}

	if modelID == "" {
		modelID = p.GetDefaultModel()
	}

	url := "https://api.openai.com/v1/chat/completions"
	body := map[string]interface{}{
		"model":      modelID,
		"messages":   []map[string]string{{"role": "user", "content": prompt}},
		"max_tokens": 16383,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+p.apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 120 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	formattedBody, err := json.MarshalIndent(body, "", "  ")
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("falha na chamada à API (%d): %s", resp.StatusCode, string(formattedBody))
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if choices, ok := result["choices"].([]interface{}); ok && len(choices) > 0 {
		firstChoice := choices[0].(map[string]interface{})
		return firstChoice["message"].(map[string]interface{})["content"].(string), nil
	}

	return "", errors.New("resposta inesperada da API OpenAI")
}
