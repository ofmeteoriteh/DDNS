package main

import (
	"context"
	"fmt"
	"os"

	"github.com/ofmeteoriteh/ddns/cloudflare"
	"github.com/ofmeteoriteh/ddns/getip"
)

func main() {
	ctx := context.Background()

	// 读取环境变量 / Read environment variables
	name := os.Getenv("DDNS_NAME")
	key := os.Getenv("CLOUDFLARE_API_KEY")
	zoneID := os.Getenv("CLOUDFLARE_ZONE_ID")

	// 获取公网 IP / Fetch public IPs
	ipv4ip, err_v4 := getip.GetIPv4IP(ctx)
	ipv6ip, err_v6 := getip.GetIPv6IP(ctx)

	if err_v4 != nil && err_v6 != nil {
		fmt.Println("获取公网 IP 失败 / Failed to get public IP:", err_v4, "|", err_v6)
		os.Exit(1)
	}

	if err_v4 == nil {
		// IPv4 成功，更新 A 记录 / IPv4 succeeded, update A record
		if err := cloudflare.DDNS(ctx, key, zoneID, name, "A", ipv4ip.String()); err != nil {
			fmt.Println("更新 A 记录失败 / Failed to update A record:", err)
		}
	}

	if err_v6 == nil {
		// IPv6 成功，更新 AAAA 记录 / IPv6 succeeded, update AAAA record
		if err := cloudflare.DDNS(ctx, key, zoneID, name, "AAAA", ipv6ip.String()); err != nil {
			fmt.Println("更新 AAAA 记录失败 / Failed to update AAAA record:", err)
		}
	}
}
