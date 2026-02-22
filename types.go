package main

const (
	OneMB       = 1_000_000
	GeminiModel = "gemini-2.5-flash"
)

type UIProvider interface {
	GetSubject() string
	ConfirmAndEdit(initialMsg string) (string, bool)
}

type CommitGen struct {
	UI     UIProvider
	Prompt string
}
