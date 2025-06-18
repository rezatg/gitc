package utils

import (
	"fmt"
	"strings"
)

func GetPromptForSingleCommit(diff, commitType, customMessageConvention, language string) string {
	language = strings.ToLower(strings.TrimSpace(language))
	if language == "" {
		language = "en"
	}

	return fmt.Sprintf(`Generate a Git commit message in %s for the following diff:

%s

Rules:
- Format:
  Line 1: <type>: <short summary> (max 50 characters, no emoji, no continuation)
  Line 2: blank
  Line 3+: (optional) short explanation, max 100 chars per line

- Use imperative mood (e.g., Add, Fix, Refactor), no past tense
- Be concise and clear
- %s
- %s
- Output must be plain text. No quotes, no Markdown, no emoji, no explanations.

Examples:
feat: add JWT middleware

Add access token verification to protected routes.

fix: prevent nil pointer on DB init

Add nil check before using DB config to avoid panic.`,
		language,
		diff,
		getTypeInstruction(commitType),
		getConventionInstruction(customMessageConvention))
}

func getTypeInstruction(commitType string) string {
	if commitType != "" {
		return fmt.Sprintf("Use type '%s'", commitType)
	}
	return "Choose appropriate type (feat, fix, docs, style, refactor, test, chore, build, ci, revert, init, security)"
}

func getConventionInstruction(convention string) string {
	if convention != "" {
		return fmt.Sprintf("Follow custom convention: %s", convention)
	}
	return "Follow Conventional Commits"
}
