package app

import (
	"commitgen/pkg/ai"
	"commitgen/pkg/config"
	"commitgen/pkg/git"
	"commitgen/pkg/prompt"
	"commitgen/pkg/ui"
	"fmt"
	"os"
	"strings"
)

type CommitGen struct {
	UI ui.Provider
	AI ai.Provider
}

func NewApp(forceNvim bool) (*CommitGen, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("⚠️  Aviso ao carregar config: %v. Usando padrão.\n", err)
	}

	var aiProvider ai.Provider
	switch cfg.Provider {
	case "gemini":
		aiProvider = &ai.GeminiProvider{
			APIKey: cfg.Gemini.APIKey,
			Model:  cfg.Gemini.Model,
		}
	default:
		aiProvider = &ai.GeminiProvider{
			APIKey: cfg.Gemini.APIKey,
			Model:  cfg.Gemini.Model,
		}
	}

	var uiProvider ui.Provider
	if forceNvim || os.Getenv("NVIM") != "" || os.Getenv("NVIM_SERVER") != "" {
		uiProvider, err = ui.NewNvimProvider()
		if err != nil {
			fmt.Printf("⚠️  Erro ao conectar ao Neovim: %v. Usando terminal.\n", err)
			uiProvider = ui.NewTerminalProvider()
		}
	} else {
		uiProvider = ui.NewTerminalProvider()
	}

	return &CommitGen{
		UI: uiProvider,
		AI: aiProvider,
	}, nil
}

func (cg *CommitGen) Run(args []string) {
	var commitMessage string
	customMsg := strings.Join(args, " ")

	if customMsg != "" {
		commitMessage = customMsg
	} else {
		diff, err := git.GetDiff()
		if err != nil {
			fmt.Println("❌ Erro git:", err)
			return
		}

		subject := cg.UI.GetSubject()
		basePrompt, _ := prompt.LoadPrompt()
		fullPrompt := prompt.BuildPrompt(basePrompt, subject, prompt.Truncate(diff, prompt.OneMB))

		fmt.Println("\n🧠 Gerando mensagem com IA...")
		commitMessage, err = cg.AI.Generate(fullPrompt)
		if err != nil {
			fmt.Println("❌ Erro IA:", err)
			return
		}
	}

	finalMsg, ok := cg.UI.ConfirmAndEdit(commitMessage)
	if ok && finalMsg != "" {
		if err := git.Commit(finalMsg); err != nil {
			fmt.Println("❌ Falha no commit:", err)
		} else {
			fmt.Println("✅ Commit realizado com sucesso!")
		}
	} else {
		fmt.Println("👋 Operação cancelada.")
	}
}
