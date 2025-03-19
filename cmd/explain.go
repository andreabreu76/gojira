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
	// Flags para o comando explain
	filePath   string
	lineStart  int
	lineEnd    int
	outputFile string
	langLevel  string
)

// explainCmd representa o comando para explicar código
var explainCmd = &cobra.Command{
	Use:   "explain",
	Short: "Explica o funcionamento de trechos de código",
	Long:  `Analisa e explica o funcionamento de trechos de código, classes, funções ou arquivos inteiros, tornando mais fácil entender código complexo ou legado.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if filePath == "" {
			return fmt.Errorf("é necessário fornecer o caminho para um arquivo")
		}

		// Verifica se o arquivo existe
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			return fmt.Errorf("o arquivo %s não existe", filePath)
		}

		// Lê o conteúdo do arquivo
		content, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("erro ao ler o arquivo: %w", err)
		}

		// Se não foram especificadas linhas, processa o arquivo inteiro
		var codeToExplain string
		if lineStart == 0 && lineEnd == 0 {
			codeToExplain = string(content)
		} else {
			lines := strings.Split(string(content), "\n")
			
			// Validação das linhas
			if lineStart < 1 {
				lineStart = 1
			}
			if lineEnd == 0 || lineEnd > len(lines) {
				lineEnd = len(lines)
			}
			if lineStart > lineEnd {
				lineStart, lineEnd = lineEnd, lineStart
			}

			// Extrai o código das linhas especificadas
			selectedLines := lines[lineStart-1:lineEnd]
			codeToExplain = strings.Join(selectedLines, "\n")
		}

		// Determina a linguagem com base na extensão do arquivo
		language := getLanguageFromExtension(filepath.Ext(filePath))

		// Constrói o prompt para a IA
		prompt := buildExplanationPrompt(codeToExplain, language)

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

		// Gera a explicação
		explanation, err := provider.GetCompletions(prompt, config.AIModel)
		if err != nil {
			return fmt.Errorf("erro ao gerar explicação: %w", err)
		}

		// Salva a explicação em um arquivo se solicitado
		if outputFile != "" {
			if err := os.WriteFile(outputFile, []byte(explanation), 0644); err != nil {
				return fmt.Errorf("erro ao salvar a explicação: %w", err)
			}
			fmt.Printf("Explicação salva em %s\n", outputFile)
		}

		// Exibe a explicação
		fmt.Println(explanation)
		return nil
	},
}

// buildExplanationPrompt cria o prompt para a IA explicar o código
func buildExplanationPrompt(code, language string) string {
	// Define o nível de experiência
	experienceLevel := "experiente"
	switch strings.ToLower(langLevel) {
	case "beginner", "iniciante":
		experienceLevel = "iniciante"
	case "intermediate", "intermediario", "intermediário":
		experienceLevel = "intermediário"
	}

	return fmt.Sprintf(
		"Explique o seguinte código %s para um desenvolvedor de nível %s. "+
			"Forneça uma análise detalhada que inclua:\n\n"+
			"1. Visão geral do que o código faz\n"+
			"2. Explicação de cada seção ou função importante\n"+
			"3. Identificação de padrões ou técnicas utilizadas\n"+
			"4. Possíveis melhorias ou otimizações\n"+
			"5. Potenciais problemas ou bugs\n\n"+
			"Código:\n```%s\n%s\n```\n\n"+
			"Formate a resposta em Markdown, usando títulos e blocos de código quando apropriado.",
		language, experienceLevel, language, code,
	)
}

// getLanguageFromExtension determina a linguagem de programação com base na extensão do arquivo
func getLanguageFromExtension(ext string) string {
	ext = strings.TrimPrefix(ext, ".")
	languageMap := map[string]string{
		"go":     "Go",
		"py":     "Python",
		"js":     "JavaScript",
		"ts":     "TypeScript",
		"java":   "Java",
		"c":      "C",
		"cpp":    "C++",
		"cs":     "C#",
		"php":    "PHP",
		"rb":     "Ruby",
		"swift":  "Swift",
		"kt":     "Kotlin",
		"rs":     "Rust",
		"sh":     "Shell",
		"sql":    "SQL",
		"html":   "HTML",
		"css":    "CSS",
		"json":   "JSON",
		"xml":    "XML",
		"yaml":   "YAML",
		"yml":    "YAML",
		"md":     "Markdown",
		"dart":   "Dart",
		"scala":  "Scala",
		"pl":     "Perl",
		"clj":    "Clojure",
		"exs":    "Elixir",
		"ex":     "Elixir",
		"f":      "Fortran",
		"f90":    "Fortran",
		"hs":     "Haskell",
		"lua":    "Lua",
		"m":      "Objective-C",
		"r":      "R",
	}

	if language, ok := languageMap[strings.ToLower(ext)]; ok {
		return language
	}
	return "desconhecida"
}

func init() {
	RootCmd.AddCommand(explainCmd)

	// Flags para o comando explain
	explainCmd.Flags().StringVarP(&filePath, "file", "f", "", "Caminho para o arquivo a ser explicado")
	explainCmd.Flags().IntVarP(&lineStart, "start", "s", 0, "Linha inicial (opcional)")
	explainCmd.Flags().IntVarP(&lineEnd, "end", "e", 0, "Linha final (opcional)")
	explainCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Arquivo para salvar a explicação (opcional)")
	explainCmd.Flags().StringVarP(&langLevel, "level", "l", "intermediate", "Nível de experiência do desenvolvedor (beginner, intermediate, expert)")

	// Marca o parâmetro de arquivo como obrigatório
	_ = explainCmd.MarkFlagRequired("file")
}