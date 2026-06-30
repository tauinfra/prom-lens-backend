-- 告警通知模块（通道、路由、标签匹配）
USE prom_lens;

CREATE TABLE IF NOT EXISTS prom_lens_alert_webhook (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  name VARCHAR(64) NOT NULL UNIQUE COMMENT '通道名，对应 Alertmanager receiver',
  url VARCHAR(512) NOT NULL COMMENT 'Lark bot webhook URL',
  description VARCHAR(255) DEFAULT NULL COMMENT '描述',
  enabled TINYINT(1) DEFAULT 1 COMMENT '是否启用',
  callback_token VARCHAR(64) NOT NULL COMMENT 'Alertmanager 回调鉴权 token',
  created_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  updated_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS prom_lens_alert_route (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  webhook_id BIGINT NOT NULL COMMENT '关联通知通道',
  enabled TINYINT(1) NOT NULL DEFAULT 1 COMMENT '是否启用该 route',
  priority INT NOT NULL DEFAULT 100 COMMENT '越小越靠前匹配',
  route_continue TINYINT(1) NOT NULL DEFAULT 0 COMMENT '匹配后是否继续后续 route',
  created_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3),
  updated_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  UNIQUE KEY uniq_alert_route_webhook (webhook_id),
  INDEX idx_alert_route_priority (priority),
  CONSTRAINT fk_alert_route_webhook
    FOREIGN KEY (webhook_id) REFERENCES prom_lens_alert_webhook(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS prom_lens_alert_route_matcher (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  route_id BIGINT NOT NULL COMMENT '所属 route',
  label VARCHAR(64) NOT NULL COMMENT '标签名',
  operator VARCHAR(8) NOT NULL DEFAULT '=' COMMENT '运算符: = != =~ !~',
  value VARCHAR(255) NOT NULL COMMENT '标签值',
  sort_order INT NOT NULL DEFAULT 0 COMMENT '同 route 内排序',
  created_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3),
  INDEX idx_alert_route_matcher_route (route_id),
  CONSTRAINT fk_alert_route_matcher_route
    FOREIGN KEY (route_id) REFERENCES prom_lens_alert_route(id) ON DELETE CASCADE,
  CONSTRAINT chk_alert_route_matcher_operator
    CHECK (operator IN ('=', '!=', '=~', '!~'))
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
