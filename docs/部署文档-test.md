# 部署文档

## 环境信息

- **项目**: my-gin
- **部署环境**: Docker 容器
- **宿主机**: Debian 12
- **目标网络**: dev-network

## 已部署服务

| 服务 | 容器名 | 端口 | 网络 | 说明 |
|------|--------|------|------|------|
| mysql-dev | mysql:8.0 | 3306 | dev-network | MySQL 数据库 |
| redis-dev | redis:latest | 6379 | dev-network | Redis 缓存 |
| my-gin-app | my-gin | 8080 | dev-network | 应用服务 |

## 部署步骤

### 1. 准备工作

```bash
# 项目目录
cd /home/script/my-gin

# 创建 .env 文件
cat > .env << 'EOF'
APP_ENV=dev
APP_JWT_SECRET=dev-secret-key-at-least-32-characters-long-for-testing
APP_DATABASE_HOST=mysql-dev
APP_DATABASE_PORT=3306
APP_DATABASE_USERNAME=root
APP_DATABASE_PASSWORD=123456
APP_DATABASE_NAME=my_app_db
APP_REDIS_HOST=redis-dev:6379
APP_REDIS_PASSWORD=
APP_SERVER_MODE=debug
EOF
```

### 2. 构建并启动应用

```bash
# 启动应用（加入 dev-network）
docker compose -f build/docker/docker-compose.yml --env-file .env up -d --build

# 查看日志
docker compose -f build/docker/docker-compose.yml logs -f app
```

### 3. 验证服务

```bash
# 健康检查
curl http://localhost:8080/health

# Metrics 端点
curl http://localhost:8080/metrics

# 查看容器状态
docker ps | grep my-gin
```

## 网络架构

```
┌─────────────────────────────────────────────────────────────┐
│                      dev-network                             │
├─────────────────┬─────────────────┬─────────────────────────────────┤
│   mysql-dev     │   redis-dev     │      my-gin-app            │
│   :3306         │   :6379         │      :8080                 │
└─────────────────┴─────────────────┴─────────────────────────────────┘
                           │
                           │ 端口映射
                           ▼
                    ┌─────────────────┐
                    │   宿主机        │
                    │   localhost:8080│
                    └─────────────────┘
```

## Prometheus 监控配置

### prometheus.yml

```yaml
scrape_configs:
  - job_name: 'prometheus'
    static_configs:
      - targets: ['localhost:9090']

  - job_name: 'my-gin-app'
    metrics_path: '/metrics'
    scrape_interval: 15s
    static_configs:
      - targets: ['172.17.0.1:8080']  # Docker 网桥网关
```

### 启动监控服务

```bash
# 启动 prometheus 和 grafana
docker compose -f prometheus-compose.yml up -d

# 访问
# Prometheus: http://localhost:9090
# Grafana: http://localhost:3000 (admin/admin)
```

## 常用命令

```bash
# 启动
docker compose -f build/docker/docker-compose.yml --env-file .env up -d

# 停止
docker compose -f build/docker/docker-compose.yml down

# 查看日志
docker compose -f build/docker/docker-compose.yml logs -f app

# 重新构建
docker compose -f build/docker/docker-compose.yml up -d --build

# 进入容器
docker exec -it my-gin-app sh

# 查看容器状态
docker ps -a | grep my-gin
```

## 故障排查

### 问题：容器无法连接 MySQL/Redis

**原因**：不在同一个网络中

**解决**：
```bash
# 确认服务在 dev-network 中
docker network inspect dev-network

# 将服务加入网络
docker network connect dev-network my-gin-app
```

### 问题：外部无法访问服务

**检查**：
```bash
# 检查端口映射
docker port my-gin-app

# 检查容器状态
docker ps | grep my-gin

# 检查日志
docker logs my-gin-app --tail 50
```

### 问题：metrics 端点无法访问

**检查**：
```bash
# 容器内测试
docker exec my-gin-app wget -O- http://localhost:8080/metrics

# 检查配置
grep -A 5 "auth:" configs/config.yaml | grep metrics
```

## 文件清单

| 文件 | 说明 |
|------|------|
| `build/docker/Dockerfile` | 镜像构建文件 |
| `build/docker/docker-compose.yml` | 应用编排文件 |
| `pkg/permission/model.conf` | Casbin 模型文件（需复制到镜像） |
| `.env` | 环境变量配置（不提交） |
| `.env.example` | 环境变量模板 |

## 下一步

- [ ] 添加 Grafana 面板配置
- [ ] 配置告警规则
- [ ] 添加日志聚合（ELK/Loki）
- [ ] 完善监控指标