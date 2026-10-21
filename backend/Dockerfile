# 构建阶段
# 与 go.mod 的 go 指令保持一致
FROM golang:1.25-alpine AS builder

# 设置工作目录
WORKDIR /app

# 安装必要的工具
RUN apk add --no-cache git ca-certificates tzdata

# 复制go mod文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o main cmd/main.go

# 运行阶段
FROM alpine:latest

# 安装必要的运行时依赖
RUN apk --no-cache add ca-certificates tzdata sqlite

# 设置工作目录
WORKDIR /app

# 从构建阶段复制可执行文件
COPY --from=builder /app/main .

# 复制配置文件
COPY --from=builder /app/configs ./configs

# 创建数据目录
RUN mkdir -p /app/data

# 设置时区
ENV TZ=Asia/Shanghai

# 容器内需要对外（供 nginx / 端口映射访问）时显式监听 0.0.0.0；二进制默认只监听 127.0.0.1
ENV PORT=8080
ENV LISTEN_ADDR=0.0.0.0
# JWT_SECRET 必须在运行时注入（至少 32 个字符），镜像内不含默认值

# 暴露端口
EXPOSE 8080

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget -q -O /dev/null --tries=1 http://127.0.0.1:8080/api/health || exit 1

# 启动应用
CMD ["./main"]