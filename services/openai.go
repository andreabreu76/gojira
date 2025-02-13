package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

func CallOpenAiCompletions(prompt string, apiKey string) (string, error) {
	if apiKey == "" {
		return "", errors.New("OPENAI_API_KEY não fornecido")
	}

	url := "https://api.openai.com/v1/chat/completions"
	body := map[string]interface{}{
		"model":    "gpt-4o",
		"messages": []map[string]string{{"role": "user", "content": prompt}},
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
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

	return "", errors.New("resposta inesperada da API")
}
