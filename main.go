package main

import (
	"context"
	"fmt"
	"os"

	"github.com/ofmeteoriteh/ddns/cloudflare"
	"github.com/ofmeteoriteh/ddns/config"
	"github.com/ofmeteoriteh/ddns/getip"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "setup" {
		runSetup()
		return
	}

	// 解析 --dry-run / Parse --dry-run flag
	dryRun := false
	for _, arg := range os.Args[1:] {
		if arg == "--dry-run" {
			dryRun = true
		}
	}

	if dryRun {
		fmt.Println("── DRY RUN / 模拟运行 ──")
	}

	cfg, err := config.Load("config.json")
	if err != nil {
		fmt.Println("配置文件不存在，请先运行 / Config not found, please run: ddns setup")
		os.Exit(1)
	}

	ctx := context.Background()

	// 获取公网 IP / Fetch public IPs
	ipv4ip, errV4 := getip.GetIPv4IP(ctx)
	ipv6ip, errV6 := getip.GetIPv6IP(ctx)

	if errV4 != nil && errV6 != nil {
		fmt.Println("获取公网 IP 失败 / Failed to get public IP:", errV4, "|", errV6)
		os.Exit(1)
	}

	if errV4 == nil {
		fmt.Println("IPv4:", ipv4ip)
	}
	if errV6 == nil {
		fmt.Println("IPv6:", ipv6ip)
	}

	// 遍历所有条目 / Iterate all entries
	for _, entry := range cfg.Entries {
		for _, t := range entry.Types {
			var ip string
			switch t {
			case "A":
				if errV4 != nil {
					fmt.Printf("  [%s] IPv4 不可用，跳过 / IPv4 unavailable, skipping\n", entry.Name)
					continue
				}
				ip = ipv4ip.String()
			case "AAAA":
				if errV6 != nil {
					fmt.Printf("  [%s] IPv6 不可用，跳过 / IPv6 unavailable, skipping\n", entry.Name)
					continue
				}
				ip = ipv6ip.String()
			}

			result, err := cloudflare.DDNS(ctx, entry.Key, entry.ZoneID, entry.Name, t, ip, entry.Proxied, dryRun)
			if err != nil {
				fmt.Printf("  [%s %s] 失败 / failed: %v\n", entry.Name, t, err)
			} else {
				fmt.Printf("  [%s %s] %s\n", entry.Name, t, result)
			}
		}
	}
}
