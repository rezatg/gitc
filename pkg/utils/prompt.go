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
1. Structure the message with:
   - First line: "<type>: <short description>" (max 50 characters)
   - Blank line
   - Body: Summarize key changes (max 100 characters per line)
2. Use present tense and imperative mood (e.g., "Add", "Fix", "Update")
3. Include specific changes (e.g., new features, configs, interfaces)
4. %s
5. %s
6. Return plain text without Markdown, code blocks, or extra characters like \
7. Do not include prefixes like "This commit" or extra explanations

Example output:
feat: add new feature

Add feature to improve performance. Update config handling.

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
