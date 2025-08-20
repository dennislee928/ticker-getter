# Ticker-Getter 票務購買平台

Ticker-Getter 是一個功能完整的票務購買平台，提供活動發布、票券銷售和票券管理功能。平台採用前後端分離架構，前端使用 Next.js，後端使用 Go (Gin) 開發。

## 功能特點

- **使用者前台**：瀏覽活動、購買票券、查看訂單歷史
- **管理員後台**：活動管理、票種管理、訂單管理、使用者管理
- **TLS 指紋識別**：防止重複購買、減少票券黃牛
- **Redis 整合**：快取和限流功能，提高系統效能和穩定性
- **Supabase 資料庫**：PostgreSQL 資料庫服務，高效可靠

## 技術架構

### 後端 (Backend)

- **語言框架**：Go + Gin
- **ORM**：GORM
- **資料庫**：PostgreSQL (Supabase)
- **快取**：Redis
- **認證**：JWT
- **API 文檔**：Swagger
- **容器化**：Docker

### 前端 (Frontend)

- **框架**：Next.js (React)
- **類型系統**：TypeScript
- **UI 框架**：TailwindCSS
- **狀態管理**：Zustand
- **HTTP 客戶端**：Axios
- **表單處理**：React Hook Form

## 專案結構

```
ticker-getter/
├── backend/                # 後端 Go 專案
│   ├── cmd/                # 入口點和配置
│   ├── config/             # 應用配置
│   ├── database/           # 資料庫遷移
│   ├── internal/           # 內部包
│   │   ├── controllers/    # HTTP 控制器
│   │   ├── dto/            # 數據傳輸對象
│   │   ├── middleware/     # 中間件
│   │   ├── models/         # 資料模型
│   │   ├── repositories/   # 資料存取層
│   │   ├── services/       # 業務邏輯層
│   │   └── vo/             # 視圖對象
│   ├── pkg/                # 可重用包
│   │   ├── cache/          # 快取工具
│   │   ├── limiter/        # 限流工具
│   │   └── utils/          # 通用工具
│   └── tests/              # 測試
│       ├── integration/    # 整合測試
│       └── unit/           # 單元測試
│
├── frontend/               # 前端 Next.js 專案
│   ├── public/             # 靜態資源
│   └── src/                # 源代碼
│       ├── api/            # API 客戶端
│       ├── app/            # Next.js App Router
│       ├── components/     # React 組件
│       ├── hooks/          # 自定義 React Hooks
│       └── lib/            # 通用工具和類型
│
├── docker-compose.yml      # Docker 配置
└── render.yaml             # Render 部署配置
```

## 安裝與運行

### 必要條件

- Go 1.21+
- Node.js 20+
- Docker & Docker Compose
- PostgreSQL 15+
- Redis 7+

### 本地開發

#### 後端

```bash
# 進入後端目錄
cd backend

# 安裝依賴
go mod download

# 設置環境變數
cp .env.example .env
# 編輯 .env 文件填入必要配置

# 啟動開發環境資料庫和 Redis
make dev

# 運行資料庫遷移
make migrate-up

# 啟動後端服務
make run
```

#### 前端

```bash
# 進入前端目錄
cd frontend

# 安裝依賴
npm install

# 設置環境變數
cp .env.example .env.local
# 編輯 .env.local 文件填入必要配置

# 啟動開發伺服器
npm run dev
```

### 使用 Docker Compose

```bash
# 構建並啟動所有服務
docker-compose up -d

# 停止所有服務
docker-compose down
```

## 部署

本專案配置為使用 Render 平台部署。

1. 在 Render 上註冊並設置專案
2. 使用 render.yaml 文件自動配置服務
3. 連接 GitHub 儲存庫進行自動部署

## 測試

```bash
# 進入後端目錄
cd backend

# 執行所有測試
./tests/run_tests.sh

# 執行單元測試
go test -v ./tests/unit/...

# 執行整合測試
go test -v ./tests/integration/...
```

## 授權協議

本專案採用 MIT 授權協議 - 詳情請參閱 [LICENSE](LICENSE) 文件。

## 聯絡方式

如有任何問題或建議，請聯繫專案維護者。
