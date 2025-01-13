package git

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func GetIgnoredFiles() []string {
	filePath := ".gitignore"
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return []string{}
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Erro ao ler o .gitignore: %v\n", err)
		return []string{}
	}

	lines := strings.Split(string(content), "\n")
	var ignoredPatterns []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "#") { // Ignora linhas vazias e comentários
			ignoredPatterns = append(ignoredPatterns, line)
		}
	}
	return ignoredPatterns
}

func IsIgnored(file string, ignoredPatterns []string) bool {
	for _, pattern := range ignoredPatterns {
		matched, err := filepath.Match(pattern, file)
		if err != nil {
			fmt.Printf("Erro ao processar o padrão %s: %v\n", pattern, err)
			continue
		}
		if matched {
			return true
		}
	}
	return false
}
