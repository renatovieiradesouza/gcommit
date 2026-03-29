#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

cd "$SCRIPT_DIR"

echo "Baixando dependências..."
go mod download

echo "Compilando gcommit..."
go build -o gcommit main.go

echo "Aplicando permissão de execução..."
chmod +x gcommit

echo "Movendo binário para /usr/local/bin..."
sudo mv gcommit /usr/local/bin/gcommit

echo "Instalação concluída com sucesso."
