-- 创建角色表
CREATE TABLE IF NOT EXISTS `roles` (
    `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '主键ID',
    `code` VARCHAR(50) NOT NULL COMMENT '角色编码',
    `name` VARCHAR(50) NOT NULL COMMENT '角色名称',
    `description` VARCHAR(255) DEFAULT '' COMMENT '角色描述',
    `status` TINYINT(1) NOT NULL DEFAULT 1 COMMENT '状态：1启用 0禁用',
    `built_in` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否内置角色：1是 0否',
    `sort_order` INT NOT NULL DEFAULT 0 COMMENT '排序',
    `created_at` DATETIME(3) DEFAULT NULL COMMENT '创建时间',
    `updated_at` DATETIME(3) DEFAULT NULL COMMENT '更新时间',
    `deleted_at` DATETIME(3) DEFAULT NULL COMMENT '删除时间（软删除）',
    UNIQUE KEY `uk_code` (`code`),
    KEY `idx_status` (`status`),
    KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='角色表';

-- 创建用户角色关联表
CREATE TABLE IF NOT EXISTS `user_roles` (
    `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '主键ID',
    `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    `role_id` BIGINT UNSIGNED NOT NULL COMMENT '角色ID',
    `created_at` DATETIME(3) DEFAULT NULL COMMENT '创建时间',
    `updated_at` DATETIME(3) DEFAULT NULL COMMENT '更新时间',
    `deleted_at` DATETIME(3) DEFAULT NULL COMMENT '删除时间（软删除）',
    UNIQUE KEY `uk_user_role` (`user_id`, `role_id`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_role_id` (`role_id`),
    KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户角色关联表';

-- 创建 Casbin 策略表
CREATE TABLE IF NOT EXISTS `casbin_rule` (
    `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '主键ID',
    `ptype` VARCHAR(255) NOT NULL DEFAULT '' COMMENT '策略类型：p=权限策略, g=角色继承',
    `v0` VARCHAR(255) NOT NULL DEFAULT '' COMMENT 'role_id(权限) 或 user_uid(继承)',
    `v1` VARCHAR(255) NOT NULL DEFAULT '' COMMENT 'resource(权限) 或 role_id(继承)',
    `v2` VARCHAR(255) NOT NULL DEFAULT '' COMMENT 'action',
    `v3` VARCHAR(255) NOT NULL DEFAULT '',
    `v4` VARCHAR(255) NOT NULL DEFAULT '',
    `v5` VARCHAR(255) NOT NULL DEFAULT '',
    `created_at` DATETIME(3) DEFAULT NULL COMMENT '创建时间',
    `updated_at` DATETIME(3) DEFAULT NULL COMMENT '更新时间',
    `deleted_at` DATETIME(3) DEFAULT NULL COMMENT '删除时间（软删除）',
    KEY `idx_ptype` (`ptype`),
    KEY `idx_ptype_v0` (`ptype`, `v0`),
    KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Casbin策略表';