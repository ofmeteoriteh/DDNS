// Package cloudflare 对接 Cloudflare API，提供 Token 验证与 DNS 记录管理 / Interacts with Cloudflare API for token verification and DNS record management
package cloudflare

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

// cloudflare 通用响应外层 / Generic Cloudflare API response wrapper
type cloudflare struct {
	Success bool `json:"success"`
	Result  result
}

// result 验证端点的单个结果 / Single result object for the verify endpoint
type result struct {
	Status string `json:"status"`
}

// VerifyAPI 验证 API Token 是否有效且处于 active 状态 / Verifies that the API token is valid and active
func VerifyAPI(ctx context.Context, key string) (bool, error) {
	client := &http.Client{Timeout: 3 * time.Second}

	// 构建验证请求 / Build verification request
	req, err := http.NewRequestWithContext(ctx, "GET", "https://api.cloudflare.com/client/v4/user/tokens/verify", nil)
	if err != nil {
		return false, err
	}

	// 设置认证头 / Set authorization header
	req.Header.Set("Authorization", "Bearer "+key)
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	// 反序列化响应 / Deserialize response
	var cf cloudflare
	if err := json.NewDecoder(resp.Body).Decode(&cf); err != nil {
		return false, err
	}

	// 检查 Token 是否有效且处于 active 状态 / Check token is valid and active
	if cf.Success && cf.Result.Status == "active" {
		return true, nil
	}

	return false, nil
}
