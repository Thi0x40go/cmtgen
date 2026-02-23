package prompt

import (
	"fmt"
	"os"
	"path/filepath"
)

const OneMB = 1_000_000

const DefaultPrompt = `# Git Commit Message Rules (Conventional Commits)

You are an expert at writing Git commit messages following the  
**Conventional Commits v1.0.0** specification.

Your job is to generate a short, clear commit message that summarizes the
changes.

The commit message **must be written in Portuguese** and follow this format:

tipo(escopo opcional): descrição

---

## Scope Rules (important)

The **scope is optional**, but when used it **must follow these rules**:

- Must be a noun describing a section of the codebase  
- Must be lowercase  
- Must not contain spaces  
- Should represent a module, feature, or layer

### Valid examples

feat(auth): adicionar login com Google
fix(api): corrigir erro de paginação

### Rules

- Always use the Conventional Commits format  
- Use the imperative mood in the description  
- The subject line must be no longer than 50 characters
- Return **only** the commit message, with no extra commentary
`

func LoadPrompt() (string, error) {
	home, err := os.UserHomeDir()
	if err == nil {
		homePath := filepath.Join(home, ".commit_message")
		if data, err := os.ReadFile(homePath); err == nil {
			return string(data), nil
		}
	}

	if data, err := os.ReadFile("commit_message_prompt.txt"); err == nil {
		return string(data), nil
	}

	return DefaultPrompt, nil
}

func Truncate(text string, max int) string {
	if len(text) > max {
		return text[:max]
	}
	return text
}

func BuildPrompt(base, subject, diff string) string {
	return fmt.Sprintf("%s\nSubject: %s\nDiff:\n%s", base, subject, diff)
}
