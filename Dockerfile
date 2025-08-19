# 多阶段构建
FROM golang:1.23 AS builder

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /data/
COPY . .
RUN go build -o dingtalk

# 使用极小的运行时镜像
FROM alpine:latest
WORKDIR /data/
COPY --from=builder /data/dingtalk .
COPY config.yaml .
EXPOSE 8080

# 设置可执行权限（可选但保险）
# RUN chmod +x /data/dingtalk

CMD ["/data/dingtalk", "--port", "8080"]
