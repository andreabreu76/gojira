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
		"You are an AI assistant trained to generate commit messages following the Conventional Commits standard. "+
			"Analyze the Git diffs below and, based on the current branch context and Git Flow rules, generate a commit message "+
			"with the following format, **without introductions or explanations**:\n\n"+
			"  %s: [%s] Concise commit message in US English\n\n"+
			"- Item 1: Brief and clear description of what was changed or added.\n"+
			"- Item 2: Another brief description of an improvement or fix.\n"+
			"- ... (add more items if necessary).\n\n"+
			"Mandatory rules:\n"+
			"- The commit message must follow the format: `type: [TICKET] Message`.\n"+
			"- `TICKET` must be extracted from the branch name and placed in brackets.\n"+
			"- The `TICKET` is usually composed of two or more uppercase letters followed by a hyphen and a sequence of digits (e.g., ABCD-1234).\n"+
			"- Extract the `TICKET` from the branch name using this pattern: `[A-Z]{2,}-\\d+`.\n"+
			"- If no valid ticket is found in the branch name, leave this section empty.\n"+
			"- `type` must be one of the following:\n"+
			"  - `feat` for new features\n"+
			"  - `fix` for bug fixes\n"+
			"  - `chore` for maintenance tasks\n"+
			"  - `refactor` for code improvements without changing behavior\n"+
			"  - `docs` for documentation changes\n"+
			"  - `test` for adding or improving tests\n"+
			"  - `style` for formatting changes\n"+
			"- The commit title must be short and clearly describe the changes.\n"+
			"- The items must mention modified files and their functions.\n"+
			"- The message should be concise and accurately reflect the changes.\n\n"+
			"Example:\n\n"+
			"fix: [ABCD-1234] Correct last-month date calculation\n\n"+
			"- Fixed algorithm to correctly calculate the last date of the month.\n"+
			"- Resolved inconsistency in monthly report generation.\n"+
			"- Added unit tests for edge cases.\n\n"+
			"Respond **exactly** in the format above, without additional explanations.",
		commitType, context)

	for file, diff := range diffs {
		prompt += fmt.Sprintf("\n\nBranch: %s\nFile: %s\nChanges:\n%s\n", branch, file, diff)
	}

	return services.CallOpenAiCompletions(prompt, commons.GetEnv("OPENAI_API_KEY"))
}
