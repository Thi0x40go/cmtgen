package main

import (
	"flag"
	"fmt"
	"strings"
)

func (cg *CommitGen) Run(customMsg string) {
	var commitMessage string

	if customMsg != "" {
		commitMessage = customMsg
	} else {
		diff, err := getGitDiff()
		if err != nil {
			fmt.Println("❌ Erro git:", err)
			return
		}

		subject := cg.UI.GetSubject()
		basePrompt, _ := loadPrompt()
		prompt := buildPrompt(basePrompt, subject, truncate(diff, OneMB))

		fmt.Println("\n🧠 Gerando mensagem com Gemini...")
		commitMessage, err = generateAI(prompt)
		if err != nil {
			fmt.Println("❌ Erro AI:", err)
			return
		}
	}

	finalMsg, ok := cg.UI.ConfirmAndEdit(commitMessage)
	if ok && finalMsg != "" {
		if err := executeCommit(finalMsg); err != nil {
			fmt.Println("❌ Falha no commit:", err)
		} else {
			fmt.Println("✅ Commit realizado com sucesso!")
		}
	} else {
		fmt.Println("👋 Operação cancelada.")
	}
}

func main() {
	// Definição clara dos flags
	useNvim := flag.Bool("nvim", false, "Usar o Neovim para revisão/edição")
	flag.Parse()

	var ui UIProvider
	var err error

	// Lógica de seleção:
	// 1. Só usa Neovim se o usuário pedir explicitamente via --nvim
	if *useNvim {
		ui, err = NewNvimProvider()
		if err != nil {
			fmt.Printf("⚠️  Erro ao conectar ao Neovim: %v. Usando terminal.\n", err)
			ui = NewTerminalProvider()
		}
	} else {
		// 2. Caso contrário, funciona "como atualmente" (Terminal)
		ui = NewTerminalProvider()
	}

	app := &CommitGen{UI: ui}

	// Trata argumentos restantes como mensagem customizada
	customMsg := strings.Join(flag.Args(), " ")

	app.Run(customMsg)
}
