-- 删除 Casbin 权限策略
DELETE FROM `casbin_rule` WHERE `ptype` = 'p' AND `v0` IN ('1', '2', '3', '4');

-- 删除默认角色
DELETE FROM `roles` WHERE `id` IN (1, 2, 3, 4);