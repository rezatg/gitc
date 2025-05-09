package utils

import (
	"fmt"
	"regexp"
	"strings"
)

// Define the mapping of commit types to Gitmojis
var typeToGitmoji = map[string]string{
	"feat":     "✨",  // New feature
	"fix":      "🚑",  // Bug fix
	"docs":     "📝",  // Documentation
	"style":    "💄",  // Code style
	"refactor": "♻️", // Code refactoring
	"perf":     "⚡️", // Performance improvements
	"test":     "✅",  // Tests
	"chore":    "🔧",  // Maintenance
	"build":    "🏗️", // Build system
	"ci":       "🤖",  // CI/CD
	"revert":   "⏪",  // Reverts
	"init":     "🎉",  // Initial commit
	"security": "🔒",  // Security fixes
}

// AddGitmojiToCommitMessage adds a Gitmoji to the commit message based on its type.
func AddGitmojiToCommitMessage(commitMessage string) string {
	// Extract the commit type (e.g., "feat" from "feat: description")
	match := regexp.MustCompile(`^[a-zA-Z]+`).FindString(commitMessage)
	if match == "" {
		return commitMessage // No valid type found, return unchanged
	}

	// Get the corresponding Gitmoji
	gitmoji, exists := typeToGitmoji[strings.ToLower(match)]
	if !exists {
		return commitMessage // Type not recognized, return unchanged
	}

	// Add Gitmoji to the start of the message
	return fmt.Sprintf("%s %s", gitmoji, commitMessage)
}
