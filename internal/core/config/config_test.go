package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInitConfig_EnvOverridesSecrets(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "config.yaml")
	cfg := []byte(`
server:
  port: 8080
app:
  jwt_secret: ""
auth:
  access_token_expires: 1h
  refresh_token_expires: 24h
  issuer: prom-lens-server
  audience: prom-lens-client
database:
  host: mysql.local
  username: prom_lens
  password: ""
  dbname: prom_lens
  port: 3306
  max_idle_conns: 5
  max_open_conns: 10
  conn_max_lifetime: 1h
logging:
  log_dir: logs
  enable_console: true
base_url: http://127.0.0.1:8080
prometheus:
  target:
    namespace: monitoring
    configmap: prometheus-targets
  rule:
    namespace: thanos
    configmap: thanos-ruler-config
alerting:
  lark_timeout: 10s
  alertmanager:
    namespace: monitoring
    configmap: alertmanager-config
    config_key: alertmanager.yml
`)
	if err := os.WriteFile(cfgPath, cfg, 0o644); err != nil {
		t.Fatal(err)
	}

	t.Setenv("PROM_LENS_DATABASE_PASSWORD", "env-db-password")
	t.Setenv("PROM_LENS_APP_JWT_SECRET", "env-jwt-secret-16c")

	got, err := InitConfig(cfgPath)
	if err != nil {
		t.Fatalf("InitConfig: %v", err)
	}
	if got.Database.Password != "env-db-password" {
		t.Fatalf("database.password = %q, want env override", got.Database.Password)
	}
	if got.App.JWTSecret != "env-jwt-secret-16c" {
		t.Fatalf("app.jwt_secret = %q, want env override", got.App.JWTSecret)
	}
}
