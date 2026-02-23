package ui

type Provider interface {
	GetSubject() string
	ConfirmAndEdit(initialMsg string) (string, bool)
}
