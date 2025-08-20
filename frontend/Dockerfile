# 使用 Node.js 20 作為構建基礎鏡像
FROM node:20-alpine AS builder

# 設置工作目錄
WORKDIR /app

# 複製 package.json 和 package-lock.json
COPY package*.json ./

# 安裝依賴
RUN npm ci

# 複製源代碼
COPY . .

# 構建應用程序
RUN npm run build

# 使用 Node.js 20 作為生產環境基礎鏡像
FROM node:20-alpine

# 設置工作目錄
WORKDIR /app

# 從構建階段複製構建結果
COPY --from=builder /app/package*.json ./
COPY --from=builder /app/.next ./.next
COPY --from=builder /app/public ./public
COPY --from=builder /app/node_modules ./node_modules

# 設置環境變量
ENV NODE_ENV=production
ENV NEXT_PUBLIC_API_URL=https://api.ticker-getter.example.com/api/v1

# 暴露端口
EXPOSE 3000

# 運行應用程序
CMD ["npm", "start"]
