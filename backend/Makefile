.PHONY: all build run clean test migrate-up migrate-down swagger

# 設置變量
APP_NAME=ticket-getter
MIGRATION_DIR=./database/migrations

all: clean build run

# 構建應用
build:
	@echo "構建應用..."
	go build -o bin/$(APP_NAME) -v .

# 運行應用
run: build
	@echo "運行應用..."
	./bin/$(APP_NAME)

# 清理構建文件
clean:
	@echo "清理構建文件..."
	rm -rf bin

# 運行測試
test:
	@echo "運行測試..."
	go test -v ./...

# 生成 Swagger 文檔
swagger:
	@echo "生成 Swagger 文檔..."
	swag init

# 運行數據庫遷移 (up)
migrate-up:
	@echo "運行數據庫遷移 (up)..."
	goose -dir $(MIGRATION_DIR) postgres "$(DATABASE_URL)" up

# 還原數據庫遷移 (down)
migrate-down:
	@echo "還原數據庫遷移 (down)..."
	goose -dir $(MIGRATION_DIR) postgres "$(DATABASE_URL)" down

# 創建新的遷移文件
migrate-create:
	@echo "創建新的遷移文件..."
	goose -dir $(MIGRATION_DIR) create $(name) sql

# 啟動開發環境 Docker 容器
dev:
	@echo "啟動開發環境..."
	docker-compose up -d postgres redis

# 停止開發環境 Docker 容器
dev-stop:
	@echo "停止開發環境..."
	docker-compose stop postgres redis

# 運行完整 Docker 環境
docker-up:
	@echo "啟動完整 Docker 環境..."
	docker-compose up -d

# 停止 Docker 環境
docker-down:
	@echo "停止 Docker 環境..."
	docker-compose down
