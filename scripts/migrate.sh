#!/usr/bin/env bash
# 执行数据库迁移（需本机 mysql 客户端）
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
MIGRATIONS_DIR="${ROOT_DIR}/config/migrations"

MYSQL_HOST="${MYSQL_HOST:-127.0.0.1}"
MYSQL_PORT="${MYSQL_PORT:-3306}"
MYSQL_USER="${MYSQL_USER:-root}"
MYSQL_PASSWORD="${MYSQL_PASSWORD:-}"
MYSQL_DB="${MYSQL_DB:-prom_lens}"

MYSQL_ARGS=(
  -h "${MYSQL_HOST}"
  -P "${MYSQL_PORT}"
  -u "${MYSQL_USER}"
)

if [[ -n "$MYSQL_PASSWORD" ]]; then
  MYSQL_ARGS+=(-p"${MYSQL_PASSWORD}")
fi

run_sql() {
  local file="$1"
  local db="${2:-}"
  echo "==> ${file}"
  if [[ -n "$db" ]]; then
    mysql "${MYSQL_ARGS[@]}" "$db" < "$file"
  else
    mysql "${MYSQL_ARGS[@]}" < "$file"
  fi
}

MODE="${1:-init}"

case "$MODE" in
  init)
    run_sql "${MIGRATIONS_DIR}/00_database.sql"
    for f in authn_init.sql audit_init.sql prometheus_init.sql alerting_init.sql authn_admin_seed.sql; do
      run_sql "${MIGRATIONS_DIR}/${f}" "${MYSQL_DB}"
    done
    ;;
  upgrade)
    run_sql "${MIGRATIONS_DIR}/upgrade_legacy.sql" "${MYSQL_DB}"
    ;;
  *)
    echo "用法: $0 [init|upgrade]"
    echo "  init    新环境全量建表（默认）"
    echo "  upgrade 旧库增量升级"
    echo ""
    echo "环境变量: MYSQL_HOST MYSQL_PORT MYSQL_USER MYSQL_PASSWORD MYSQL_DB"
    exit 1
    ;;
esac

echo "==> 迁移完成"
