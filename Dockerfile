FROM golang:1.21-alpine AS builder

WORKDIR /app

# 复制前端构建产物
COPY backend/frontend/dist /app/frontend/dist

# 复制Go源码
COPY backend /app

# 构建后端
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o nya-bird-lg-go

# 使用最小化的镜像
FROM alpine:latest

WORKDIR /app

# 安装必要的运行时依赖
RUN apk --no-cache add ca-certificates tzdata

# 从构建阶段复制二进制文件
COPY --from=builder /app/nya-bird-lg-go /app/nya-bird-lg-go

# 复制README和配置示例
COPY README.md /app/
COPY backend/.env.sample /app/.env.sample

# 创建卷用于持久化配置
VOLUME /app/config

# 暴露应用端口
EXPOSE 5000 8000

# 设置环境变量
ENV BIRDLG_HOST=0.0.0.0 \
    BIRDLG_LISTEN=5000 \
    BIRDLG_PROXY_PORT=8000

# 运行应用
ENTRYPOINT ["/app/nya-bird-lg-go"]