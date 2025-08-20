package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/lipeichen/ticket-getter/cmd/api"
)

// @title           票務購買系統 API
// @version         1.0
// @description     票務購買平台的 RESTful API
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.example.com/support
// @contact.email  support@example.com

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description 在 Authorization header 提供 Bearer token
func main() {
	// 載入環境變數
	err := godotenv.Load()
	if err != nil {
		log.Println("未找到 .env 文件，使用環境變數")
	}

	// 啟動 API 伺服器
	api.Run()
}
