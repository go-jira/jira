package jiracli

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"time"
)

type Transport struct {
	shadow *http.Transport
}

func NewUnixProxyTransport(path string) *Transport {
	dial := func(network, addr string) (net.Conn, error) {
		return net.Dial("unix", path)
	}

	shadow := &http.Transport{
		Dial:                  dial,
		DialTLS:               dial,
		DisableKeepAlives:     true,
		ResponseHeaderTimeout: 30 * time.Second,
		ExpectContinueTimeout: 10 * time.Second,
	}

	return &Transport{shadow}
}

func UnixProxy(path string) *Transport {
	return NewUnixProxyTransport(os.ExpandEnv(path))
}

func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	req2 := *req
	url2 := *req.URL
	req2.URL = &url2
	req2.URL.Opaque = fmt.Sprintf("//%s%s", req.URL.Host, req.URL.EscapedPath())
	return t.shadow.RoundTrip(&req2)
}
