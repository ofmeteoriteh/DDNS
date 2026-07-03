package cloudflare

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// DDNSRecord 列出 DNS 记录的响应结构 / Response structure for listing DNS records
type DDNSRecord struct {
	Success    bool `json:"success"`
	Result     []Record
	ResultInfo ResultInfo
}

// Record DNS 记录核心字段 / Core fields of a DNS record
type Record struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Type    string `json:"type"`
	Content string `json:"content"`
	Proxied bool   `json:"proxied"`
	TTL     int    `json:"ttl"`
}

// ResultInfo 分页信息 / Pagination info
type ResultInfo struct {
	Count      int `json:"count"`
	Page       int `json:"page"`
	PerPage    int `json:"per_page"`
	TotalCount int `json:"total_count"`
	TotalPages int `json:"total_pages"`
}

// getDNSRecords 列出指定 Zone 下的 DNS 记录 / Lists DNS records under the specified zone
func getDNSRecords(ctx context.Context, key string, url string) (DDNSRecord, error) {
	client := &http.Client{Timeout: 3 * time.Second}

	// 构建 GET 请求 / Build GET request
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return DDNSRecord{}, err
	}

	// 设置认证头 / Set authorization header
	req.Header.Set("Authorization", "Bearer "+key)
	resp, err := client.Do(req)
	if err != nil {
		return DDNSRecord{}, err
	}
	defer resp.Body.Close()

	// 反序列化响应（result 为数组） / Deserialize response (result is an array)
	var dnsRecords DDNSRecord
	err = json.NewDecoder(resp.Body).Decode(&dnsRecords)
	if err != nil {
		return DDNSRecord{}, err
	}

	return dnsRecords, nil
}

// postDNSRecord 创建一条新的 DNS 记录 / Creates a new DNS record
func postDNSRecord(ctx context.Context, key string, url string, record Record) error {
	client := &http.Client{Timeout: 3 * time.Second}

	// 序列化请求体 / Serialize request body
	body, err := json.Marshal(record)
	if err != nil {
		return err
	}

	// 构建 POST 请求 / Build POST request
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(body))
	if err != nil {
		return err
	}

	// 设置认证头和内容类型 / Set authorization and content-type headers
	req.Header.Set("Authorization", "Bearer "+key)
	req.Header.Set("Content-Type", "application/json")
	_, err = client.Do(req)
	if err != nil {
		return err
	}

	return nil
}

// putDNSRecord 覆盖一条已有的 DNS 记录 / Overwrites an existing DNS record
func putDNSRecord(ctx context.Context, key string, url string, record Record) error {
	client := &http.Client{Timeout: 3 * time.Second}

	// 序列化请求体 / Serialize request body
	body, err := json.Marshal(record)
	if err != nil {
		return err
	}

	// 构建 PUT 请求 / Build PUT request
	req, err := http.NewRequestWithContext(ctx, "PUT", url, bytes.NewReader(body))
	if err != nil {
		return err
	}

	// 设置认证头和内容类型 / Set authorization and content-type headers
	req.Header.Set("Authorization", "Bearer "+key)
	req.Header.Set("Content-Type", "application/json")
	_, err = client.Do(req)
	if err != nil {
		return err
	}

	return nil
}

// DDNS 编排完整流程：列出 → 不存在则创建，存在且 IP 变化则覆盖 / Orchestrates the full flow: list → create if absent, overwrite if IP changed
func DDNS(ctx context.Context, key string, zone_id string, name string, recordType string, ip string) error {
	// 拼接带过滤条件的 URL / Build URL with query filters
	url1 := "https://api.cloudflare.com/client/v4/zones/" + zone_id + "/dns_records?name=" + name + "&type=" + recordType

	// 列出现有记录 / List existing records
	dns, err := getDNSRecords(ctx, key, url1)
	if err != nil {
		return fmt.Errorf("get DNS records: %w", err)
	}

	if len(dns.Result) == 0 {
		// 无记录，创建新记录 / No record found, create new one
		record := Record{
			Name:    name,
			Type:    recordType,
			Content: ip,
			TTL:     1,
		}
		err = postDNSRecord(ctx, key, url1, record)
	} else if dns.Result[0].Content != ip {
		// IP 变化，覆盖记录 / IP changed, overwrite record
		record := dns.Result[0]
		record.Content = ip
		url2 := "https://api.cloudflare.com/client/v4/zones/" + zone_id + "/dns_records/" + record.ID
		err = putDNSRecord(ctx, key, url2, record)
	}
	// IP 未变化，跳过 / IP unchanged, skip
	if err != nil {
		return fmt.Errorf("update DNS record: %w", err)
	}

	return nil
}
