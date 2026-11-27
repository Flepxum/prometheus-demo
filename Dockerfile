
# 使用构建机器本身的架构编译 Go
FROM --platform=$BUILDPLATFORM golang:1.26.5-alpine AS builder

ARG TARGETOS
ARG TARGETARCH

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# 根据目标架构编译
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -o prometheus-demo .

# 最终运行环境
FROM alpine:3.22

WORKDIR /app

COPY --from=builder /app/prometheus-demo .

EXPOSE 8080

CMD ["./prometheus-demo"]
