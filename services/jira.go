package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gojira/utils/commons"
	"io"
	"net/http"
	"strings"
)

// JiraIssueType representa um tipo de tarefa no Jira
type JiraIssueType string

const (
	JiraEpic JiraIssueType = "Epic"
	JiraTask JiraIssueType = "Task"
	JiraBug  JiraIssueType = "Bug"
)

// JiraIssue representa uma tarefa no Jira
type JiraIssue struct {
	Key         string       `json:"key,omitempty"`
	Summary     string       `json:"summary"`
	Description string       `json:"description"`
	Type        JiraIssueType `json:"type"`
	ProjectKey  string       `json:"projectKey"`
}

// GetJiraIssue busca uma tarefa no Jira pelo ID
func GetJiraIssue(issueID string) (*JiraIssue, error) {
	config, err := commons.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("erro ao carregar configuração: %w", err)
	}

	if config.JiraURL == "" || config.JiraToken == "" {
		return nil, fmt.Errorf("URL do Jira ou token de autenticação não configurados")
	}

	url := fmt.Sprintf("%s/rest/api/2/issue/%s", config.JiraURL, issueID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+config.JiraToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("erro ao buscar tarefa no Jira: %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	fields, ok := result["fields"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("formato de resposta do Jira inválido")
	}

	summary, _ := fields["summary"].(string)
	description, _ := fields["description"].(string)
	
	issueTypeField, ok := fields["issuetype"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("formato de tipo de tarefa do Jira inválido")
	}
	
	issueTypeName, _ := issueTypeField["name"].(string)
	
	projectField, ok := fields["project"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("formato de projeto do Jira inválido")
	}
	
	projectKey, _ := projectField["key"].(string)

	var issueType JiraIssueType
	switch strings.ToLower(issueTypeName) {
	case "epic":
		issueType = JiraEpic
	case "bug":
		issueType = JiraBug
	default:
		issueType = JiraTask
	}

	return &JiraIssue{
		Key:         issueID,
		Summary:     summary,
		Description: description,
		Type:        issueType,
		ProjectKey:  projectKey,
	}, nil
}

// CreateJiraIssue cria uma nova tarefa no Jira
func CreateJiraIssue(issue *JiraIssue) (string, error) {
	config, err := commons.LoadConfig()
	if err != nil {
		return "", fmt.Errorf("erro ao carregar configuração: %w", err)
	}

	if config.JiraURL == "" || config.JiraToken == "" {
		return "", fmt.Errorf("URL do Jira ou token de autenticação não configurados")
	}

	// Se o projeto não for especificado, usa o padrão da configuração
	if issue.ProjectKey == "" {
		issue.ProjectKey = config.DefaultJira
	}

	if issue.ProjectKey == "" {
		return "", fmt.Errorf("projeto Jira não especificado")
	}

	url := fmt.Sprintf("%s/rest/api/2/issue", config.JiraURL)
	
	// Mapeia o tipo de tarefa para o ID correspondente
	var issueTypeId string
	switch issue.Type {
	case JiraEpic:
		issueTypeId = "10000" // ID típico para Epic
	case JiraBug:
		issueTypeId = "10006" // ID típico para Bug
	default:
		issueTypeId = "10001" // ID típico para Task
	}

	// Constrói o corpo da requisição
	bodyMap := map[string]interface{}{
		"fields": map[string]interface{}{
			"project": map[string]string{
				"key": issue.ProjectKey,
			},
			"summary": issue.Summary,
			"description": issue.Description,
			"issuetype": map[string]string{
				"id": issueTypeId,
			},
		},
	}

	jsonBody, err := json.Marshal(bodyMap)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+config.JiraToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	if resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("erro ao criar tarefa no Jira: %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	issueKey, ok := result["key"].(string)
	if !ok {
		return "", fmt.Errorf("erro ao obter chave da tarefa criada")
	}

	return issueKey, nil
}