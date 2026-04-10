#!/usr/bin/env bash
set -euo pipefail

PORT="${1:-8080}"
cd "$(dirname "$0")/.."

LOCAL_IP=$(hostname -I 2>/dev/null | awk '{print $1}')
if [[ -z "${LOCAL_IP:-}" ]]; then
  LOCAL_IP="127.0.0.1"
fi

echo "Iniciando Agent HITL em modo web na porta ${PORT}..."
echo "Acesse no navegador:"
echo "  - Local: http://localhost:${PORT}"
echo "  - Rede:  http://${LOCAL_IP}:${PORT}"

go run ./hitl -web -addr ":${PORT}"
