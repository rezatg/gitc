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

	return fmt.Sprintf(`You're a git commit message generator. Based on this diff:
%s

Generate a professional git commit message in %s following these rules:
1. Format: "<type>: <short description>\n\n<body>"
2. Use present tense and imperative mood (e.g., "Add", "Fix", "Update")
3. Title: Max 50 characters
4. Body: Summarize all key changes in the diff, max 100 characters per line
5. Include specific changes (e.g., new features, configs, or interfaces)
6. %s
7. %s
8. Return plain text without Markdown, code blocks, or extra characters like \\\
9. Do not include prefixes like "This commit" or extra explanations

Only return the commit message, no explanations.`,
		diff,
		language,
		getTypeInstruction(commitType),
		getConventionInstruction(customMessageConvention))
}

func getTypeInstruction(commitType string) string {
	if commitType != "" {
		return fmt.Sprintf("Use type '%s'", commitType)
	}
	return "Choose appropriate type (feat, fix, docs, style, refactor, test, chore)"
}

func getConventionInstruction(convention string) string {
	if convention != "" {
		return fmt.Sprintf("Follow custom convention: %s", convention)
	}
	return "Follow Conventional Commits"
}
