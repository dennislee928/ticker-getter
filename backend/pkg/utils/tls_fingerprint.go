package utils

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"
	"sort"
	"strings"
)

// ExtractTLSFingerprint 從 HTTP 請求提取 TLS 指紋
func ExtractTLSFingerprint(r *http.Request) string {
	// 從請求頭獲取相關信息
	headers := map[string]string{
		"User-Agent":       r.Header.Get("User-Agent"),
		"Accept":           r.Header.Get("Accept"),
		"Accept-Language":  r.Header.Get("Accept-Language"),
		"Accept-Encoding":  r.Header.Get("Accept-Encoding"),
		"Sec-Ch-Ua":        r.Header.Get("Sec-Ch-Ua"),
		"Sec-Ch-Ua-Mobile": r.Header.Get("Sec-Ch-Ua-Mobile"),
		"Sec-Ch-Ua-Platform": r.Header.Get("Sec-Ch-Ua-Platform"),
		"Sec-Fetch-Dest":     r.Header.Get("Sec-Fetch-Dest"),
		"Sec-Fetch-Mode":     r.Header.Get("Sec-Fetch-Mode"),
		"Sec-Fetch-Site":     r.Header.Get("Sec-Fetch-Site"),
	}

	// 從請求獲取 IP 地址
	ip := r.RemoteAddr
	forwardedFor := r.Header.Get("X-Forwarded-For")
	if forwardedFor != "" {
		// 使用第一個代理前的 IP
		ips := strings.Split(forwardedFor, ",")
		ip = strings.TrimSpace(ips[0])
	}
	
	// 可能的情況下從 TLS 連接獲取更多信息
	var tlsInfo string
	if r.TLS != nil {
		tlsVersion := r.TLS.Version
		cipherSuite := r.TLS.CipherSuite
		tlsInfo = fmt.Sprintf("TLSv%d-CS:%d", tlsVersion, cipherSuite)
	}

	// 創建有序的字符串列表
	var parts []string
	for k, v := range headers {
		if v != "" {
			parts = append(parts, fmt.Sprintf("%s=%s", k, v))
		}
	}
	
	// 添加 IP 和 TLS 信息
	parts = append(parts, fmt.Sprintf("IP=%s", ip))
	if tlsInfo != "" {
		parts = append(parts, tlsInfo)
	}
	
	// 排序以確保結果一致
	sort.Strings(parts)
	
	// 計算 SHA-256 哈希
	fingerprint := sha256.Sum256([]byte(strings.Join(parts, "|")))
	
	// 返回 base64 編碼的指紋
	return base64.StdEncoding.EncodeToString(fingerprint[:])
}

// GetClientIPFingerprint 獲取客戶端 IP 和 TLS 指紋的組合
func GetClientIPFingerprint(r *http.Request) string {
	// 獲取客戶端 IP
	ip := r.RemoteAddr
	forwardedFor := r.Header.Get("X-Forwarded-For")
	if forwardedFor != "" {
		ips := strings.Split(forwardedFor, ",")
		ip = strings.TrimSpace(ips[0])
	}
	
	// 從客戶端獲取指紋
	clientFingerprint := r.Header.Get("X-TLS-Fingerprint")
	
	// 如果客戶端沒有提供指紋，生成一個
	if clientFingerprint == "" {
		clientFingerprint = ExtractTLSFingerprint(r)
	}
	
	// 組合 IP 和指紋
	combinedFingerprint := fmt.Sprintf("%s|%s", ip, clientFingerprint)
	
	// 計算 SHA-256 哈希
	hash := sha256.Sum256([]byte(combinedFingerprint))
	
	// 返回 base64 編碼的結果
	return base64.StdEncoding.EncodeToString(hash[:])
}
