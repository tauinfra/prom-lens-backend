#!/usr/bin/env bash
# 构建 Docker 镜像
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

APP_NAME="${APP_NAME:-prom-lens-backend}"
IMAGE_REGISTRY="${IMAGE_REGISTRY:-docker.io/tau0125}"
IMAGE_TAG="${IMAGE_TAG:-}"

if [[ -z "$IMAGE_TAG" ]]; then
  BRANCH="$(git rev-parse --abbrev-ref HEAD 2>/dev/null | tr '/' '-' || echo local)"
  IMAGE_TAG="${BRANCH}-$(date +%Y%m%d%H%M%S)"
fi

if [[ -n "$IMAGE_REGISTRY" ]]; then
  IMAGE="${IMAGE_REGISTRY%/}/${APP_NAME}:${IMAGE_TAG}"
else
  IMAGE="${APP_NAME}:${IMAGE_TAG}"
fi

VERSION="${VERSION:-$(git describe --tags --always --dirty 2>/dev/null || echo dev)}"
BUILD_TIME="$(date -u +%Y-%m-%dT%H:%M:%SZ)"

echo "==> 构建镜像: ${IMAGE}"
docker build \
  -f deployments/docker/Dockerfile \
  --build-arg VERSION="${VERSION}" \
  --build-arg BUILD_TIME="${BUILD_TIME}" \
  -t "${IMAGE}" \
  .

echo "==> 完成: ${IMAGE}"
docker images "${IMAGE}" --format '大小: {{.Size}}'

if [[ "${PUSH:-0}" == "1" ]]; then
  echo "==> 推送镜像"
  docker push "${IMAGE}"
fi

