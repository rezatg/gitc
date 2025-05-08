package utils

import (
	"fmt"
	"strings"
)

// GetPromptForSingleCommit generates a prompt for creating a git commit message.
func GetPromptForSingleCommit(diff, commitType, customMessageConvention, language string) string {
	// Normalize language input
	language = strings.ToLower(strings.TrimSpace(language))
	if language == "" {
		language = "en"
	}

	// Start building the prompt
	prompt := fmt.Sprintf(
		"Generate a professional git commit message in %s based on the provided diff. Follow these rules:\n",
		language,
	)

	// Add Conventional Commits rules
	prompt += `- Format the message according to Conventional Commits: "<type>: <short description>\n\n<optional body>".
- Use present tense and imperative mood (e.g., "Add", "Fix", "Update").
- Keep the first line under 50 characters.
- Keep the body (if needed) under 100 characters, summarizing key changes.
- Base the message strictly on the diff content.
- Return plain text without Markdown, code blocks, or extra characters like \.
- Do not include prefixes like "This commit" or extra explanations.
- If the diff is empty or contains no meaningful changes, return an empty string.
`

	// Add commit type if provided
	if commitType != "" {
		commitType = strings.TrimSpace(commitType)
		prompt += fmt.Sprintf("- Use commit type '%s' (e.g., '%s: description').\n", commitType, commitType)
	} else {
		prompt += "- Choose an appropriate commit type (e.g., feat, fix, docs, style, refactor, test, chore) based on the diff.\n"
	}

	// Add custom message convention if provided
	if customMessageConvention != "" {
		prompt += fmt.Sprintf(
			"- Apply the following custom rules (in JSON format): %s.\n",
			customMessageConvention,
		)
	}

	// Add the diff
	prompt += "\nDiff:\n" + diff

	return prompt
}
