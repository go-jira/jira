package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/op/go-logging"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

var log = logging.MustGetLogger("jira.cli")

type Cli struct {
	endpoint   *url.URL
	opts       map[string]string
	cookieFile string
	ua         *http.Client
}

func New(opts map[string]string) *Cli {
	homedir := os.Getenv("HOME")
	cookieJar, _ := cookiejar.New(nil)
	endpoint, _ := opts["endpoint"]
	url, _ := url.Parse(strings.TrimRight(endpoint, "/"))

	if project, ok := opts["project"]; ok {
		opts["project"] = strings.ToUpper(project)
	}

	cli := &Cli{
		endpoint:   url,
		opts:       opts,
		cookieFile: fmt.Sprintf("%s/.jira.d/cookies.js", homedir),
		ua:         &http.Client{Jar: cookieJar},
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
	cookies := make([]*http.Cookie, 0)
	err = json.Unmarshal(bytes, &cookies)
	if err != nil {
		log.Error("Failed to parse json from file %s: %s", c.cookieFile, err)
	}
	log.Debug("Loading Cookies: %s", cookies)
	return cookies
}

func (c *Cli) post(uri string, content string) (*http.Response, error) {
	return c.makeRequestWithContent("POST", uri, content)
}

func (c *Cli) put(uri string, content string) (*http.Response, error) {
	return c.makeRequestWithContent("PUT", uri, content)
}

func (c *Cli) makeRequestWithContent(method string, uri string, content string) (*http.Response, error) {
	buffer := bytes.NewBufferString(content)
	req, _ := http.NewRequest(method, uri, buffer)

	log.Info("%s %s", req.Method, req.URL.String())
	if log.IsEnabledFor(logging.DEBUG) {
		logBuffer := bytes.NewBuffer(make([]byte, 0, len(content)))
		req.Write(logBuffer)
		log.Debug("%s", logBuffer)
		// need to recreate the buffer since the offset is now at the end
		// need to be able to rewind the buffer offset, dont know how yet
		req, _ = http.NewRequest(method, uri, bytes.NewBufferString(content))
	}

	if resp, err := c.makeRequest(req); err != nil {
		return nil, err
	} else {
		if resp.StatusCode == 401 {
			if err := c.CmdLogin(); err != nil {
				return nil, err
			}
			req, _ = http.NewRequest(method, uri, bytes.NewBufferString(content))
			return c.makeRequest(req)
		}
		return resp, err
	}
}

func (c *Cli) get(uri string) (*http.Response, error) {
	req, _ := http.NewRequest("GET", uri, nil)
	log.Info("%s %s", req.Method, req.URL.String())
	if log.IsEnabledFor(logging.DEBUG) {
		logBuffer := bytes.NewBuffer(make([]byte, 0))
		req.Write(logBuffer)
		log.Debug("%s", logBuffer)
	}

	if resp, err := c.makeRequest(req); err != nil {
		return nil, err
	} else {
		if resp.StatusCode == 401 {
			if err := c.CmdLogin(); err != nil {
				return nil, err
			}
			return c.makeRequest(req)
		}
		return resp, err
	}
}

func (c *Cli) makeRequest(req *http.Request) (resp *http.Response, err error) {
	req.Header.Set("Content-Type", "application/json")
	if resp, err = c.ua.Do(req); err != nil {
		log.Error("Failed to %s %s: %s", req.Method, req.URL.String(), err)
		return nil, err
	} else {
		if resp.StatusCode < 200 || resp.StatusCode >= 300 && resp.StatusCode != 401 {
			log.Error("response status: %s", resp.Status)
		}

		runtime.SetFinalizer(resp, func(r *http.Response) {
			r.Body.Close()
		})

		if _, ok := resp.Header["Set-Cookie"]; ok {
			c.saveCookies(resp.Cookies())
		}
	}
	return resp, nil
}

func (c *Cli) getTemplate(name string) string {
	if override, ok := c.opts["template"]; ok {
		if _, err := os.Stat(override); err == nil {
			return readFile(override)
		} else {
			if file, err := FindClosestParentPath(fmt.Sprintf(".jira.d/templates/%s", override)); err == nil {
				return readFile(file)
			}
			if dflt, ok := all_templates[override]; ok {
				return dflt
			}
		}
	}
	if file, err := FindClosestParentPath(fmt.Sprintf(".jira.d/templates/%s", name)); err != nil {
		// create-bug etc are special, if we dont find it in the path
		// then just return a generic create template
		if strings.HasPrefix(name, "create-") {
			if file, err := FindClosestParentPath(".jira.d/templates/create"); err != nil {
				return all_templates["create"]
			} else {
				return readFile(file)
			}
		}
		return all_templates[name]
	} else {
		return readFile(file)
	}
}

func (c *Cli) editTemplate(template string, tmpFilePrefix string, templateData map[string]interface{}, templateProcessor func(string) error) error {

	tmpdir := fmt.Sprintf("%s/.jira.d/tmp", os.Getenv("HOME"))
	if err := mkdir(tmpdir); err != nil {
		return err
	}

	fh, err := ioutil.TempFile(tmpdir, tmpFilePrefix)
	if err != nil {
		log.Error("Failed to make temp file in %s: %s", tmpdir, err)
		return err
	}
	defer fh.Close()

	tmpFileName := fmt.Sprintf("%s.yml", fh.Name())
	if err := os.Rename(fh.Name(), tmpFileName); err != nil {
		log.Error("Failed to rename %s to %s: %s", fh.Name(), fmt.Sprintf("%s.yml", fh.Name()), err)
		return err
	}

	err = runTemplate(template, templateData, fh)
	if err != nil {
		return err
	}

	fh.Close()

	editor, ok := c.opts["editor"]
	if !ok {
		editor = os.Getenv("JIRA_EDITOR")
		if editor == "" {
			editor = os.Getenv("EDITOR")
			if editor == "" {
				editor = "vim"
			}
		}
	}

	editing := true
	if val, ok := c.opts["edit"]; ok && val == "false" {
		editing = false
	}
	
	for true {
		if editing {
			log.Debug("Running: %s %s", editor, tmpFileName)
			cmd := exec.Command(editor, tmpFileName)
			cmd.Stdout, cmd.Stderr, cmd.Stdin = os.Stdout, os.Stderr, os.Stdin
			if err := cmd.Run(); err != nil {
				log.Error("Failed to edit template with %s: %s", editor, err)
				if promptYN("edit again?", true) {
					continue
				}
				return err
			}
		}

		edited := make(map[string]interface{})
		if fh, err := ioutil.ReadFile(tmpFileName); err != nil {
			log.Error("Failed to read tmpfile %s: %s", tmpFileName, err)
			if editing && promptYN("edit again?", true) {
				continue
			}
			return err
		} else {
			if err := yaml.Unmarshal(fh, &edited); err != nil {
				log.Error("Failed to parse YAML: %s", err)
				if editing && promptYN("edit again?", true) {
					continue
				}
				return err
			}
		}

		if fixed, err := yamlFixup(edited); err != nil {
			return err
		} else {
			edited = fixed.(map[string]interface{})
		}

		if _, ok := templateData["meta"]; ok {
			mf := templateData["meta"].(map[string]interface{})["fields"]
			if f, ok := edited["fields"].(map[string]interface{}); ok {
				for k, _ := range f {
					if _, ok := mf.(map[string]interface{})[k]; !ok {
						err := fmt.Errorf("Field %s is not editable", k)
						log.Error("%s", err)
						if editing && promptYN("edit again?", true) {
							continue
						}
						return err
					}
				}
			}
		}

		json, err := jsonEncode(edited)
		if err != nil {
			return err
		}

		if err := templateProcessor(json); err != nil {
			log.Error("%s", err)
			if editing && promptYN("edit again?", true) {
				continue
			}
		}
		return nil
	}
	return nil
}

func (c *Cli) Browse(issue string) error {
	if val, ok := c.opts["browse"]; ok && val == "true" {
		if runtime.GOOS == "darwin" {
			return exec.Command("open", fmt.Sprintf("%s/browse/%s", c.endpoint, issue)).Run()
		} else if runtime.GOOS == "linux" {
			return exec.Command("xdg-open", fmt.Sprintf("%s/browse/%s", c.endpoint, issue)).Run()
		}
	}
	return nil
}
