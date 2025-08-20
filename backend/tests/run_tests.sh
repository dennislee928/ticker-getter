#!/bin/bash

# 運行單元測試
echo "正在運行單元測試..."
cd .. && go test -v ./tests/unit/...

# 運行整合測試
echo "正在運行整合測試..."
cd .. && go test -v ./tests/integration/...

# 輸出測試覆蓋率報告
echo "生成測試覆蓋率報告..."
cd .. && go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
echo "測試覆蓋率報告已生成: coverage.html"
