package unit

import (
	"net/http/httptest"
	"testing"

	"github.com/lipeichen/ticket-getter/pkg/utils"
)

func TestExtractTLSFingerprint(t *testing.T) {
	// 創建一個 HTTP 請求
	req := httptest.NewRequest("GET", "http://example.com", nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	// 提取指紋
	fingerprint := utils.ExtractTLSFingerprint(req)

	// 驗證指紋非空
	if fingerprint == "" {
		t.Error("Expected non-empty fingerprint, got empty string")
	}

	// 使用相同的 HTTP 請求提取指紋，應該得到相同結果
	fingerprint2 := utils.ExtractTLSFingerprint(req)
	if fingerprint != fingerprint2 {
		t.Errorf("Expected consistent fingerprints, got different values: %s and %s", fingerprint, fingerprint2)
	}

	// 修改 HTTP 請求頭，應該得到不同的指紋
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	fingerprint3 := utils.ExtractTLSFingerprint(req)
	if fingerprint == fingerprint3 {
		t.Error("Expected different fingerprints for different headers, got the same value")
	}
}

func TestGetClientIPFingerprint(t *testing.T) {
	// 創建一個 HTTP 請求
	req := httptest.NewRequest("GET", "http://example.com", nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36")
	req.Header.Set("X-Forwarded-For", "192.168.1.1")
	
	// 獲取客戶端 IP 和指紋組合
	ipFingerprint := utils.GetClientIPFingerprint(req)
	
	// 驗證 IP 指紋非空
	if ipFingerprint == "" {
		t.Error("Expected non-empty IP fingerprint, got empty string")
	}
	
	// 使用相同的 HTTP 請求獲取指紋，應該得到相同結果
	ipFingerprint2 := utils.GetClientIPFingerprint(req)
	if ipFingerprint != ipFingerprint2 {
		t.Errorf("Expected consistent IP fingerprints, got different values: %s and %s", ipFingerprint, ipFingerprint2)
	}
	
	// 修改 IP，應該得到不同的指紋
	req.Header.Set("X-Forwarded-For", "192.168.1.2")
	ipFingerprint3 := utils.GetClientIPFingerprint(req)
	if ipFingerprint == ipFingerprint3 {
		t.Error("Expected different IP fingerprints for different IPs, got the same value")
	}
}
