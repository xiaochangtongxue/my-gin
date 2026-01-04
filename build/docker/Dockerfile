# 多阶段构建 Dockerfile
# 阶段1: 构建阶段
FROM golang:1.24-alpine AS builder

# 安装构建依赖
RUN apk add --no-cache git ca-certificates tzdata

# 设置工作目录
WORKDIR /build

# 复制 go.mod 和 go.sum 并下载依赖（利用 Docker 缓存）
COPY go.mod go.sum ./

# 设置 GOPROXY 加速（国内强烈建议）
ENV GOPROXY=https://goproxy.cn,direct
RUN go mod download

# 复制源代码
COPY . .

# 构建参数
ARG VERSION=dev
ARG BUILD_TIME
ARG GIT_COMMIT

# 构建应用（静态链接，减小镜像体积）
# -ldflags 参数注入版本信息
# -w -s 去除调试信息减小二进制大小
# -tags netgo 禁用 cgo 实现静态链接
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s -X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME} -X main.GitCommit=${GIT_COMMIT}" \
    -tags netgo \
    -o /build/server ./cmd/server

# 验证二进制文件
RUN /build/server --version || true

# ========================================
# 阶段2: 运行阶段
# ========================================
FROM alpine:3.19

# 安装运行时依赖
# ca-certificates: HTTPS 请求需要
# tzdata: 时区数据
RUN apk --no-cache add ca-certificates tzdata && \
    # 设置时区为上海
    ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone && \
    # 创建非 root 用户
    addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

# 设置工作目录
WORKDIR /app



# 从构建阶段复制二进制文件
COPY --from=builder /build/server .

# 复制配置文件
COPY configs/ ./configs/

# 复制 Casbin 模型文件（运行时需要读取）
COPY pkg/permission/model.conf ./pkg/permission/model.conf

# 创建日志目录
RUN mkdir -p logs && \
    chown -R appuser:appuser /app

# 切换到非 root 用户
USER appuser

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=10s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# 暴露端口
EXPOSE 8080

# 设置停止信号为 SIGTERM（优雅关闭）
STOPSIGNAL SIGTERM

# 启动应用
# 配置文件路径可通过环境变量 CONFIG_PATH 覆盖
CMD ["./server", "-c", "configs/config.yaml"]