-- 审计模块
USE prom_lens;

CREATE TABLE IF NOT EXISTS prom_lens_audit_log (
  id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  username VARCHAR(32) NOT NULL COMMENT '登录账号',
  url_path VARCHAR(255) NOT NULL COMMENT 'URL地址',
  method VARCHAR(32) NOT NULL COMMENT '请求方法',
  ip_address VARCHAR(64) NOT NULL COMMENT '登录地址',
  agent VARCHAR(255) NOT NULL COMMENT '客户端',
  status_code INT NOT NULL COMMENT '状态码',
  success TINYINT(1) NOT NULL COMMENT '是否成功',
  params JSON DEFAULT NULL COMMENT '请求参数',
  response JSON NOT NULL COMMENT '响应结果',
  created_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  INDEX idx_audit_log_username (username),
  INDEX idx_audit_log_success (success),
  INDEX idx_audit_log_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS prom_lens_audit_auth_log (
  id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  username VARCHAR(64) NOT NULL COMMENT '登录账号',
  ip_address VARCHAR(64) NOT NULL COMMENT '登录地址',
  `system` VARCHAR(64) DEFAULT NULL COMMENT '系统版本',
  agent VARCHAR(255) NOT NULL COMMENT '客户端',
  status_code INT NOT NULL COMMENT '状态码',
  success TINYINT(1) NOT NULL COMMENT '是否成功',
  error_code INT DEFAULT NULL COMMENT '错误码',
  error_msg VARCHAR(255) DEFAULT NULL COMMENT '错误信息',
  created_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  INDEX idx_auth_log_username (username),
  INDEX idx_auth_log_success (success),
  INDEX idx_auth_log_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
