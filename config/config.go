// Package config 提供 DDNS 配置的加载与保存 / Provides config loading and saving for DDNS
package config

import (
	"encoding/json"
	"os"
)

// Config 完整配置结构 / Full configuration structure
type Config struct {
	Keys    []APIKey `json:"keys"`
	Zones   []Zone   `json:"zones"`
	Entries []Entry  `json:"entries"`
}

// APIKey Cloudflare API 密钥 / Cloudflare API key
type APIKey struct {
	Name string `json:"name"`
	Key  string `json:"key"`
}

// Zone Cloudflare DNS Zone / Cloudflare DNS zone
type Zone struct {
	Name   string `json:"name"`
	Domain string `json:"domain"`
	ZoneID string `json:"zone_id"`
}

// Entry DNS 记录条目 / DNS record entry
type Entry struct {
	Name    string   `json:"name"`
	ZoneID  string   `json:"zone_id"`
	Key     string   `json:"key"`
	Types   []string `json:"types"`
	Proxied bool     `json:"proxied"`
}

// Load 从 config.json 加载配置 / Loads config from config.json
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// Save 保存配置到 config.json / Saves config to config.json
func Save(path string, cfg *Config) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}
