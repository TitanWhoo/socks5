package socks5_test

import (
	"context"
	"encoding/hex"
	"github.com/TitanWhoo/socks5"
	"github.com/miekg/dns"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"testing"
)

func ExampleServer() {
	s, err := socks5.NewServer("127.0.0.1:1080", "127.0.0.1", "titan", "", []string{"198.18.0.1/32"}, 60, 60)
	if err != nil {
		log.Println(err)
		return
	}
	// You can pass in custom Handler
	s.ListenAndServe(nil)
	// #Output:
}

func TestExampleServer_start(t *testing.T) {
	ExampleServer()
}

func ExampleClient_tcp() {
	go ExampleServer()
	c, err := socks5.NewClient("127.0.0.1:1080", "", "", 0, 60)
	if err != nil {
		log.Println(err)
		return
	}
	client := &http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				return c.Dial(network, addr)
			},
		},
	}
	res, err := client.Get("https://httpbin.org/ip")
	if err != nil {
		log.Println(err)
		return
	}
	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("tcp", string(b))
	// Output:
}

func ExampleClient_udp() {
	go ExampleServer()
	c, err := socks5.NewClient("127.0.0.1:1080", "", "", 0, 60)
	if err != nil {
		log.Println(err)
		return
	}
	conn, err := c.Dial("udp", "8.8.8.8:53")
	if err != nil {
		log.Println(err)
		return
	}
	b, err := hex.DecodeString("0001010000010000000000000a74787468696e6b696e6703636f6d0000010001")
	if err != nil {
		log.Println(err)
		return
	}
	if _, err := conn.Write(b); err != nil {
		log.Println(err)
		return
	}
	b = make([]byte, 2048)
	n, err := conn.Read(b)
	if err != nil {
		log.Println(err)
		return
	}
	m := &dns.Msg{}
	if err := m.Unpack(b[0:n]); err != nil {
		log.Println(err)
		return
	}
	log.Println(m.String())
	// Output:
}
