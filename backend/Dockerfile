# 使用 golang 1.21 作為構建基礎鏡像
FROM golang:1.21-alpine AS builder

# 安裝基本工具
RUN apk update && apk add --no-cache git ca-certificates tzdata && update-ca-certificates

# 設置工作目錄
WORKDIR /app

# 複製 go mod 和 sum 文件
COPY go.mod go.sum ./

# 下載依賴
RUN go mod download

# 複製源代碼
COPY . .

# 構建應用程序
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# 使用輕量級 alpine 鏡像
FROM alpine:latest

# 從構建階段複製相關文件
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/main /app/main
COPY --from=builder /app/.env.example /app/.env

# 設置工作目錄
WORKDIR /app

# 設置時區
ENV TZ=Asia/Taipei

# 暴露端口
EXPOSE 8080

# 運行應用程序
CMD ["./main"]
