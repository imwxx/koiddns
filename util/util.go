package util

import (
	"errors"
	"io"
	"net"
	"net/http"
	"strings"
	"time"
)

// DefaultIPFetchTimeout 获取公网 IP 的 HTTP 超时时间
const DefaultIPFetchTimeout = 10 * time.Second

// GetPublicIP 从默认源获取公网 IP，带超时与校验。失败时返回 error。
func GetPublicIP() (string, error) {
	return GetPublicIPFrom("http://ipinfo.io/ip", DefaultIPFetchTimeout)
}

// GetPublicIPFrom 从指定 URL 获取公网 IP，超时后校验是否为合法 IP。
func GetPublicIPFrom(url string, timeout time.Duration) (string, error) {
	client := &http.Client{Timeout: timeout}
	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	ip := strings.TrimSpace(string(body))
	if ip == "" {
		return "", errors.New("empty response from IP service")
	}
	if net.ParseIP(ip) == nil {
		return "", errors.New("invalid IP: " + ip)
	}
	return ip, nil
}
