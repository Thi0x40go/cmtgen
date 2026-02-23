# CommitGen

Gerador de mensagens de commit inteligente usando IA (Gemini) com integração ao Neovim.

## Configuração

A aplicação pode ser configurada via arquivo JSON em `~/.commitgen.json`:

```json
{
  "provider": "gemini",
  "gemini": {
    "api_key": "SUA_CHAVE_AQUI",
    "model": "gemini-2.5-flash"
  }
}
```

### Neovim RPC

Para que a integração com o Neovim funcione perfeitamente, especialmente em terminais externos, adicione o seguinte ao seu `autocmd.lua`:

```lua
-- Garante que o endereço do servidor RPC esteja sempre disponível
vim.api.nvim_create_autocmd("VimEnter", {
  callback = function()
    if vim.v.servername ~= "" then
      vim.env.NVIM_SERVER = vim.v.servername
    end
  end
})
```

## Como Usar

### Instalação
```bash
go build -o commitgen ./cmd/commitgen
sudo mv commitgen /usr/local/bin/
```

### Execução
- **Modo Terminal**: `commitgen`
- **Modo Neovim**: `commitgen --nvim` ou apenas rodar dentro do terminal do Neovim.
- **Mensagem Customizada**: `commitgen --nvim "feat: minha mensagem"`
