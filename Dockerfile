# ---------- Builder Stage ----------
FROM golang:1.23-alpine AS builder

# 获取目标平台信息
ARG TARGETPLATFORM
ARG GOARCH
RUN echo "Building for $TARGETPLATFORM"

# 根据目标平台设置正确的 GOARCH
RUN if [ "$TARGETPLATFORM" = "linux/amd64" ]; then \
      echo "GOARCH=amd64" >> /etc/environment; \
    elif [ "$TARGETPLATFORM" = "linux/arm64" ]; then \
      echo "GOARCH=arm64" >> /etc/environment; \
    else \
      echo "Unsupported platform: $TARGETPLATFORM" && exit 1; \
    fi

WORKDIR /app

# 拷贝源码
COPY backend/ .  

# 安装构建工具、编译并压缩二进制
RUN apk add --no-cache git upx && \
    . /etc/environment && \
    CGO_ENABLED=0 GOOS=linux \
    go build -ldflags="-s -w" -o nya-bird-lg-go . && \
    upx --ultra-brute nya-bird-lg-go

# ---------- Runtime Stage ----------
FROM scratch AS runtime

# 拷贝已压缩的静态二进制
COPY --from=builder /app/nya-bird-lg-go /nya-bird-lg-go

# 以非 root（UID/GID 1000）运行
USER 1000:1000

# 暴露端口与环境变量
EXPOSE 5000 8000
ENV BIRDLG_HOST=0.0.0.0 \
    BIRDLG_LISTEN=5000 \
    BIRDLG_PROXY_PORT=8000

ENTRYPOINT ["/nya-bird-lg-go"]