package functions

import (
	"fmt"
	"gojira/services"
	"gojira/utils/commons"
	"gojira/utils/git"
)

//goland:noinspection GoPrintFunctions
func GenerateCommitMessage(diffs map[string]string, branch string) (string, error) {
	commitType, context, err := git.ParseBranchForCommitType(branch)
	if err != nil {
		return "", err
	}

	prompt := fmt.Sprintf(
		"Você é um assistente de IA. Analise os seguintes diffs do Git e, considerando o contexto da branch atual "+
			"e as regras do Git Flow, forneça uma mensagem de commit no seguinte formato, **sem introduções ou explicações adicionais**:\n\n"+
			" %s(%s) Refactor project structure for better modularity\n\n"+
			"- Item 1: Descrição breve e clara do que foi alterado ou adicionado.\n"+
			"- Item 2: Descrição breve e clara de outra alteração ou melhoria.\n"+
			"- ... (adicione mais itens conforme necessário para refletir as mudanças).\n\n"+
			"Certifique-se de que:\n"+
			"- É imprencidivel o uso de gitemoji no titulo do commit de acordo com o tipo de commit ou que reflita o contexto.\n"+
			"- Utilize prioritariamente o gitflow e o padrão de commits do projeto.\n"+
			"- O título seja objetivo e resuma a essência das alterações.\n"+
			"- Os itens mencionem os novos arquivos criados (se houver) e suas funções, sem se limitar a extensões específicas.\n"+
			"- A mensagem seja clara e reflita as mudanças de forma concisa.\n\n"+
			"Por exemplo:\n\n"+
			":gitemoji: commit_type(branch_name) Refactor project structure for better modularity\n\n"+
			"- Created new files to separate concerns and improve organization.\n"+
			"- Modularized the codebase by introducing utility and service layers.\n"+
			"- Enhanced main entry point to integrate the refactored structure.\n\n"+
			"Responda **preferencialmente** no formato acima.",
		"{{commit objective like feat, fix, chore}}", commitType, context)
	for file, diff := range diffs {
		prompt += fmt.Sprintf("Branch: %s\nFile: %s\nChanges:\n%s\n\n", branch, file, diff)
	}

	return services.CallOpenAiCompletions(prompt, commons.GetEnv("OPENAI_API_KEY"))
}
