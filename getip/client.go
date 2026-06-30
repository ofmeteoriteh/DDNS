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

func fetchIP(ctx context.Context, wantV6 bool) (netip.Addr, error) {
	ctx, cancel := context.WithTimeout(ctx, 35*time.Second)
	defer cancel()

	const UA = "ofmh-ddns-client/1.0.0"
	var addr netip.Addr

	var selectedSources []string
	if wantV6 {
		selectedSources = v6Sources
	} else {
		selectedSources = v4Sources
	}

	client := &http.Client{Timeout: 5 * time.Second}

	for _, source := range selectedSources {
		req, err := http.NewRequestWithContext(ctx, "GET", source, nil)
		if err != nil {
			continue
		}

		req.Header.Set("User-Agent", UA)

		resp, err := client.Do(req)
		if err != nil {
			continue
		}

		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			continue
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			continue
		}

		addr, err = netip.ParseAddr(strings.TrimSpace(string(body)))
		if err != nil {
			continue
		}

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

func GetIPv4IP(ctx context.Context) (netip.Addr, error) {
	ipv4ip, err := fetchIP(ctx, false)
	return ipv4ip, err
}

func GetIPv6IP(ctx context.Context) (netip.Addr, error) {
	ipv6ip, err := fetchIP(ctx, true)
	return ipv6ip, err
}
