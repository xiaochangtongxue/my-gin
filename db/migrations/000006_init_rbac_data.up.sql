-- 插入默认角色
INSERT INTO `roles` (`id`, `code`, `name`, `description`, `status`, `built_in`, `sort_order`, `created_at`, `updated_at`) VALUES
(1, 'super_admin', '超级管理员', '拥有所有权限', 1, 1, 1, NOW(3), NOW(3)),
(2, 'admin', '管理员', '系统管理员', 1, 1, 2, NOW(3), NOW(3)),
(3, 'user', '普通用户', '注册用户默认角色', 1, 1, 3, NOW(3), NOW(3)),
(4, 'guest', '访客', '未登录用户', 1, 1, 4, NOW(3), NOW(3));

-- 插入 Casbin 权限策略 (v0 = role_id)
INSERT INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`, `created_at`, `updated_at`) VALUES
-- super_admin (id=1) 所有权限
('p', '1', '/api/v1/*', '*', NOW(3), NOW(3)),

-- admin (id=2) 管理权限
('p', '2', '/api/v1/admin/*', '*', NOW(3), NOW(3)),
('p', '2', '/api/v1/users', 'read', NOW(3), NOW(3)),
('p', '2', '/api/v1/users', 'create', NOW(3), NOW(3)),
('p', '2', '/api/v1/users', 'update', NOW(3), NOW(3)),

-- user (id=3) 普通用户权限
('p', '3', '/api/v1/articles', 'read', NOW(3), NOW(3)),
('p', '3', '/api/v1/articles', 'create', NOW(3), NOW(3)),
('p', '3', '/api/v1/profile', 'read', NOW(3), NOW(3)),
('p', '3', '/api/v1/profile', 'update', NOW(3), NOW(3)),

-- guest (id=4) 访客权限
('p', '4', '/api/v1/public/*', 'read', NOW(3), NOW(3));