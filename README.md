# 🧠 gcommit 2

`gcommit` é um utilitário de linha de comando em Go que substitui o uso manual de `git commit`. Ele analisa as mudanças staged no seu repositório e usa a OpenAI para gerar automaticamente uma mensagem de commit curta e significativa.

## 🚀 Funcionalidades

- Adiciona arquivos automaticamente para staged (`git add`)
- Coleta o diff com `git diff --cached`
- Usa a API da OpenAI para gerar uma mensagem de commit
- Realiza o commit com a mensagem gerada
- Otimizado para baixo uso de tokens
- Quando usado com `-a`, detecta a branch atual automaticamente
- Antes do push, executa `git pull --rebase origin <branch-atual>`
- Após sincronizar a branch local com a remota, executa `git push origin <branch-atual>`

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
```

## ▶️ Uso

Para adicionar os arquivos, gerar a mensagem automaticamente e criar o commit:

```bash
gcommit
```

Para fazer commit e sincronizar a branch atual com a remota antes do push:

```bash
gcommit -a
```

Fluxo do `gcommit -a`:

1. Executa `git add .`
2. Gera a mensagem de commit com base no diff staged
3. Executa `git commit`
4. Executa `git pull --rebase origin <branch-atual>`
5. Executa `git push origin <branch-atual>`
