-- Prometheus 规则与采集目标模块
USE prom_lens;

CREATE TABLE IF NOT EXISTS prom_lens_prometheus_group (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  name VARCHAR(64) NOT NULL UNIQUE COMMENT '规则组名称',
  type VARCHAR(16) NOT NULL COMMENT '规则组类型',
  description VARCHAR(64) DEFAULT NULL COMMENT '规则组描述',
  created_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  updated_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  CONSTRAINT chk_prometheus_group_type CHECK (type IN ('ALERTING RULES', 'ALERTING RECORDS'))
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS prom_lens_prometheus_rule (
  id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  name VARCHAR(64) NOT NULL COMMENT '告警名称',
  group_id BIGINT NOT NULL COMMENT '规则组',
  summary VARCHAR(255) NOT NULL COMMENT '告警描述',
  description LONGTEXT NOT NULL COMMENT '告警详情',
  expr LONGTEXT NOT NULL COMMENT '告警规则',
  `for` VARCHAR(4) NOT NULL COMMENT '持续时间',
  labels JSON NOT NULL COMMENT '规则标签',
  extra_annotations JSON NOT NULL DEFAULT (JSON_OBJECT()) COMMENT '扩展 annotations（不含 summary/description）',
  status TINYINT(1) DEFAULT 1 COMMENT '规则状态',
  created_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  updated_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  UNIQUE KEY uniq_prom_rule_group_name (group_id, name),
  INDEX idx_prometheus_rule_group_id (group_id),
  CONSTRAINT fk_prometheus_rule_group FOREIGN KEY (group_id) REFERENCES prom_lens_prometheus_group(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS prom_lens_prometheus_record (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  name VARCHAR(64) NOT NULL COMMENT '记录名称',
  group_id BIGINT NOT NULL COMMENT '规则组',
  expr LONGTEXT NOT NULL COMMENT '记录规则',
  created_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  updated_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  UNIQUE KEY uniq_prom_record_group_name (group_id, name),
  INDEX idx_prometheus_record_group_id (group_id),
  CONSTRAINT fk_prometheus_record_group FOREIGN KEY (group_id) REFERENCES prom_lens_prometheus_group(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS prom_lens_prometheus_target_group (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  name VARCHAR(64) NOT NULL UNIQUE COMMENT 'Target 分组',
  description VARCHAR(64) DEFAULT NULL COMMENT 'Target 分组描述',
  labels JSON NOT NULL COMMENT '分组标签',
  created_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  updated_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS prom_lens_prometheus_target (
  id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  group_id BIGINT NOT NULL COMMENT 'Target 分组',
  ip_address VARCHAR(255) NOT NULL COMMENT 'IP 地址',
  port BIGINT NOT NULL COMMENT '端口',
  labels JSON NOT NULL COMMENT '实例标签',
  enabled TINYINT(1) DEFAULT 1 COMMENT '监控状态',
  created_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  updated_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  UNIQUE KEY uniq_prom_target_addr (group_id, ip_address, port),
  INDEX idx_prometheus_target_group_id (group_id),
  CONSTRAINT fk_prometheus_target_group FOREIGN KEY (group_id) REFERENCES prom_lens_prometheus_target_group(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
