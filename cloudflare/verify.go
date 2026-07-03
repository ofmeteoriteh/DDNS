package cloudflare

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

type cloudflare struct {
	Success bool
	Result  result
}

type result struct {
	Status string
}

func VerifyAPI(ctx context.Context, key string) (bool, error) {
	client := &http.Client{Timeout: 3 * time.Second}

	req, err := http.NewRequestWithContext(ctx, "GET", "https://api.cloudflare.com/client/v4/user/tokens/verify", nil)
	if err != nil {
		return false, err
	}

	req.Header.Set("Authorization", "Bearer "+key)
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	var cf cloudflare

	if err := json.NewDecoder(resp.Body).Decode(&cf); err != nil {
		return false, err
	}

	if cf.Success && cf.Result.Status == "active" {
		return true, nil
	}

	return false, nil
}
