#!/usr/bin/env bash
# 本地开发启动
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

CONFIG="${CONFIG:-config/config.yaml}"
EXAMPLE="${ROOT_DIR}/config/config.example.yaml"

if [[ ! -f "$CONFIG" ]]; then
  echo "未找到 ${CONFIG}，从模板复制..."
  cp "$EXAMPLE" "$CONFIG"
  echo "请编辑 ${CONFIG} 后重新运行"
  exit 1
fi

mkdir -p logs

echo "==> 启动 prom-lens-backend"
go run ./cmd/server -c "$CONFIG"
