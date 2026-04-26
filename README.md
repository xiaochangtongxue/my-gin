# my-gin

生产级 Gin 后端脚手架，面向“持续新增业务接口 + 前端稳定对接”的场景。项目基于 Go、Gin、GORM、Wire、Viper、Zap、Redis、Casbin 和 Prometheus，提供认证、RBAC、限流、健康检查、Swagger/OpenAPI、统一响应和业务错误封装。

## 当前架构

项目已经从单纯分层目录整理为“应用容器 + 垂直模块注册层 + 共享基础包”的结构：

```text
my-gin/
├── cmd/server/              # 入口、Wire 注入定义与生成文件
├── internal/app/            # App 容器与 RouteModule 接口
├── internal/modules/        # 业务模块注册层
│   ├── auth/                # 认证、验证码路由与 provider set
│   ├── docs/                # Swagger 路由
│   ├── health/              # 健康检查路由
│   ├── metrics/             # Prometheus 路由
│   ├── permission/          # RBAC 管理路由
│   └── security/            # 解锁、防暴力破解管理路由
├── internal/handler/        # HTTP Handler
├── internal/service/        # 业务逻辑
├── internal/repository/     # 数据访问
├── internal/model/          # GORM 模型
├── internal/dto/            # 请求/响应 DTO
├── internal/middleware/     # Gin 中间件
├── internal/permission/     # Casbin RBAC 权限实现与内嵌模型
├── pkg/                     # 可复用基础设施包
├── api/                     # Swagger/OpenAPI 生成产物
├── configs/                 # 配置文件
├── db/migrations/           # 嵌入式数据库迁移
└── docs/                    # 项目文档
```

新增业务模块时，优先在 `internal/modules/<module>` 增加 `RouteModule` 和 Wire `ProviderSet`，再接入 `cmd/server/wire.go` 的 `moduleSet` 与 `ProvideRouteModules`。现阶段 Handler/Service/Repository 仍在共享分层目录，后续模块可逐步垂直迁移。

## 快速开始

准备 MySQL、Redis，并设置必填环境变量：

```powershell
$env:APP_JWT_SECRET="your-256-bit-secret-key-here-32chars"
$env:APP_DATABASE_HOST="127.0.0.1"
$env:APP_DATABASE_PORT="3306"
$env:APP_DATABASE_USERNAME="root"
$env:APP_DATABASE_PASSWORD="your-password"
$env:APP_DATABASE_NAME="my_app_db"
$env:APP_REDIS_HOST="127.0.0.1"
$env:APP_REDIS_PORT="6379"
```

运行服务：

```bash
go mod tidy
go run cmd/server/main.go -c configs/config.yaml
```

构建：

```bash
go build -o bin/server.exe ./cmd/server
./bin/server.exe -c configs/config.yaml
```

外部依赖采用严格启动策略：配置、数据库、Redis、RBAC 初始化失败会直接启动失败；数据库迁移失败也会终止启动。

## 常用命令

```bash
# 全量测试
go test ./...

# Wire 依赖注入生成
go run github.com/google/wire/cmd/wire ./cmd/server

# Swagger/OpenAPI 生成
go run github.com/swaggo/swag/cmd/swag init -g cmd/server/main.go -o api

# 健康检查
curl http://localhost:8080/health
curl http://localhost:8080/health/ready

# Swagger
curl http://localhost:8080/swagger/index.html
```

## API 路径

认证接口统一在 `/api/v1/auth` 下：

- `GET /api/v1/auth/captcha`
- `POST /api/v1/auth/register`
- `POST /api/v1/auth/login`
- `POST /api/v1/auth/refresh`
- `POST /api/v1/auth/logout`

管理接口继续使用 `/api/v1/admin/*`，当前包含角色、权限、用户角色、账号解锁等接口。

前端契约见 [docs/api-contract.md](docs/api-contract.md)，OpenAPI 生成产物见 `api/swagger.yaml`。

## 响应与错误

所有业务响应使用统一 envelope：

```json
{
  "code": 0,
  "message": "success",
  "data": {},
  "timestamp": 1700000000,
  "request_id": "xxx"
}
```

Service 层应返回 `pkg/errors.BusinessError`，Handler 层统一使用 `response.Error(c, err)`。新增错误码时需要同步：

1. `pkg/response/code.go` 的常量
2. `codeMessages` 映射
3. `pkg/response/response.go` 的 HTTP 状态码映射
4. 前端契约文档中的错误码表

## 配置说明

配置从 `configs/config.yaml` 读取，支持 `APP_` 环境变量覆盖。重要变量：

| 变量名 | 说明 | 必填 |
|---|---|---|
| `APP_JWT_SECRET` | JWT 签名密钥，至少 32 位 | 是 |
| `APP_DATABASE_HOST` | MySQL 主机 | 是 |
| `APP_DATABASE_PORT` | MySQL 端口 | 是 |
| `APP_DATABASE_USERNAME` | MySQL 用户名 | 是 |
| `APP_DATABASE_PASSWORD` | MySQL 密码 | 是 |
| `APP_DATABASE_NAME` | MySQL 数据库名 | 是 |
| `APP_REDIS_HOST` | Redis 主机，可不带端口 | 是 |
| `APP_REDIS_PORT` | Redis 端口 | 是 |

`APP_DATABASE_DATABASE` 作为兼容别名仍可使用，但推荐统一使用 `APP_DATABASE_NAME`。

## 文档

- [docs/api-contract.md](docs/api-contract.md) - 前端 API 对接契约
- [docs/项目需求文档.md](docs/项目需求文档.md) - 项目需求文档
