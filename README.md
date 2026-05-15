# CommitGen

Gerador de mensagens de commit inteligente usando IA (Gemini) com suporte a múltiplos idiomas e integração ao Neovim.

## Instalacao

O projeto utiliza um Makefile para facilitar o processo de build e instalacao global.

```bash
# Para compilar o binario localmente
make build

# Para instalar globalmente no seu GOPATH/bin
make install
```

Certifique-se de que o diretorio bin do seu Go (geralmente ~/go/bin) esteja no seu PATH.

## Configuracao

A aplicacao pode ser configurada via arquivo JSON localizado em `~/.config/commitgen/config.json`. Caso o arquivo nao exista, o sistema tentara ler do caminho legado `~/.commitgen.json` antes de utilizar os valores padrao.

Exemplo de configuracao:

```json
{
  "language": "pt",
  "provider": "gemini",
  "gemini": {
    "api_key": "SUA_CHAVE_AQUI",
    "model": "gemini-2.5-flash"
  }
}
```

### Prioridade da API Key

O sistema segue a seguinte ordem de prioridade para a API Key:
1. Variavel de ambiente GEMINI_API_KEY.
2. Arquivo de configuracao config.json.
3. Valor padrao (leitura direta da env var).

## Como Usar

### Execucao
- **Modo Padrao**: `commitgen` (gera a mensagem baseada no diff do git).
- **Idioma customizado**: `commitgen --lang en` (sobrescreve a configuracao do arquivo).
- **Modo Neovim**: `commitgen --nvim` (utiliza a interface RPC do Neovim para confirmacao).
- **Mensagem Direta**: `commitgen "feat: minha mensagem customizada"`

### Integracao Neovim RPC

Para garantir que a comunicacao funcione corretamente ao chamar o commitgen de dentro do terminal do Neovim, adicione o seguinte ao seu arquivo de configuracao do Neovim (init.lua ou autocmd.lua):

```lua
-- Garante que o endereco do servidor RPC esteja sempre disponivel no ambiente
vim.api.nvim_create_autocmd("VimEnter", {
  callback = function()
    if vim.v.servername ~= "" then
      vim.env.NVIM_SERVER = vim.v.servername
    end
  end
})
```
