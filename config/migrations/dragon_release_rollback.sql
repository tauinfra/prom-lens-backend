-- 回滚功能：为 release 表增加回滚目标 release id
ALTER TABLE valyria_dragon_release
ADD COLUMN target_release_id BIGINT UNSIGNED NULL DEFAULT NULL
COMMENT '回滚目标 release id' AFTER description;
