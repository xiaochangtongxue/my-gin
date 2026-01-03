-- 回滚：删除 uid 和 mobile 字段
ALTER TABLE `users` DROP INDEX IF EXISTS `idx_mobile`;
ALTER TABLE `users` DROP COLUMN IF EXISTS `mobile`;

ALTER TABLE `users` DROP INDEX IF EXISTS `idx_uid`;
ALTER TABLE `users` DROP COLUMN IF EXISTS `uid`;

ALTER TABLE `users` MODIFY COLUMN `username` VARCHAR(50) NOT NULL COMMENT '用户名';