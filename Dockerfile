# FILEPATH:/Users/yangyang/go/src/github.com/github-wyy/m-scheduler-extender/Dockerfile
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o m-scheduler-extender .

FROM alpine:3.18
WORKDIR /app
# 创建证书目录
RUN mkdir -p /app/ssl
# 复制可执行文件
COPY --from=builder /app/m-scheduler-extender .
# 复制证书文件
COPY --from=builder /app/ssl/cert.pem /app/ssl/cert.pem
COPY --from=builder /app/ssl/key.pem /app/ssl/key.pem
# 暴露HTTP和HTTPS端口
EXPOSE 8010 8443
ENTRYPOINT ["./m-scheduler-extender"]
