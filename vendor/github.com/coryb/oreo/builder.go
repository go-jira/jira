package oreo

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type ReqBuilder struct {
	request *http.Request
}

func RequestBuilder(u *url.URL) *ReqBuilder {
	return &ReqBuilder{
		request: &http.Request{
			Method:     "GET",
			URL:        u,
			Proto:      "HTTP/1.1",
			ProtoMajor: 1,
			ProtoMinor: 1,
			Header:     make(http.Header),
			Body:       nil,
			Host:       u.Host,
		},
	}
}

func (b *ReqBuilder) WithHeader(name, value string) *ReqBuilder {
	b.request.Header.Add(name, value)
	return b
}

func (b *ReqBuilder) WithContentType(value string) *ReqBuilder {
	b.request.Header.Add("Content-Type", value)
	return b
}

func (b *ReqBuilder) WithUserAgent(value string) *ReqBuilder {
	b.request.Header.Add("User-Agent", value)
	return b
}

func (b *ReqBuilder) WithMethod(method string) *ReqBuilder {
	b.request.Method = method
	return b
}

func (b *ReqBuilder) WithJSON(data string) *ReqBuilder {
	contentType := "application/json"
	return b.WithContentType(contentType).WithHeader("Accept", contentType).WithBody(strings.NewReader(data))
}

func (b *ReqBuilder) WithXML(data string) *ReqBuilder {
	contentType := "application/xml"
	return b.WithContentType(contentType).WithHeader("Accept", contentType).WithBody(strings.NewReader(data))
}

func (b *ReqBuilder) WithPostForm(data url.Values) *ReqBuilder {
	return b.WithContentType("application/x-www-form-urlencoded").WithBody(strings.NewReader(data.Encode()))
}

func (b *ReqBuilder) WithBody(body io.Reader) *ReqBuilder {
	rc, ok := body.(io.ReadCloser)
	if !ok && body != nil {
		rc = ioutil.NopCloser(body)
	}
	b.request.Body = rc
	return b
}

func (b *ReqBuilder) WithCompression() *ReqBuilder {
	if b.request.Body == nil {
		panic(fmt.Errorf("oreo usage error: WithCompression called before WithBody"))
	}
	buf := bytes.NewBufferString("")
	w := gzip.NewWriter(buf)
	_, err := io.Copy(w, b.request.Body)
	if err != nil {
		panic(err)
	}
	w.Close()
	b.request.Body.Close()
	b.request.Body = ioutil.NopCloser(buf)
	b.request.Header.Add("Content-Encoding", "gzip")
	return b
}

func (b *ReqBuilder) WithAuth(username, password string) *ReqBuilder {
	b.request.SetBasicAuth(username, password)
	return b
}

func (b *ReqBuilder) Build() *http.Request {
	return b.request
}
