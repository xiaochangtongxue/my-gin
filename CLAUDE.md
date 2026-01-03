# CLAUDE.md

此文件为 Claude Code (claude.ai/code) 提供在此代码库中工作的指导。

## 项目概述

这是一个基于 Go 1.21+ 构建的生产级 Gin 框架脚手架 (`my-gin`)，采用清洁架构模式，分层为 Handler → Service → Repository。

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

### Wire 依赖注入生成
```bash
# 修改 wire.go 后重新生成
wire gen ./cmd/server
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
  - 支持自定义验证规则（mobile、password 等）

- **`pkg/jwt`** - JWT 认证管理
  - 支持双 Token 机制（Access Token + Refresh Token）
  - Token 生成、解析、验证
  - 支持自定义 Claims（包含 UserID、Username）

- **`pkg/limiter`** - 限流器抽象层
  - 接口定义在 `limiter.go`
  - 内存限流实现（令牌桶算法，基于 `golang.org/x/time/rate`）
  - Redis 分布式限流实现
  - 支持 Global/IP/User 三级限流

- **`pkg/metrics`** - Prometheus 指标收集
  - HTTP 请求指标（RequestTotal、RequestDuration、RequestsInFlight）
  - 数据库指标（连接数、查询耗时）
  - Redis 指标（命令总数、请求耗时）
  - 应用指标（启动时间、版本信息）

- **`pkg/captcha`** - 验证码生成
  - 基于 base64Captcha 生成图形验证码
  - 支持数字/字符验证码
  - 验证码存储在缓存中，支持过期时间

- **`pkg/notify`** - 通知抽象层
  - 定义通知接口（`Notifier`）
  - 提供 NoopNotifier 空实现

- **`pkg/utils/uid.go`** - UID 生成器
  - 使用 Feistel Cipher 算法将数据库自增 ID 混淆为 14 位纯数字 UID
  - 可逆混淆，支持 UID 解析回原始 ID

### 配置

配置通过 Viper 从 `configs/config.yaml` 加载。主要配置项：

- `server` - HTTP 服务器设置（mode、host、port、timeouts）
- `database` - MySQL/GORM 连接池设置
- `redis` - Redis 连接设置
- `logger` - 日志级别和文件轮转设置
- `jwt` - JWT 过期时间（**密钥通过环境变量 APP_JWT_SECRET 读取**）
- `password` - 密码强度配置
- `captcha` - 验证码配置
- `bruteforce` - 防暴力破解配置
- `middleware` - 中间件配置

### 环境变量

| 变量名 | 说明 | 示例值 | 必填 |
|--------|------|--------|------|
| `APP_JWT_SECRET` | JWT 签名密钥（至少32位） | `your-256-bit-secret-key-here` | **是** |

**重要**: JWT 密钥必须通过环境变量设置，切勿写入配置文件或提交到代码仓库。启动前必须先设置环境变量：

```bash
export APP_JWT_SECRET="your-256-bit-secret-key-here"
go run cmd/server/main.go
```

### 依赖注入（Wire）

项目使用 Google Wire 进行依赖注入，主要文件：

- `cmd/server/wire.go` - Wire 依赖注入定义（`//go:build wireinject`）
- `cmd/server/wire_gen.go` - Wire 自动生成的代码（无需手动编辑）

**Wire ProviderSets 组织结构**：
- `infraSet` - 基础设施（Config、DB、Redis、JWT）
- `repoSet` - Repository 层
- `serviceSet` - Service 层
- `handlerSet` - Handler 层
- `engineSet` - Gin Engine 及中间件

**修改依赖后需要重新生成**：
```bash
wire gen ./cmd/server
```

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
- `10400` - 验证码错误或已过期
- `10401` - 用户名或密码错误
- `10403` - 权限不足
- `10423` - 账号已锁定
- `10429` - 请求过于频繁（限流）

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

**⚠️ 错误码开发规范（重要）**：
1. 在 `pkg/response/code.go` 中定义新错误码常量
2. 在 `pkg/response/code.go` 的 `codeMessages` 中添加错误消息映射
3. **在 `pkg/response/response.go` 的 `getStatusCode()` 函数中添加 HTTP 状态码映射**（常遗漏！）
4. 业务错误码 → HTTP 状态码映射规则：
   - `1xxxx` → 400 Bad Request
   - `10002,10008,10009,10401` → 401 Unauthorized
   - `10003,10303,10423` → 403 Forbidden
   - `10004` → 404 Not Found
   - `10429` → 429 Too Many Requests
   - 其他 → 500 Internal Server Error

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
3. **核心中间件**（已完成）- recovery、request_id、logger、cors、csrf、security、auth
4. **分层架构实现**（已完成）- Wire 依赖注入、handler/service/repository
5. **可观测性**（部分完成）- 健康检查、metrics 已实现，Swagger 待实现
6. **权限保护**（进行中）- RBAC、限流（已实现）、防暴力破解
7. **工具包** - 加密、ID 生成器、时间/字符串工具
8. **异步任务** - worker pools、定时任务、文件处理
9. **容器化** - Docker、K8s 配置
10. **CI/CD** - GitLab CI、测试、文档

### 中间件执行顺序

中间件按以下顺序执行（在 `internal/middleware/middleware.go` 中定义）：
1. Recovery（必须最前，捕获 panic）
2. RequestID（尽早执行，确保后续中间件能获取）
3. Security（安全响应头）
4. CORS（跨域）
5. XSS（可选）
6. CSRF（可选）
7. Auth（JWT 认证）
8. Metrics（Prometheus 指标收集）
9. RateLimit（限流）
10. Logger（请求日志，放最后记录完整信息）

## 限流器使用

限流中间件已集成，支持三级限流：
- **全局限流**：所有请求共享配额
- **IP 限流**：按客户端 IP 限流
- **用户限流**：按认证用户 ID 限流

配置示例：
```yaml
middleware:
  ratelimit:
    enable: true
    global:
      rate: 100
      burst: 200
    ip:
      rate: 20
      burst: 40
    user:
      rate: 10
      burst: 20
```

## 用户认证与防暴力破解

### 用户表结构

用户表包含以下关键字段：
- `id` - 数据库自增主键
- `uid` - 14位纯数字（通过 Feistel Cipher 混淆生成，对外展示）
- `username` - 用户名
- `mobile` - 手机号（唯一）
- `password` - bcrypt 哈希密码

### 认证流程

```
POST /api/v1/login
{ mobile, password, captcha_id?, captcha_code? }
```

登录流程：
1. 根据 mobile 查询用户
2. 检查账号是否被锁定（Redis: `login:locked:{uid}`）
3. 检查是否需要验证码（失败 1 次后触发）
4. 验证密码（bcrypt）
5. 失败则记录次数并检查是否锁定

### 防暴力破解策略

| 策略 | 触发条件 | 锁定时间 |
|------|----------|----------|
| IP 限制 | 同一IP 5次/15分钟 | 锁定 30分钟 |
| 账号限制 | 同一账号 3次/15分钟 | 锁定 30分钟 |
| IP 黑名单 | 累计 10次失败 | 锁定 1小时 |

### Redis 数据结构

| Key | 用途 | TTL |
|-----|------|-----|
| `captcha:{id}` | 验证码答案 | 5分钟 |
| `login:need_captcha:{ip}:{uid}` | 需要验证码标记 | 15分钟 |
| `login:fail:{ip}:{uid}` | 失败次数 | 15分钟 |
| `login:locked:{uid}` | 账号锁定 | 30分钟 |
| `login:blacklist:{ip}` | IP黑名单 | 1小时 |

### 已实现 API

- `POST /api/v1/register` - 用户注册
- `POST /api/v1/login` - 用户登录
- `POST /api/v1/refresh` - 刷新 Token
- `GET /api/v1/captcha` - 获取验证码
- `POST /api/v1/logout` - 登出（Token 加入黑名单）
- `POST /api/v1/admin/unlock-account` - 管理员解锁账号（临时开放，后续需 RBAC）

## 重要说明

- 始终使用 `config.Get()` 访问配置 - 禁止硬编码值
- 数据库迁移使用 `golang-migrate` - 迁移文件应放在 `db/migrations/`
- 应用在 SIGINT/SIGTERM 信号时优雅关闭
- 数据库和 Redis 初始化失败仅记录日志，不会阻止启动（fail-soft 设计）