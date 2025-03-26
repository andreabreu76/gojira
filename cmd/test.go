package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"gojira/services/ai"
	"gojira/utils/commons"
	"os"
	"path/filepath"
	"strings"
)

var (
	// Flags para o comando test
	sourceFile    string
	testFile      string
	testFramework string
	coverage      string
)

// testCmd representa o comando para gerar testes
var testCmd = &cobra.Command{
	Use:   "test-gen",
	Short: "Gera testes automaticamente para um arquivo de código",
	Long:  `Analisa o código-fonte e gera testes automaticamente utilizando inteligência artificial, cobrindo funções, classes e métodos.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if sourceFile == "" {
			return fmt.Errorf("é necessário fornecer o caminho para um arquivo fonte")
		}

		// Verifica se o arquivo fonte existe
		if _, err := os.Stat(sourceFile); os.IsNotExist(err) {
			return fmt.Errorf("o arquivo fonte %s não existe", sourceFile)
		}

		// Determina o arquivo de teste se não foi especificado
		if testFile == "" {
			dir, fileName := filepath.Split(sourceFile)
			ext := filepath.Ext(fileName)
			baseName := fileName[:len(fileName)-len(ext)]

			// Cria nome do arquivo de teste baseado na convenção da linguagem
			switch ext {
			case ".go":
				testFile = filepath.Join(dir, baseName+"_test.go")
			case ".py":
				testFile = filepath.Join(dir, "test_"+baseName+".py")
			case ".js", ".ts":
				testFile = filepath.Join(dir, baseName+".test"+ext)
			case ".java":
				testFile = filepath.Join(dir, baseName+"Test.java")
			default:
				testFile = filepath.Join(dir, baseName+"_test"+ext)
			}
		}

		// Lê o conteúdo do arquivo fonte
		sourceContent, err := os.ReadFile(sourceFile)
		if err != nil {
			return fmt.Errorf("erro ao ler o arquivo fonte: %w", err)
		}

		// Determina a linguagem com base na extensão do arquivo
		ext := filepath.Ext(sourceFile)
		language := getLanguageFromExt(ext)

		// Infere o framework de testes se não foi especificado
		if testFramework == "" {
			testFramework = inferTestFramework(language)
		}

		// Constrói o prompt para a IA
		prompt := buildTestGenerationPrompt(string(sourceContent), language, testFramework, coverage)

		// Verifica se o arquivo de teste já existe
		var existingTests string
		if _, err := os.Stat(testFile); err == nil {
			existingTestsBytes, err := os.ReadFile(testFile)
			if err == nil {
				existingTests = string(existingTestsBytes)
				prompt += fmt.Sprintf("\n\nO arquivo de teste já existe com o seguinte conteúdo. "+
					"Integre seus novos testes com os existentes, mantendo a cobertura atual e "+
					"adicionando os novos casos de teste:\n\n```\n%s\n```", existingTests)
			}
		}

		// Carrega configuração
		config, err := commons.LoadConfig()
		if err != nil {
			return fmt.Errorf("erro ao carregar configuração: %w", err)
		}

		// Obtém o provedor de IA configurado
		provider, exists := ai.GetProvider(config.AIProvider)
		if !exists {
			provider = ai.GetDefaultProvider()
		}

		// Gera os testes
		fmt.Println("Gerando testes...")
		testCode, err := provider.GetCompletions(prompt, config.AIModel)
		if err != nil {
			return fmt.Errorf("erro ao gerar testes: %w", err)
		}

		// Extrai apenas o código de teste (remove explicações e markdown)
		testCode = extractCodeFromMarkdown(testCode, language)

		// Salva os testes no arquivo
		if err := os.WriteFile(testFile, []byte(testCode), 0644); err != nil {
			return fmt.Errorf("erro ao salvar os testes: %w", err)
		}

		fmt.Printf("Testes gerados com sucesso e salvos em %s\n", testFile)
		return nil
	},
}

// buildTestGenerationPrompt cria o prompt para a IA gerar os testes
func buildTestGenerationPrompt(sourceCode, language, framework, coverage string) string {
	return fmt.Sprintf(
		"Analise o seguinte código %s e gere testes automatizados usando o framework %s. "+
			"Gere testes com uma cobertura %s, incluindo casos de teste para comportamento normal e "+
			"casos de borda. Para cada função/método, gere pelo menos um teste positivo e um negativo "+
			"quando aplicável.\n\n"+
			"Código fonte:\n```%s\n%s\n```\n\n"+
			"Gere testes completos, bem estruturados e prontos para execução. "+
			"Inclua imports/requires necessários e configure corretamente o ambiente de teste. "+
			"Os testes devem seguir as melhores práticas para %s e %s.",
		language, framework, coverage, language, sourceCode, language, framework,
	)
}

// getLanguageFromExt determina a linguagem de programação com base na extensão do arquivo
func getLanguageFromExt(ext string) string {
	ext = strings.ToLower(ext)
	switch ext {
	case ".go":
		return "Go"
	case ".py":
		return "Python"
	case ".js":
		return "JavaScript"
	case ".ts":
		return "TypeScript"
	case ".java":
		return "Java"
	case ".rb":
		return "Ruby"
	case ".php":
		return "PHP"
	case ".cs":
		return "C#"
	case ".cpp", ".cc", ".cxx":
		return "C++"
	case ".c":
		return "C"
	case ".rs":
		return "Rust"
	case ".swift":
		return "Swift"
	case ".kt":
		return "Kotlin"
	default:
		return "desconhecida"
	}
}

// inferTestFramework infere o framework de testes baseado na linguagem
func inferTestFramework(language string) string {
	switch language {
	case "Go":
		return "testing (padrão Go)"
	case "Python":
		return "pytest"
	case "JavaScript":
		return "Jest"
	case "TypeScript":
		return "Jest"
	case "Java":
		return "JUnit"
	case "Ruby":
		return "RSpec"
	case "PHP":
		return "PHPUnit"
	case "C#":
		return "xUnit"
	case "C++":
		return "Google Test"
	case "C":
		return "Unity"
	case "Rust":
		return "Rust Test"
	case "Swift":
		return "XCTest"
	case "Kotlin":
		return "JUnit"
	default:
		return "framework padrão"
	}
}

// extractCodeFromMarkdown extrai o código de teste de uma resposta markdown
func extractCodeFromMarkdown(markdown, language string) string {
	// Se a resposta já parece ser só código (sem markdown), retorna ela mesma
	if !strings.Contains(markdown, "```") {
		return markdown
	}

	// Identifica blocos de código no markdown
	var codeBlocks []string
	lines := strings.Split(markdown, "\n")
	inCodeBlock := false
	currentBlock := []string{}
	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)

		// Início de um bloco de código
		if strings.HasPrefix(trimmedLine, "```") {
			if inCodeBlock {
				// Fim de um bloco
				inCodeBlock = false
				codeBlocks = append(codeBlocks, strings.Join(currentBlock, "\n"))
				currentBlock = []string{}
			} else {
				// Início de um bloco
				inCodeBlock = true
				// Ignora a linha de abertura
			}
		} else if inCodeBlock {
			// Linha dentro do bloco de código
			currentBlock = append(currentBlock, line)
		}
	}

	// Se não encontrou blocos de código, retorna o texto original
	if len(codeBlocks) == 0 {
		return markdown
	}

	// Retorna o bloco de código maior (provavelmente é o código de teste completo)
	largestBlock := ""
	for _, block := range codeBlocks {
		if len(block) > len(largestBlock) {
			largestBlock = block
		}
	}

	return largestBlock
}

func init() {
	RootCmd.AddCommand(testCmd)

	// Flags para o comando test
	testCmd.Flags().StringVarP(&sourceFile, "source", "s", "", "Arquivo fonte para o qual gerar testes")
	testCmd.Flags().StringVarP(&testFile, "output", "o", "", "Arquivo de saída para os testes (opcional)")
	testCmd.Flags().StringVarP(&testFramework, "framework", "f", "", "Framework de testes a ser usado (opcional)")
	testCmd.Flags().StringVarP(&coverage, "coverage", "c", "alta", "Nível de cobertura desejado (básica, média, alta)")

	// Marca o parâmetro de arquivo fonte como obrigatório
	_ = testCmd.MarkFlagRequired("source")
}
