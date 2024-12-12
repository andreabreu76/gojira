#!/bin/bash

set -e

REPO_OWNER="yooga-tecnologia"
REPO_NAME="gojira"
BINARY_NAME="gojira"
INSTALL_DIR="/usr/local/bin"

OS="$(uname | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"
if [ "$ARCH" = "x86_64" ]; then
  ARCH="amd64"
elif [ "$ARCH" = "arm64" ] || [ "$ARCH" = "aarch64" ]; then
  ARCH="arm64"
else
  echo "Arquitetura $ARCH não suportada."
  exit 1
fi

LATEST_URL="https://github.com/$REPO_OWNER/$REPO_NAME/releases/latest/download/${BINARY_NAME}-${OS}-${ARCH}"
echo "Baixando $BINARY_NAME de $LATEST_URL..."

curl -L "$LATEST_URL" -o "$BINARY_NAME"
chmod +x "$BINARY_NAME"

echo "Instalando em $INSTALL_DIR..."
sudo mv "$BINARY_NAME" "$INSTALL_DIR/"

if command -v "$BINARY_NAME" > /dev/null; then
  echo "$BINARY_NAME instalado com sucesso!"
else
  echo "Houve um problema na instalação."
  exit 1
fi