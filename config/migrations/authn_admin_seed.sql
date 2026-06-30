-- 默认管理员：admin / admin（密码为 bcrypt 哈希）
USE prom_lens;

INSERT INTO prom_lens_authn_user (
  username, password, nickname, email, phone, is_active, is_superuser, is_ldap, dn, creator, created_at, updated_at
)
VALUES (
  'admin',
  '$2a$10$yTpLvqG7Fdf4ZA3GMG7m..KxCeKvdmkdB/vJkzf1BZGnGgWXPo6.a',
  'admin',
  'admin@local',
  NULL,
  1,
  1,
  0,
  NULL,
  'system',
  NOW(3),
  NOW(3)
)
ON DUPLICATE KEY UPDATE
  password = VALUES(password),
  nickname = VALUES(nickname),
  is_active = VALUES(is_active),
  is_superuser = VALUES(is_superuser),
  updated_at = NOW(3);
