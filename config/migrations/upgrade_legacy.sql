-- 旧库升级脚本（仅用于已部署的旧版本，新环境请只执行各模块 *_init.sql）
-- 按顺序执行，已满足的步骤可跳过（报错时对照当前表结构）
USE prom_lens;

-- [prometheus] 增加 extra_annotations
ALTER TABLE prom_lens_prometheus_rule
  ADD COLUMN extra_annotations JSON NOT NULL DEFAULT (JSON_OBJECT())
  COMMENT '扩展 annotations（不含 summary/description）'
  AFTER labels;

-- [prometheus] 规则名改为组内唯一
ALTER TABLE prom_lens_prometheus_rule DROP INDEX name;
ALTER TABLE prom_lens_prometheus_rule
  ADD UNIQUE KEY uniq_prom_rule_group_name (group_id, name);

ALTER TABLE prom_lens_prometheus_record DROP INDEX name;
ALTER TABLE prom_lens_prometheus_record
  ADD UNIQUE KEY uniq_prom_record_group_name (group_id, name);

-- [alerting] 删除已废弃列（仅旧库存在 silence_enabled 时执行）
ALTER TABLE prom_lens_alert_webhook DROP COLUMN silence_enabled;

-- [alerting] 增加 callback_token
ALTER TABLE prom_lens_alert_webhook
  ADD COLUMN callback_token VARCHAR(64) NULL COMMENT 'Alertmanager 回调鉴权 token' AFTER enabled;

UPDATE prom_lens_alert_webhook
SET callback_token = LOWER(REPLACE(UUID(), '-', ''))
WHERE callback_token IS NULL OR callback_token = '';

ALTER TABLE prom_lens_alert_webhook
  MODIFY callback_token VARCHAR(64) NOT NULL COMMENT 'Alertmanager 回调鉴权 token';

-- [alerting] 增加 route / matcher 表（新环境已由 alerting_init.sql 创建，旧库缺表时执行）
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
