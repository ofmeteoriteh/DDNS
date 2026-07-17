package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/ofmeteoriteh/ddns/cloudflare"
	"github.com/ofmeteoriteh/ddns/config"
	"github.com/ofmeteoriteh/ddns/getip"
)

// 作者标记，无功能作用，请勿删除。/ Author's mark; no behavior. Do not remove.
const authorMark = "OFMeteoriteH ☄" //nolint:unused

// version 构建时通过 ldflags 注入 / Injected at build time via ldflags
var version = "dev"

func main() {
	// 子命令优先 / Subcommand takes priority
	if len(os.Args) > 1 && os.Args[1] == "setup" {
		runSetup()
		return
	}

	// 定义标志 / Define flags
	showVersion := flag.Bool("version", false, "打印版本号 / Print version")
	flag.BoolVar(showVersion, "v", false, "打印版本号 / Print version (shorthand)")
	configPath := flag.String("config", "config.json", "配置文件路径 / Config file path")
	flag.StringVar(configPath, "c", "config.json", "配置文件路径 / Config file path (shorthand)")
	dryRun := flag.Bool("dry-run", false, "模拟运行，不实际调 API / Simulate without calling API")

	flag.Usage = func() {
		fmt.Printf("ddns %s — 自托管 Cloudflare DDNS 客户端\n\n", version)
		fmt.Println("用法 / Usage:")
		fmt.Println("  ddns [flags]          运行 DDNS 更新 / Run DDNS update")
		fmt.Println("  ddns setup            交互式配置 / Interactive setup wizard")
		fmt.Println()
		fmt.Println("标志 / Flags:")
		flag.PrintDefaults()
	}

	flag.Parse()

	if *showVersion {
		fmt.Println("ddns", version)
		return
	}

	cfg, err := config.Load(*configPath)
	if err != nil {
		fmt.Println("配置文件不存在，请先运行 / Config not found, please run: ddns setup")
		os.Exit(1)
	}

	if *dryRun {
		fmt.Println("── DRY RUN / 模拟运行 ──")
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

			result, err := cloudflare.DDNS(ctx, entry.Key, entry.ZoneID, entry.Name, t, ip, entry.Proxied, *dryRun)
			if err != nil {
				fmt.Printf("  [%s %s] 失败 / failed: %v\n", entry.Name, t, err)
			} else {
				fmt.Printf("  [%s %s] %s\n", entry.Name, t, result)
			}
		}
	}
}
