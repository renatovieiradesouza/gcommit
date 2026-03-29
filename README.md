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
- Quando usado com `-c`, executa primeiro o fluxo do `-a`
- Depois do push principal, atualiza o `change_log.txt` e cria um commit `chore` com `[skip ci]`

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

Para fazer commit e registrar a alteração no `change_log.txt`:

```bash
gcommit -c
```

Para fazer commit, sincronizar e depois registrar no change log:

```bash
gcommit -c -a
```

As flags também podem ser combinadas:

```bash
gcommit -ac
gcommit -ca
```

Mesmo quando o usuário usa `-ca`, o fluxo interno continua executando primeiro o comportamento de `-a` e depois o de `-c`.

Fluxo do `gcommit -a`:

1. Executa `git add .`
2. Gera a mensagem de commit com base no diff staged
3. Executa `git commit`
4. Executa `git pull --rebase origin <branch-atual>`
5. Executa `git push origin <branch-atual>`

Fluxo do `gcommit -c`:

1. Executa o commit principal da mudança do usuário
2. Executa `git pull --rebase origin <branch-atual>`
3. Executa `git push origin <branch-atual>`
4. Cria ou incrementa o arquivo `change_log.txt`
5. Executa um novo commit com a mensagem `chore: update change log [skip ci]`
6. Executa `git push origin <branch-atual>` para enviar o commit do change log

O `change_log.txt` registra:

1. `Change`: a mensagem do commit
2. `Date`: a data atual no formato `dd/mm/aaaa às HHhmm`
3. `Author`: o autor do último commit

Cada entrada é separada pelo texto `==================`.
