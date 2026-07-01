#!/usr/bin/env bash
# 部署到 Kubernetes（kubectl + kustomize）
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
KUSTOMIZE_DIR="${KUSTOMIZE_DIR:-${ROOT_DIR}/deployments/kubernetes}"

if ! command -v kubectl >/dev/null 2>&1; then
  echo "kubectl 未安装"
  exit 1
fi

SECRET_FILE="${KUSTOMIZE_DIR}/secret.yaml"
if [[ ! -f "$SECRET_FILE" ]]; then
  echo "未找到 ${SECRET_FILE}"
  echo "请先: cp ${KUSTOMIZE_DIR}/secret.yaml.example ${SECRET_FILE} 并填入密钥"
  exit 1
fi

TMP_DIR="$(mktemp -d)"
trap 'rm -rf "$TMP_DIR"' EXIT

cp -R "${KUSTOMIZE_DIR}/." "${TMP_DIR}/"
rm -f "${TMP_DIR}/secret.yaml.example"

if [[ -n "${IMAGE:-}" ]]; then
  IMAGE_REPO="${IMAGE%:*}"
  IMAGE_TAG="${IMAGE##*:}"
  if [[ "$(uname)" == "Darwin" ]]; then
    sed -i '' "s|newName: prom-lens-backend|newName: ${IMAGE_REPO}|" "${TMP_DIR}/kustomization.yaml"
    sed -i '' "s|newTag: latest|newTag: ${IMAGE_TAG}|" "${TMP_DIR}/kustomization.yaml"
  else
    sed -i "s|newName: prom-lens-backend|newName: ${IMAGE_REPO}|" "${TMP_DIR}/kustomization.yaml"
    sed -i "s|newTag: latest|newTag: ${IMAGE_TAG}|" "${TMP_DIR}/kustomization.yaml"
  fi
fi

echo "==> 应用 Kubernetes 资源"
kubectl apply -k "${TMP_DIR}"
kubectl apply -f "${SECRET_FILE}"

echo "==> 等待 Deployment 就绪"
kubectl rollout status deployment/prom-lens-backend -n prom-lens --timeout=120s

echo "==> 部署完成"
kubectl get pods,svc -n prom-lens -l app.kubernetes.io/name=prom-lens-backend
