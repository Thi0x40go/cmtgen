package main

import (
	"commitgen/pkg/app"
	"flag"
	"fmt"
	"os"
)

func main() {
	// Flags de orquestração
	useNvim := flag.Bool("nvim", false, "Habilitar interface via Neovim RPC")
	flag.Parse()

	// Inicializa a aplicação (Bootstrap)
	// O pacote app resolve as dependências de AI, Config e UI
	application, err := app.NewApp(*useNvim)
	if err != nil {
		fmt.Printf("❌ Erro fatal na inicialização: %v\n", err)
		os.Exit(1)
	}

	// Executa o fluxo principal
	// Repassa os argumentos extras como mensagem customizada
	application.Run(flag.Args())
}
