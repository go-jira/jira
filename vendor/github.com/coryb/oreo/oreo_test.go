package oreo

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func init() {
	TraceRequestBody = true
	TraceResponseBody = true
}

func TestOreoGet(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		fmt.Fprintf(w, "OK")
	}))
	defer ts.Close()

	c := New()
	resp, err := c.Get(ts.URL)
	assert.Nil(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)
	body, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, []byte("OK"), body)
}

func TestOreoHead(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "HEAD", r.Method)
		fmt.Fprintf(w, "OK")
	}))
	defer ts.Close()

	c := New()
	resp, err := c.Head(ts.URL)
	assert.Nil(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)
	body, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, []byte(""), body)
	assert.Equal(t, int64(2), resp.ContentLength)
}

func TestOreoPost(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		body, _ := ioutil.ReadAll(r.Body)
		assert.Equal(t, []byte("DATA"), body)
		contentLength := r.Header["Content-Type"][0]
		assert.Equal(t, "text/plain", contentLength)

		fmt.Fprintf(w, "OK")
	}))
	defer ts.Close()

	c := New()
	resp, err := c.Post(ts.URL, "text/plain", strings.NewReader("DATA"))
	assert.Nil(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)
	body, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, []byte("OK"), body)
	assert.Equal(t, int64(2), resp.ContentLength)
}

func TestOreoPostForm(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		body, _ := ioutil.ReadAll(r.Body)
		assert.Equal(t, []byte("key=value"), body)
		contentLength := r.Header["Content-Type"][0]
		assert.Equal(t, "application/x-www-form-urlencoded", contentLength)

		fmt.Fprintf(w, "OK")
	}))
	defer ts.Close()

	c := New()
	data := url.Values{}
	data.Add("key", "value")
	resp, err := c.PostForm(ts.URL, data)
	assert.Nil(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)
	body, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, []byte("OK"), body)
	assert.Equal(t, int64(2), resp.ContentLength)
}

func TestOreoPostJSON(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		body, _ := ioutil.ReadAll(r.Body)
		assert.Equal(t, []byte(`{"key":"value"}`), body)
		contentLength := r.Header["Content-Type"][0]
		assert.Equal(t, "application/json", contentLength)

		fmt.Fprintf(w, "OK")
	}))
	defer ts.Close()

	c := New()
	resp, err := c.PostJSON(ts.URL, `{"key":"value"}`)
	assert.Nil(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)
	body, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, []byte("OK"), body)
	assert.Equal(t, int64(2), resp.ContentLength)
}

func TestOreoPut(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PUT", r.Method)
		body, _ := ioutil.ReadAll(r.Body)
		assert.Equal(t, []byte("DATA"), body)
		contentLength := r.Header["Content-Type"][0]
		assert.Equal(t, "text/plain", contentLength)

		fmt.Fprintf(w, "OK")
	}))
	defer ts.Close()

	c := New()
	resp, err := c.Put(ts.URL, "text/plain", strings.NewReader("DATA"))
	assert.Nil(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)
	body, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, []byte("OK"), body)
	assert.Equal(t, int64(2), resp.ContentLength)
}

func TestOreoPutJSON(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PUT", r.Method)
		body, _ := ioutil.ReadAll(r.Body)
		assert.Equal(t, []byte(`{"key":"value"}`), body)
		contentLength := r.Header["Content-Type"][0]
		assert.Equal(t, "application/json", contentLength)

		fmt.Fprintf(w, "OK")
	}))
	defer ts.Close()

	c := New()
	resp, err := c.PutJSON(ts.URL, `{"key":"value"}`)
	assert.Nil(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)
	body, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, []byte("OK"), body)
	assert.Equal(t, int64(2), resp.ContentLength)
}

func TestOreoDelete(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "DELETE", r.Method)
		fmt.Fprintf(w, "OK")
	}))
	defer ts.Close()

	c := New()
	resp, err := c.Delete(ts.URL)
	assert.Nil(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)
	body, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, []byte("OK"), body)
}

func TestOreoWithRetries(t *testing.T) {
	attempts := 0
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		http.Error(w, "error", http.StatusInternalServerError)
	}))
	defer ts.Close()

	c := New().WithRetries(2)
	resp, err := c.Get(ts.URL)
	assert.Nil(t, err)
	assert.Equal(t, 3, attempts)
	assert.NotNil(t, resp)
	assert.Equal(t, 500, resp.StatusCode)
}

func TestOreoWithTimeout(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
		http.Error(w, "error", http.StatusInternalServerError)
	}))
	defer ts.Close()

	c := New().WithTimeout(1 * time.Second).WithRetries(2)
	start := time.Now().Unix()
	resp, err := c.Get(ts.URL)
	end := time.Now().Unix()
	assert.Nil(t, resp)
	assert.Error(t, err)
	assert.True(t, end-start >= 2, "duration more than 2x timeout")
}

func TestOreoWithLinearTimeout(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
		http.Error(w, "error", http.StatusInternalServerError)
	}))
	defer ts.Close()

	c := New().WithTimeout(1 * time.Second).WithBackoff(LINEAR_BACKOFF).WithRetries(2)

	start := time.Now().Unix()
	resp, err := c.Get(ts.URL)
	end := time.Now().Unix()
	assert.Nil(t, resp)
	assert.Error(t, err)
	assert.True(t, end-start >= 3, "duration more than 1*timeout + 2*timeout")
}

func TestOreoWithCookieFile(t *testing.T) {
	request := 0
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		request++
		switch request {
		case 1:
			cookie := &http.Cookie{
				Name:  "key1",
				Value: "val1",
			}
			http.SetCookie(w, cookie)
		case 2:
			cookie := r.Header["Cookie"][0]
			assert.Equal(t, "key1=val1", cookie)
		}
		fmt.Fprintf(w, "OK")
	}))
	defer ts.Close()

	tmpFile, err := ioutil.TempFile("", "oreo-cookies")
	assert.Nil(t, err)
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()
	os.Remove(tmpFile.Name())

	c := New().WithCookieFile(tmpFile.Name())
	// first request will get a cookie set on response
	resp, err := c.Get(ts.URL)
	assert.NotNil(t, resp)
	assert.Nil(t, err)

	/// this request should automatically send cookie back to server
	resp, err = c.Get(ts.URL)
	assert.NotNil(t, resp)
	assert.Nil(t, err)
}

func TestOreoWithTransport(t *testing.T) {
	// set tcp connect timeout to 5s
	var netTransport = &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 5 * time.Second,
		}).Dial,
	}

	c := New().WithTransport(netTransport).WithRetries(0)

	// test against google dns servers, we will get a tcp connection
	//  failure (timeout due to firewall) to a non dns port on those hosts
	start := time.Now().Unix()
	resp, err := c.Get("http://8.8.8.8:9999")
	end := time.Now().Unix()
	assert.Nil(t, resp)
	assert.Error(t, err)
	lapse := end - start
	msg := fmt.Sprintf("duration between 5-6s timeout, got: %d", lapse)
	assert.True(t, lapse >= 5 && lapse <= 6, msg)

}

func TestOreoWithPostCallback(t *testing.T) {
	requests := 0
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		_, ok := r.Header["Authorization"]
		if ok {
			fmt.Fprintf(w, "OK")
		} else {
			http.Error(w, "error", http.StatusUnauthorized)
		}
	}))
	defer ts.Close()

	var c *Client
	called := 0
	callback := func(req *http.Request, resp *http.Response) (*http.Response, error) {
		called++
		// if we get a 401 then add auth headers and try the request again
		if resp.StatusCode == 401 {
			req.SetBasicAuth("user", "pass")
			return c.Do(req)
		}
		return resp, nil
	}

	c = New().WithPostCallback(callback)

	resp, err := c.Get(ts.URL)
	assert.NotNil(t, resp)
	assert.Nil(t, err)
	assert.Equal(t, 1, called)
	assert.Equal(t, 2, requests)
}

func TestOreoWithPreCallback(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, ok := r.Header["Authorization"]
		assert.True(t, ok)
		fmt.Fprintf(w, "OK")
	}))
	defer ts.Close()

	callback := func(req *http.Request) (*http.Request, error) {
		req.SetBasicAuth("user", "pass")
		return req, nil
	}

	c := New().WithPreCallback(callback)

	resp, err := c.Get(ts.URL)
	assert.NotNil(t, resp)
	assert.Nil(t, err)
}

func TestOreoWithRedirect(t *testing.T) {
	requests := 0
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		if requests == 1 {
			http.Redirect(w, r, "/redirect", http.StatusMovedPermanently)
			return
		}
		fmt.Fprintf(w, "OK")
	}))
	defer ts.Close()

	c := New()

	resp, err := c.Get(ts.URL)
	assert.NotNil(t, resp)
	assert.Nil(t, err)
	assert.Equal(t, 2, requests)
}

func TestOreoWithNoRedirect(t *testing.T) {
	requests := 0
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		if requests == 1 {
			http.Redirect(w, r, "/redirect/", http.StatusMovedPermanently)
		} else {
			fmt.Fprintf(w, "OK")
		}
	}))
	defer ts.Close()

	c := New().WithCheckRedirect(NoRedirect)

	resp, err := c.Get(ts.URL)
	assert.NotNil(t, resp)
	assert.Nil(t, err)
	assert.Equal(t, 1, requests)
}

func TestOreoWithImmutability(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "OK")
	}))
	defer ts.Close()

	result := ""
	callback1 := func(req *http.Request) (*http.Request, error) {
		result = "callback1"
		return req, nil
	}

	callback2 := func(req *http.Request) (*http.Request, error) {
		result = "callback2"
		return req, nil
	}

	c1 := New().WithPreCallback(callback1)
	c2 := c1.WithPreCallback(callback2)

	resp, err := c1.Get(ts.URL)
	assert.NotNil(t, resp)
	assert.Nil(t, err)
	assert.Equal(t, "callback1", result)

	resp, err = c2.Get(ts.URL)
	assert.NotNil(t, resp)
	assert.Nil(t, err)
	assert.Equal(t, "callback2", result)

	resp, err = c1.Get(ts.URL)
	assert.NotNil(t, resp)
	assert.Nil(t, err)
	assert.Equal(t, "callback1", result)
}

func TestOreoPostCompressed(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "gzip", r.Header.Get("Content-Encoding"))

		reader, err := gzip.NewReader(r.Body)
		assert.Nil(t, err)
		defer reader.Close()
		buf := bytes.NewBufferString("")
		_, err = io.Copy(buf, reader)
		assert.Nil(t, err)

		assert.Equal(t, []byte("DATA"), buf.Bytes())
		contentLength := r.Header["Content-Type"][0]
		assert.Equal(t, "text/plain", contentLength)

		fmt.Fprintf(w, "OK")
	}))
	defer ts.Close()

	c := New()
	parsed, _ := url.Parse(ts.URL)
	req := RequestBuilder(parsed).WithMethod("POST").WithContentType("text/plain").WithBody(strings.NewReader("DATA")).WithCompression().Build()
	resp, err := c.Do(req)
	assert.Nil(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)
	body, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, []byte("OK"), body)
	assert.Equal(t, int64(2), resp.ContentLength)
}
