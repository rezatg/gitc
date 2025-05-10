package git

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
)

// ------------------- processDiff -------------------

func TestProcessDiff_RemovesUnnecessaryLines(t *testing.T) {
	input := `
diff --git a/main.go b/main.go
index 123456..abcdef 100644
--- a/main.go
+++ b/main.go
@@ -1,4 +1,4 @@
 package main
-import "fmt"
+import "log"
 
 func main() {
 	fmt.Println("Hello")
 }
`
	expectedContains := "@@"
	expectedDoesNotContain := []string{
		"index ", "--- ", "+++ ", " package main",
	}

	result := processDiff(input)

	if !strings.Contains(result, expectedContains) {
		t.Errorf("expected diff to contain '%s'", expectedContains)
	}

	for _, notExpected := range expectedDoesNotContain {
		if strings.Contains(result, notExpected) {
			t.Errorf("expected diff NOT to contain '%s'", notExpected)
		}
	}
}

// ------------------- getGitRoot -------------------

func TestGetGitRoot_NotInRepo(t *testing.T) {
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	tmpDir := t.TempDir()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to change dir: %v", err)
	}

	_, err := getGitRoot()
	if err == nil {
		t.Fatal("expected error when not in a git repo")
	}
	if !strings.Contains(err.Error(), "not a git repository") {
		t.Errorf("unexpected error: %v", err)
	}
}

// ------------------- GetDiffStaged -------------------

func TestGetDiffStaged_NoChanges(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Check if we're in a git repo
	_, err := exec.Command("git", "rev-parse", "--is-inside-work-tree").Output()
	if err != nil {
		t.Skip("skipping test: not inside a git repository")
	}

	// Make sure there are no staged changes
	_ = exec.Command("git", "reset").Run()

	diff, err := GetDiffStaged(ctx, nil)
	if err == nil {
		t.Error("expected error for no staged changes, got nil")
	}
	if !errors.Is(err, errors.New("no staged changes found")) && diff != "" {
		t.Errorf("unexpected diff output: %s", diff)
	}
}
