# 🧠 gcommit

`gcommit` é um utilitário de linha de comando em Go que substitui o uso manual de `git commit`. Ele analisa as mudanças staged no seu repositório e usa a OpenAI para gerar automaticamente uma mensagem de commit curta e significativa.

## 🚀 Funcionalidades

- Adiciona arquivos automaticamente para staged (`git add`)
- Coleta o diff com `git diff --cached`
- Usa a API da OpenAI para gerar uma mensagem de commit
- Realiza o commit com a mensagem gerada
- Otimizado para baixo uso de tokens
- Realiza o push baseado na branch atual

## 🛠 Requisitos

- Go 1.21 ou superior
- Git instalado
- Chave de API da OpenAI com acesso ao modelo `gpt-3.5-turbo` ou `gpt-4o`

## ⚙️ Instalação

```bash
git clone https://github.com/seuusuario/gcommit.git
cd gcommit
go build -o gcommit main.go
mv gcommit ~/bin/
