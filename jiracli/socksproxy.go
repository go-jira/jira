package jiracli

import (
	"net"
	"net/http"
	"time"

	"golang.org/x/net/proxy"
)

func socksProxy(address string) *http.Transport {
	return newSocksProxyTransport(address)
}

func newSocksProxyTransport(address string) *http.Transport {
	dialer, err := proxy.SOCKS5("tcp", address, nil, proxy.Direct)
	if err != nil {
		// TODO: whoops, return error?
		panic(err)
	}
	dial := func(network, addr string) (net.Conn, error) {
		return dialer.Dial(network, addr)
	}

	return &http.Transport{
		Dial:                  dial,
		DisableKeepAlives:     true,
		ResponseHeaderTimeout: 30 * time.Second,
		ExpectContinueTimeout: 10 * time.Second,
	}
}
