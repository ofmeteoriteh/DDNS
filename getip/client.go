// Package getip 通过多个公网源获取本机出口 IPv4/IPv6 地址 / Fetches public IPv4/IPv6 addresses from multiple sources
package getip

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/netip"
	"strings"
	"time"
)

// 纯单栈源列表 / Single-stack source lists
var (
	v4Sources = []string{
		"https://ipv4.icanhazip.com",
		"https://checkip.amazonaws.com",
		"https://api-ipv4.ip.sb/ip",
		"https://v4.api.ipinfo.io/ip",
		"https://ipv4.myexternalip.com/raw",
		"https://v4.ident.me",
	}

	v6Sources = []string{
		"https://ipv6.icanhazip.com",
		"https://api-ipv6.ip.sb/ip",
		"https://ipv6.myexternalip.com/raw",
		"https://v6.api.ipinfo.io/ip",
		"https://v6.ident.me",
	}
)

// fetchIP 从多个源获取 IP，并且多重校验 / Fetches IP from multiple sources with multi-layer validation
func fetchIP(ctx context.Context, wantV6 bool) (netip.Addr, error) {
	// ctx 超时设置 / Context timeout setup
	ctx, cancel := context.WithTimeout(ctx, 35*time.Second)
	defer cancel()

	// User-Agent
	const UA = "ofmh-ddns-client/1.0.0"
	var addr netip.Addr

	// 判断需要选择的源 / Select source list based on IP version
	var selectedSources []string
	if wantV6 {
		selectedSources = v6Sources
	} else {
		selectedSources = v4Sources
	}

	// 创建有超时的 HTTP 客户端 / Create HTTP client with timeout
	client := &http.Client{Timeout: 5 * time.Second}

	// 循环获取 IP 地址，并且校验 / Loop through sources to fetch and validate IP
	for _, source := range selectedSources {
		// 构建请求 / Build request
		req, err := http.NewRequestWithContext(ctx, "GET", source, nil)
		if err != nil {
			continue
		}

		// 设置 User-Agent / Set User-Agent header
		req.Header.Set("User-Agent", UA)

		// 发送请求 / Send request
		resp, err := client.Do(req)
		if err != nil {
			continue
		}

		// 校验是否成功 / Check response status
		if resp.StatusCode != http.StatusOK {
			_ = resp.Body.Close()
			continue
		}

		// 读取响应体 / Read response body
		body, err := io.ReadAll(resp.Body)
		_ = resp.Body.Close()
		if err != nil {
			continue
		}

		// 去除首尾空白后解析 IP 地址 / Trim whitespace and parse IP address
		addr, err = netip.ParseAddr(strings.TrimSpace(string(body)))
		if err != nil {
			continue
		}

		// 校验 IP 地址族 / Validate IP address family
		if wantV6 && !addr.Is6() {
			continue
		}

		if !wantV6 && !addr.Is4() {
			continue
		}

		return addr, nil
	}
	return addr, fmt.Errorf("failed to fetch IP address from all sources")
}

// GetIPv4IP 获取公网 IPv4 地址 / Fetches public IPv4 address
func GetIPv4IP(ctx context.Context) (netip.Addr, error) {
	ipv4ip, err := fetchIP(ctx, false)
	return ipv4ip, err
}

// GetIPv6IP 获取公网 IPv6 地址 / Fetches public IPv6 address
func GetIPv6IP(ctx context.Context) (netip.Addr, error) {
	ipv6ip, err := fetchIP(ctx, true)
	return ipv6ip, err
}
