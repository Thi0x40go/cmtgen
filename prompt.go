package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func loadPrompt() (string, error) {
	// 1. Tentar ler do ~/.commit_message
	home, err := os.UserHomeDir()
	if err == nil {
		homePath := filepath.Join(home, ".commit_message")
		if data, err := os.ReadFile(homePath); err == nil {
			return string(data), nil
		}
	}

	// 2. Tentar ler do arquivo local (retrocompatibilidade)
	if data, err := os.ReadFile("commit_message_prompt.txt"); err == nil {
		return string(data), nil
	}

	// 3. Fallback para a string padrão no código
	return DefaultPrompt, nil
}

func truncate(text string, max int) string {
	if len(text) > max {
		return text[:max]
	}
	return text
}

func buildPrompt(base, subject, diff string) string {
	return fmt.Sprintf("%s\nSubject: %s\nDiff:\n%s", base, subject, diff)
}
