package socks5

import (
	"fmt"
	"testing"
)

func TestParseAddress(t *testing.T) {
	t.Log(ParseAddress("127.0.0.1:80"))
	t.Log(ParseAddress("[::1]:80"))
	t.Log(ParseAddress("a.com:80"))
}

func TestGetRandomIPFromCIDR(t *testing.T) {
	// IPv4示例
	ipv4, err := GetRandomIPFromCidrs([]string{"192.168.1.0/24"})
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Random IPv4: %s\n", ipv4)

	// IPv6示例
	ipv6, err := GetRandomIPFromCidrs([]string{"2001:db8::/64"})
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Random IPv6: %s\n", ipv6)
}
