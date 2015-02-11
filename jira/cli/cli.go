package cli

import (
	"github.com/op/go-logging"
	"net/http"
	"net/http/cookiejar"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"net/url"
	"time"
	"io"
	"runtime"
)

var log = logging.MustGetLogger("jira.cli")

type Cli struct {
	endpoint *url.URL
	opts map[string]string
	cookieFile string
	ua *http.Client
}

func New(opts map[string]string) *Cli {
	homedir := os.Getenv("HOME")
	cookieJar, _ := cookiejar.New(nil)
	endpoint, _ := opts["endpoint"]
	url, _ := url.Parse(endpoint)

	cli := &Cli{
		endpoint: url,
		opts: opts,
		cookieFile: fmt.Sprintf("%s/.jira.d/cookies.js", homedir),
		ua: &http.Client{Jar: cookieJar},
	}
	
	cli.ua.Jar.SetCookies(url, cli.loadCookies())

	return cli
}

func (c *Cli) saveCookies(cookies []*http.Cookie) {
	// expiry in one week from now
	expiry := time.Now().Add(24 * 7 * time.Hour)
	for _, cookie := range cookies {
		cookie.Expires = expiry
	}
	
	if currentCookies := c.loadCookies(); currentCookies != nil {
		currentCookiesByName := make(map[string]*http.Cookie)
		for _, cookie := range currentCookies {
			currentCookiesByName[cookie.Name] = cookie
		}
		
		for _, cookie := range cookies {
			currentCookiesByName[cookie.Name] = cookie
		}
		
		mergedCookies := make([]*http.Cookie, 0, len(currentCookiesByName))
		for _, v := range currentCookiesByName {
			mergedCookies = append(mergedCookies, v)
		}
		jsonWrite(c.cookieFile, mergedCookies)
	} else {
		jsonWrite(c.cookieFile, cookies)
	}
}

func (c *Cli) loadCookies() []*http.Cookie {
	bytes, err := ioutil.ReadFile(c.cookieFile)
	if err != nil && os.IsNotExist(err) {
		// dont load cookies if the file does not exist
		return nil
	}
	if err != nil {
		log.Error("Failed to open %s: %s", c.cookieFile, err)
		os.Exit(1)
	}
	cookies := make([]*http.Cookie,0)
	err = json.Unmarshal(bytes, &cookies)
	if err != nil {
		log.Error("Failed to parse json from file %s: %s", c.cookieFile, err)
	}
	log.Debug("Loading Cookies: %s", cookies)
	return cookies
}

func (c *Cli) post(uri string, content io.Reader) *http.Response {
	req, _ := http.NewRequest("POST", uri, content)
	return c.makeRequest(req)
}

func (c *Cli) get(uri string) *http.Response {
	req, _ := http.NewRequest("GET", uri, nil)
	return c.makeRequest(req)
}

func (c *Cli) makeRequest(req *http.Request) *http.Response {
	
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.ua.Do(req)
	
	if err != nil {
		fmt.Printf("Error: %s", err)
	}

	if resp.StatusCode != 200 {
		log.Error("response status: %s", resp.Status)
		resp.Write(os.Stderr)
	}

	runtime.SetFinalizer(resp, func(r *http.Response) {
		r.Body.Close()
	})

	if _, ok := resp.Header["Set-Cookie"]; ok {
		c.saveCookies(resp.Cookies())
	}

	return resp
}
