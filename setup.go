package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/ofmeteoriteh/ddns/config"
)

// runSetup 交互式配置向导 / Interactive setup wizard
func runSetup() {
	cfg := &config.Config{}

	// 检查已有配置 / Check existing config
	if _, err := os.Stat("config.json"); err == nil {
		var overwrite bool
		survey.AskOne(&survey.Confirm{
			Message: "config.json 已存在，是否覆盖？/ config.json exists, overwrite?",
		}, &overwrite)
		if !overwrite {
			return
		}
	}

	// ── API Keys ──
	fmt.Println("\n── API Keys ──")
	for {
		var name, key string
		survey.AskOne(&survey.Input{Message: "API Key 名称 / Name (e.g. main):"}, &name, survey.WithValidator(survey.Required))
		survey.AskOne(&survey.Input{Message: "API Key / Token:"}, &key, survey.WithValidator(survey.Required))
		cfg.Keys = append(cfg.Keys, config.APIKey{Name: name, Key: key})

		var more bool
		survey.AskOne(&survey.Confirm{Message: "添加更多 API Key？/ Add another API Key?"}, &more)
		if !more {
			break
		}
	}

	// ── Zones ──
	fmt.Println("\n── Zones ──")
	for {
		var name, domain, zoneID string
		survey.AskOne(&survey.Input{Message: "Zone 名称 / Name (e.g. ofqfw-top):"}, &name, survey.WithValidator(survey.Required))
		survey.AskOne(&survey.Input{Message: "域名 / Domain (e.g. ofqfw.top):"}, &domain, survey.WithValidator(survey.Required))
		survey.AskOne(&survey.Input{Message: "Zone ID:"}, &zoneID, survey.WithValidator(survey.Required))
		cfg.Zones = append(cfg.Zones, config.Zone{Name: name, Domain: domain, ZoneID: zoneID})

		var more bool
		survey.AskOne(&survey.Confirm{Message: "添加更多 Zone？/ Add another Zone?"}, &more)
		if !more {
			break
		}
	}

	// ── Entries ──
	fmt.Println("\n── DNS Entries ──")
	zoneNames := make([]string, len(cfg.Zones))
	for i, z := range cfg.Zones {
		zoneNames[i] = z.Name
	}
	keyNames := make([]string, len(cfg.Keys))
	for i, k := range cfg.Keys {
		keyNames[i] = k.Name
	}

	for {
		var prefix string
		survey.AskOne(&survey.Input{
			Message: "前缀 / Prefix (e.g. legendvps-singapore-ddns，留空为根域名 / leave empty for root):",
		}, &prefix)

		var zoneChoice string
		survey.AskOne(&survey.Select{
			Message: "域名 / Domain:",
			Options: zoneNames,
		}, &zoneChoice)

		var keyChoice string
		survey.AskOne(&survey.Select{
			Message: "API Key:",
			Options: keyNames,
		}, &keyChoice)

		var types []string
		survey.AskOne(&survey.MultiSelect{
			Message: "记录类型 / Record types:",
			Options: []string{"A (IPv4)", "AAAA (IPv6)"},
		}, &types)

		var proxied bool
		survey.AskOne(&survey.Confirm{
			Message: "启用 Cloudflare 代理（橙色云）？/ Enable Cloudflare proxy (orange cloud)?",
		}, &proxied)

		// 拼接完整域名 / Build full domain name
		zone := findZone(cfg.Zones, zoneChoice)
		name := zone
		if prefix != "" {
			name = prefix + "." + zone
		}

		// 转换类型 / Convert types
		recordTypes := convertTypes(types)

		cfg.Entries = append(cfg.Entries, config.Entry{
			Name:    name,
			ZoneID:  findZoneID(cfg.Zones, zoneChoice),
			Key:     findKey(cfg.Keys, keyChoice),
			Types:   recordTypes,
			Proxied: proxied,
		})

		var more bool
		survey.AskOne(&survey.Confirm{Message: "添加更多条目？/ Add another entry?"}, &more)
		if !more {
			break
		}
	}

	// 保存配置 / Save config
	if err := config.Save("config.json", cfg); err != nil {
		fmt.Println("保存配置失败 / Failed to save config:", err)
		return
	}
	fmt.Println("\n配置已保存到 / Config saved to config.json")

	// systemd service
	var genSystemd bool
	survey.AskOne(&survey.Confirm{
		Message: "生成 systemd service 文件？/ Generate systemd service file?",
	}, &genSystemd)

	if genSystemd {
		var binPath string
		survey.AskOne(&survey.Input{
			Message: "二进制路径 / Binary path:",
			Default: "/opt/ddns/ddns",
		}, &binPath)

		service := fmt.Sprintf(`[Unit]
Description=DDNS Client
After=network-online.target
Wants=network-online.target

[Service]
Type=oneshot
ExecStart=%s
WorkingDirectory=/opt/ddns
Restart=on-failure
RestartSec=5s

[Install]
WantedBy=multi-user.target
`, binPath)

		if err := os.WriteFile("ddns.service", []byte(service), 0644); err != nil {
			fmt.Println("生成 service 文件失败 / Failed to generate service file:", err)
			return
		}
		fmt.Println("已生成 / Generated ddns.service")
		fmt.Println("  cp ddns.service /etc/systemd/system/")
		fmt.Println("  systemctl daemon-reload")
		fmt.Println("  systemctl enable --now ddns")
	}
}

// findZone 根据名称查找 Zone 域名 / Find zone domain by name
func findZone(zones []config.Zone, name string) string {
	for _, z := range zones {
		if z.Name == name {
			return z.Domain
		}
	}
	return ""
}

// findZoneID 根据名称查找 Zone ID / Find zone ID by name
func findZoneID(zones []config.Zone, name string) string {
	for _, z := range zones {
		if z.Name == name {
			return z.ZoneID
		}
	}
	return ""
}

// findKey 根据名称查找 API Key / Find API key by name
func findKey(keys []config.APIKey, name string) string {
	for _, k := range keys {
		if k.Name == name {
			return k.Key
		}
	}
	return ""
}

// convertTypes 将用户选择的类型转为 API 值 / Convert user-selected types to API values
func convertTypes(types []string) []string {
	var result []string
	for _, t := range types {
		if strings.HasPrefix(t, "A ") || t == "A" {
			result = append(result, "A")
		} else if strings.HasPrefix(t, "AAAA") {
			result = append(result, "AAAA")
		}
	}
	return result
}
