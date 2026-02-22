package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type TerminalProvider struct {
	reader *bufio.Reader
}

func NewTerminalProvider() *TerminalProvider {
	return &TerminalProvider{reader: bufio.NewReader(os.Stdin)}
}

func (p *TerminalProvider) GetSubject() string {
	fmt.Print("Digite o subject do commit (ou deixe vazio): ")
	subject, _ := p.reader.ReadString('\n')
	return strings.TrimSpace(subject)
}

func (p *TerminalProvider) ConfirmAndEdit(msg string) (string, bool) {
	fmt.Printf("\n--- Mensagem Sugerida ---\n%s\n-------------------------\n", msg)
	fmt.Print("\nDeseja fazer o commit? (y/e/n) [e=editar]: ")
	input, _ := p.reader.ReadString('\n')
	choice := strings.TrimSpace(strings.ToLower(input))

	if choice == "y" {
		return msg, true
	}
	if choice == "e" {
		fmt.Print("Digite a nova mensagem: ")
		edited, _ := p.reader.ReadString('\n')
		return strings.TrimSpace(edited), true
	}
	return "", false
}
