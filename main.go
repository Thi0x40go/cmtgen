package main

import (
	"commitgen/pkg/app"
	"commitgen/pkg/config"
	"flag"
	"fmt"
	"os"
)

func main() {
	// Flags de orquestração
	useNvim := flag.Bool("nvim", false, "Habilitar interface via Neovim RPC")
	lang := flag.String("lang", "", "Definir idioma (pt, en)")
	flag.Parse()

	// Carrega a configuração (SOLID: Injeção de dependência via main)
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("⚠️  Erro ao carregar configuração: %v. Usando padrões.\n", err)
	}

	// Sobrescreve idioma se passado via flag
	if *lang != "" {
		cfg.Language = *lang
	}

	// Inicializa a aplicação (Bootstrap)
	application, err := app.NewApp(cfg, *useNvim)
	if err != nil {
		fmt.Printf("❌ Erro fatal na inicialização: %v\n", err)
		os.Exit(1)
	}

	// Executa o fluxo principal
	application.Run(flag.Args())
}
