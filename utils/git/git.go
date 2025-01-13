package git

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

func IsGitRepository() (bool, error) {
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	cmd.Stderr = nil
	err := cmd.Run()
	if err != nil {
		return false, errors.New("o diretório atual não foi identificado como um repositório Git")
	}
	return true, nil
}

func GetBranchName() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return "", errors.New("erro ao obter o nome da branch")
	}
	return strings.TrimSpace(string(output)), nil
}

func GetGitDiff() (map[string]string, error) {
	ignoredFiles := GetIgnoredFiles()

	cmd := exec.Command("git", "status", "--porcelain")
	output, err := cmd.Output()
	if err != nil {
		return nil, errors.New("erro ao executar git status")
	}

	modifiedFiles := []string{}
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		// Verifica se a linha começa com "M ", "A " ou "?? " (ignora espaços iniciais)
		if strings.HasPrefix(line, " M") || strings.HasPrefix(line, "A ") || strings.HasPrefix(line, "?? ") {
			file := strings.TrimSpace(line[3:]) // Remove o prefixo e qualquer espaço extra
			if !IsIgnored(file, ignoredFiles) {
				modifiedFiles = append(modifiedFiles, file)
			}
		}
	}

	if len(modifiedFiles) == 0 {
		return nil, errors.New("nenhum arquivo modificado ou não rastreado encontrado")
	}

	fmt.Printf("Arquivos detectados: %v\n\n", modifiedFiles)

	diffs := make(map[string]string)
	for _, file := range modifiedFiles {
		cmd = exec.Command("git", "diff", file)
		diffOutput, err := cmd.Output()
		if err != nil {
			return nil, fmt.Errorf("erro ao obter diff para o arquivo %s", file)
		}
		diffs[file] = string(diffOutput)
	}

	if len(diffs) == 0 {
		return nil, errors.New("nenhuma diferença encontrada nos arquivos selecionados")
	}

	return diffs, nil
}

func ParseBranchForCommitType(branch string) (string, string, error) {
	parts := strings.Split(branch, "/")
	if len(parts) < 2 {
		commitType := "hotfix"
		context := branch
		return commitType, context, nil
	}

	commitType := parts[0]
	context := strings.Join(parts[1:], "/")

	return commitType, context, nil
}
