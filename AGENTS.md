# AGENTS.md

此文件为 Codex 在本仓库中工作提供指导。

## 项目概览

`my-gin` 是一个生产级 Gin 后端脚手架，模块路径为：

```text
github.com/xiaochangtongxue/my-gin
```

当前架构是“应用容器 + 垂直模块注册层 + 共享分层业务包”：

- `cmd/server/`：应用入口、Wire 注入定义和生成文件
- `internal/app/`：`App` 容器和 `RouteModule` 接口
- `internal/modules/`：模块级路由注册与 Wire ProviderSet
- `internal/handler/`、`internal/service/`、`internal/repository/`：现有共享分层业务代码
- `internal/model/`、`internal/dto/`：GORM 模型和请求/响应 DTO
- `internal/middleware/`：Gin 中间件
- `internal/permission/`：Casbin RBAC 实现，默认模型通过 `embed` 内嵌
- `pkg/`：配置、日志、响应、错误、数据库、缓存、JWT、限流、指标等基础包
- `api/`：Swagger/OpenAPI 生成产物
- `docs/api-contract.md`：前端 API 契约

新增业务模块时，优先新增 `internal/modules/<name>/module.go`，实现 `RegisterRoutes(*gin.Engine)`，提供 Wire `ProviderSet`，再接入 `cmd/server/wire.go` 的 `moduleSet` 和 `ProvideRouteModules`。不要把新路由散落到 `cmd/server/main.go`。

## 常用命令

```bash
# 运行
go run cmd/server/main.go -c configs/config.yaml

# 构建
go build -o bin/server.exe ./cmd/server

# 测试
go test ./...

# Wire 生成
go run github.com/google/wire/cmd/wire ./cmd/server

# Swagger/OpenAPI 生成
go run github.com/swaggo/swag/cmd/swag init -g cmd/server/main.go -o api
```

健康检查：

```bash
curl http://localhost:8080/health
curl http://localhost:8080/health/ready
curl http://localhost:8080/swagger/index.html
```

## 配置与启动

配置由 `pkg/config` 基于 Viper 从 `configs/config.yaml` 加载，支持 `APP_` 环境变量覆盖。

关键环境变量：

- `APP_JWT_SECRET`：必填，至少 32 位
- `APP_DATABASE_HOST`
- `APP_DATABASE_PORT`
- `APP_DATABASE_USERNAME`
- `APP_DATABASE_PASSWORD`
- `APP_DATABASE_NAME`
- `APP_REDIS_HOST`
- `APP_REDIS_PORT`

数据库、Redis、RBAC 初始化采用严格启动策略：失败即返回错误并终止启动。不要引入本地内存缓存降级启动，除非用户明确要求改变策略。

## API 路径

认证接口统一使用 `/api/v1/auth/*`：

- `GET /api/v1/auth/captcha`
- `POST /api/v1/auth/register`
- `POST /api/v1/auth/login`
- `POST /api/v1/auth/refresh`
- `POST /api/v1/auth/logout`

管理接口继续使用 `/api/v1/admin/*`。前端契约见 `docs/api-contract.md`。

## 错误处理规范

Service 层返回 `pkg/errors.BusinessError`：

```go
return nil, errors.New(response.CodeInvalidParam, "参数错误")
return nil, errors.Wrap(err, response.CodeDBError, "查询失败")
```

Handler 层统一：

```go
if err != nil {
    response.Error(c, err)
    return
}
```

不要在 Handler 中通过 `err.Error()` 字符串比较判断业务错误。

新增错误码时必须同步：

1. `pkg/response/code.go` 常量
2. `codeMessages`
3. `pkg/response/response.go` 的 HTTP 状态码映射
4. `docs/api-contract.md` 错误码表

## Wire 规范

- `cmd/server/wire.go` 保留 `infraSet`、`moduleSet` 和 `ProvideRouteModules`
- 模块自己的依赖放到 `internal/modules/<name>/module.go`
- 修改 ProviderSet 后运行 `go run github.com/google/wire/cmd/wire ./cmd/server`
- `cmd/server/wire_gen.go` 是生成文件，不手改

## 文档与 Swagger

- Swagger 生成目录统一为 `api/`
- 不再使用 `cmd/server/docs/` 或 `docs/swagger.*`
- 修改路由注解后运行 Swagger 生成命令
- 前端对接说明维护在 `docs/api-contract.md`

## 代码规范

- 遵循 `gofmt`
- 保持改动聚焦，避免顺手重构无关代码
- 配置值从 `config.Config` 读取，禁止新增魔法值
- 新功能要补针对性测试，至少保证 `go test ./...` 通过
- 不要回滚用户已有未提交改动
