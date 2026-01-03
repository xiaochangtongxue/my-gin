-- 修改用户表，添加 uid 和 mobile 字段
-- 注意：如果已有数据，需要先处理数据迁移

-- 添加 uid 字段（对外展示的唯一标识）
SET @col_exists = (SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'users' AND COLUMN_NAME = 'uid');
SET @sql = IF(@col_exists = 0,
    'ALTER TABLE `users` ADD COLUMN `uid` BIGINT UNSIGNED NULL UNIQUE COMMENT ''对外用户ID（14位纯数字）'' AFTER `id`',
    'SELECT ''Column uid already exists''');
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- 添加 uid 索引
SET @index_exists = (SELECT COUNT(*) FROM INFORMATION_SCHEMA.STATISTICS
    WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'users' AND INDEX_NAME = 'idx_uid');
SET @sql = IF(@index_exists = 0,
    'ALTER TABLE `users` ADD INDEX `idx_uid` (`uid`)',
    'SELECT ''Index idx_uid already exists''');
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- 添加 mobile 字段（手机号）
SET @col_exists = (SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'users' AND COLUMN_NAME = 'mobile');
SET @sql = IF(@col_exists = 0,
    'ALTER TABLE `users` ADD COLUMN `mobile` CHAR(11) NULL UNIQUE COMMENT ''手机号'' AFTER `uid`',
    'SELECT ''Column mobile already exists''');
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- 添加 mobile 索引
SET @index_exists = (SELECT COUNT(*) FROM INFORMATION_SCHEMA.STATISTICS
    WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'users' AND INDEX_NAME = 'idx_mobile');
SET @sql = IF(@index_exists = 0,
    'ALTER TABLE `users` ADD INDEX `idx_mobile` (`mobile`)',
    'SELECT ''Index idx_mobile already exists''');
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- 修改 username 字段长度和注释
ALTER TABLE `users` MODIFY COLUMN `username` VARCHAR(20) NOT NULL COMMENT '用户名';

-- 如果已有数据，需要为每个用户生成 uid
-- UPDATE `users` SET `uid` = 10000000000000 + id WHERE `uid` IS NULL;
-- ALTER TABLE `users` MODIFY COLUMN `uid` BIGINT UNSIGNED NOT NULL;