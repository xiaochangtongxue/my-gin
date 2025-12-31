# CLAUDE.md

此文件为 Claude Code (claude.ai/code) 提供在此代码库中工作的指导。

## 项目概述

这是一个基于 Go 1.24+ 构建的生产级 Gin 框架脚手架 (`my-gin`)，采用清洁架构模式，分层为 Handler → Service → Repository。

**模块路径**: `github.com/xiaochangtongxue/my-gin`

## 命令

### 运行应用
```bash
# 默认配置
go run cmd/server/main.go

# 指定配置文件
go run cmd/server/main.go -c configs/config.yaml

# 构建后运行
go build -o bin/server.exe ./cmd/server
./bin/server.exe
```

### 依赖管理
```bash
go mod tidy
```

### 健康检查
```bash
curl http://localhost:8080/health
curl http://localhost:8080/api/v1/ping
```

### 数据库迁移
```bash
# 迁移文件使用 embed.FS 嵌入，无需额外参数

# 执行所有待执行的迁移
# 在代码中调用 database.Up()

# 回滚最后一次迁移
# 在代码中调用 database.Down()

# 查看迁移状态
# 在代码中调用 database.Status()
```

迁移文件存放在 `db/migrations/` 目录下，编译时会自动嵌入到二进制文件中。

## 架构

### 分层架构

代码库遵循严格的关注点分离，依赖方向：`cmd` → `internal` → `pkg`

- **`cmd/server/`** - 应用程序入口
- **`internal/handler/`** - HTTP 请求处理器（控制器）
- **`internal/service/`** - 业务逻辑层
- **`internal/repository/`** - 数据访问层（DAO）
- **`internal/model/`** - 数据库模型（GORM）
- **`internal/dto/`** - 请求/响应 DTO（`dto/req`、`dto/resp`）
- **`internal/middleware/`** - Gin 中间件
- **`internal/router/`** - 路由注册

### 可复用包 (`pkg/`)

- **`pkg/config`** - 基于 Viper 的配置管理
  - 支持多环境配置（dev/test/prod）
  - 环境变量覆盖，前缀为 `APP_`
  - 启动时配置校验
  - 配置结构定义在 `struct.go`

- **`pkg/logger`** - Zap 结构化日志
  - 通过 Lumberjack 实现日志轮转
  - 双输出：文件 + 控制台
  - 支持字段化结构化日志

- **`pkg/response`** - 标准化 JSON 响应封装
  - 响应格式包含：`code`、`message`、`data`、`timestamp`、`request_id`
  - 错误码定义在 `code.go`（10000+ 范围）
  - 使用 `response.Success()`、`response.Fail()`、`response.ParamError()`

- **`pkg/errors`** - 业务错误类型
  - `errors.New(code, message)` - 创建新业务错误
  - `errors.Wrap(err, code, message)` - 包装现有错误
  - 预定义错误：`ErrInvalidParam`、`ErrUnauthorized`、`ErrNotFound` 等
  - 支持通过 `WithCaller()` 追踪调用位置

- **`pkg/database`** - GORM 封装
  - 可配置的连接池
  - 通过 `database.Transaction()` 支持事务
  - 使用 `golang-migrate` 支持数据库迁移
  - `getDSN()` 从配置读取（无硬编码）

- **`pkg/cache`** - 缓存抽象层
  - 接口定义在 `interface.go`
  - Redis 和内存缓存实现
  - 统一的 `Get/Set/Del` API

- **`pkg/validator`** - 请求参数验证
  - 封装 `go-playground/validator`
  - 支持自定义验证规则

### 配置

配置通过 Viper 从 `configs/config.yaml` 加载。主要配置项：

- `server` - HTTP 服务器设置（mode、host、port、timeouts）
- `database` - MySQL/GORM 连接池设置（host、port、username、password、database、charset、parse_time、max_idle_conns、max_open_conns、conn_max_lifetime、conn_max_idle_time、slow_threshold）
- `redis` - Redis 连接设置（host、port、password、db、pool_size、min_idle_conn）
- `logger` - 日志级别和文件轮转设置（level、file_name、max_size、max_backups、max_age、compress、console）
- `jwt` - JWT 密钥和过期时间（secret、expire_time、issuer）

### 响应格式

所有 API 响应遵循以下结构：
```json
{
  "code": 0,
  "message": "success",
  "data": {},
  "timestamp": 1700000000,
  "request_id": "xxx"
}
```

错误码（定义在 `pkg/response/code.go`）：
- `0` - 成功
- `10001` - 参数错误
- `10002` - 未授权
- `10003` - 禁止访问
- `10004` - 资源不存在
- `10005` - 服务器错误
- `10006` - 数据库错误
- `10007` - Redis 错误
- `10008` - Token 过期
- `10009` - Token 无效

### 错误处理

使用 `pkg/errors` 处理业务错误：
```go
// 创建新错误
err := errors.New(response.CodeInvalidParam, "用户名不能为空")

// 包装现有错误
err := errors.Wrap(dbErr, response.CodeDBError, "查询用户失败")

// 使用预定义错误
if user == nil {
    return nil, errors.ErrNotFound
}

// 添加调用位置信息便于调试
err := errors.Wrap(dbErr, response.CodeDBError, "查询失败").WithCaller()
```

## 代码规范

- 遵循 `gofmt` 格式化
- 包名使用小写
- 接口名以 `-er` 后缀结尾
- 函数注释使用 Godoc 格式
- 分支命名：`feature/`、`bugfix/`、`hotfix/`
- 提交信息遵循 Conventional Commits 规范
- **禁止魔法值** - 始终从配置读取

## 开发阶段

项目按 10 个开发阶段组织（详见 `docs/项目需求文档.md`）：

1. **基础设施**（已完成）- config、logger、response、errors
2. **数据访问**（已完成）- database、cache、validator
3. **核心中间件**（进行中）- recovery、request ID、auth、security
4. **分层架构实现** - handler/service/repository 实现
5. **可观测性** - Swagger、健康检查、metrics
6. **权限保护** - RBAC、限流、防暴力破解
7. **工具包** - 加密、ID 生成器、时间/字符串工具
8. **异步任务** - worker pools、定时任务、文件处理
9. **容器化** - Docker、K8s 配置
10. **CI/CD** - GitLab CI、测试、文档

## 模块状态

当前完成度：
- [x] 阶段 1：基础设施（config、logger、response、errors）
- [x] 阶段 2：数据访问（database、cache、validator）
- [ ] 阶段 3-10：进行中

## 重要说明

- 始终使用 `config.Get()` 访问配置 - 禁止硬编码值
- 数据库迁移使用 `golang-migrate` - 迁移文件应放在 `db/migrations/`
- 应用在 SIGINT/SIGTERM 信号时优雅关闭
- 数据库和 Redis 初始化失败仅记录日志，不会阻止启动（fail-soft 设计）