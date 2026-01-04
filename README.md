# my-gin

生产级 Gin 框架脚手架

## 快速开始

### 安装依赖

```bash
go mod tidy
```

### 运行服务

```bash
# 默认配置
go run cmd/server/main.go

# 指定配置文件
go run cmd/server/main.go -c configs/config.yaml

# 编译后运行
go build -o bin/server.exe ./cmd/server
./bin/server.exe
```

### 测试接口

```bash
# 健康检查
curl http://localhost:8080/health

# Ping 接口
curl http://localhost:8080/api/v1/ping
```

## 项目结构

```
my-gin/
├── cmd/server/          # 应用入口
├── internal/            # 私有代码
│   ├── handler/         # HTTP 处理器
│   ├── service/         # 业务逻辑
│   ├── repository/      # 数据访问
│   ├── model/           # 数据模型
│   ├── dto/             # 数据传输对象
│   ├── middleware/      # 中间件
│   └── router/          # 路由配置
├── pkg/                 # 可复用库
│   ├── config/          # 配置管理
│   ├── logger/          # 日志系统
│   ├── response/        # 统一响应
│   └── errors/          # 错误处理
├── configs/             # 配置文件
└── docs/                # 项目文档
```

## 配置说明

配置文件位于 `configs/config.yaml`，支持环境变量覆盖（前缀 `APP_`）。

### 环境变量

| 变量名 | 说明 | 示例值 | 必填 |
|--------|------|--------|------|
| `APP_JWT_SECRET` | JWT 签名密钥（至少32位） | `your-256-bit-secret-key-here` | **是** |

### 启动示例

```bash
# 设置 JWT 密钥后启动
export APP_JWT_SECRET="your-256-bit-secret-key-here"
go run cmd/server/main.go

# Windows PowerShell
$env:APP_JWT_SECRET="your-256-bit-secret-key-here"
go run cmd/server/main.go
```

**重要**: JWT 密钥必须通过环境变量设置，切勿写入配置文件或提交到代码仓库。

## 开发进度

- [x] 阶段 1：基础设施层
- [x] 阶段 2：数据访问层
- [ ] 阶段 3：核心中间件
- [ ] 阶段 4：分层架构实现
- [ ] 阶段 5：API 文档与健康检查
- [ ] 阶段 6：权限与保护
- [ ] 阶段 7：工具包与加密
- [ ] 阶段 8：容器化与部署
- [ ] 阶段 9：CI/CD、测试与文档

详细规划请查看 [docs/项目需求文档.md](docs/项目需求文档.md)