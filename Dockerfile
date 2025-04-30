# ---------- Builder Stage ----------
FROM golang:1.23-alpine AS builder
WORKDIR /app

# 完整复制后端源码，包括go.mod
COPY backend/ ./

# 添加必要的构建工具
RUN apk add --no-cache git

# 下载依赖
RUN go mod download

# 构建静态二进制
RUN CGO_ENABLED=0 GOOS=linux go build -o nya-bird-lg-go .

# ---------- Runtime Stage ----------
FROM alpine:3.18 AS runtime
WORKDIR /app

# 安装运行时依赖、创建用户和目录
RUN apk add --no-cache ca-certificates tzdata && \
    addgroup -S appgroup && adduser -S appuser -G appgroup && \
    mkdir -p /app/config && chown appuser:appgroup /app/config

# 切换到非 root 用户
USER appuser

# 复制二进制与说明文档
COPY --from=builder --chown=appuser:appgroup /app/nya-bird-lg-go ./

VOLUME /app/config

EXPOSE 5000 8000

ENV BIRDLG_HOST=0.0.0.0 \
    BIRDLG_LISTEN=5000 \
    BIRDLG_PROXY_PORT=8000

ENTRYPOINT ["./nya-bird-lg-go"]