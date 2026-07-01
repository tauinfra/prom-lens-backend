#!/usr/bin/env bash
# 推送镜像（需先 build，或设置 IMAGE 环境变量）
set -euo pipefail

if [[ -z "${IMAGE:-}" ]]; then
  echo "请设置 IMAGE，例如: IMAGE=registry.example.com/prom-lens-backend:tag $0"
  exit 1
fi

echo "==> 推送 ${IMAGE}"
docker push "${IMAGE}"
