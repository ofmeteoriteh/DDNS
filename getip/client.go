package getip

import (
	"context"
	"io"
	"net/http"
	"net/netip"
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

func fetchIP(ctx context.Context, sources []string, wantV6 bool) (netip.Addr, error) {
	const UA = "ofmh-ddns-client/1.0.0"
	var addr netip.Addr

	if wantV6 == false {
		sources = v4Sources
	} else {
		sources = v6Sources
	}

	for _, source := range sources {

	}
	return addr, err
}
