package main

import (
	"context"
	"fmt"

	"github.com/ofmeteoriteh/ddns/getip"
)

func main() {
	ctx := context.Background()

	ipv4ip, err_v4 := getip.GetIPv4IP(ctx)
	ipv6ip, err_v6 := getip.GetIPv6IP(ctx)

	fmt.Println(ipv4ip, ipv6ip)
	fmt.Errorf(err_v4.Error(), err_v6.Error())
}
