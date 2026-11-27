# Go 编译环境
FROM golang:1.26.5-alpine AS builder

WORKDIR /app

# 先复制依赖文件，利用 Docker 缓存
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 编译 Linux 可执行文件
RUN CGO_ENABLED=0 GOOS=linux go build -o prometheus-demo .

# 运行环境
FROM alpine:3.22

WORKDIR /app

COPY --from=builder /app/prometheus-demo .

EXPOSE 8080

CMD ["./prometheus-demo"]