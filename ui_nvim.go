package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type NvimProvider struct {
	socket string
}

func NewNvimProvider() (*NvimProvider, error) {
	socket := os.Getenv("NVIM")
	if socket == "" {
		return nil, fmt.Errorf("variável $NVIM não encontrada")
	}
	return &NvimProvider{socket: socket}, nil
}

func (p *NvimProvider) callLua(expr string) (string, error) {
	cmd := exec.Command("nvim", "--server", p.socket, "--remote-expr", expr)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("%v: %s", err, string(out))
	}
	return strings.TrimSpace(string(out)), nil
}

func (p *NvimProvider) GetSubject() string {
	fmt.Println("👉 Aguardando assunto (subject) no Neovim...")

	// Executa input diretamente e retorna o resultado
	expr := `luaeval("vim.fn.input('Assunto do Commit: ')")`
	result, err := p.callLua(expr)
	if err != nil {
		fmt.Printf("⚠️ Erro ao obter assunto: %v. Usando vazio.\n", err)
		return ""
	}
	return result
}

func (p *NvimProvider) ConfirmAndEdit(msg string) (string, bool) {
	fmt.Println("🚀 Revisando mensagem no Neovim...")

	// Injetamos a lógica completa em uma função anônima auto-executável
	// _A é o argumento (msg) passado pelo luaeval
	luaCode := `
(function(msg)
  local edited = vim.fn.input('Mensagem de Commit: ', msg)
  if edited == '' then return 'CANCEL' end
  
  local choice = vim.fn.confirm('Deseja confirmar o commit?', '&Sim\n&Não', 1)
  if choice == 1 then 
    return edited 
  end
  return 'CANCEL'
end)(_A)
`
	expr := fmt.Sprintf(`luaeval(%q, %q)`, strings.TrimSpace(luaCode), msg)
	result, err := p.callLua(expr)
	if err != nil {
		fmt.Printf("\n❌ Erro no Neovim RPC:\n%v\n", err)
		return "", false
	}

	if result == "CANCEL" || result == "nil" || result == "" {
		return "", false
	}

	return result, true
}
