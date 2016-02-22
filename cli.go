package jira

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/kballard/go-shellquote"
	"gopkg.in/coryb/yaml.v2"
	"gopkg.in/op/go-logging.v1"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"
	"time"
)

var (
	log = logging.MustGetLogger("jira")
	VERSION string
)

type Cli struct {
	endpoint   *url.URL
	opts       map[string]interface{}
	cookieFile string
	ua         *http.Client
}

func New(opts map[string]interface{}) *Cli {
	homedir := os.Getenv("HOME")
	cookieJar, _ := cookiejar.New(nil)
	endpoint, _ := opts["endpoint"].(string)
	url, _ := url.Parse(strings.TrimRight(endpoint, "/"))

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{},
	}

	if project, ok := opts["project"].(string); ok {
		opts["project"] = strings.ToUpper(project)
	}

	if insecureSkipVerify, ok := opts["insecure"].(bool); ok {
		transport.TLSClientConfig.InsecureSkipVerify = insecureSkipVerify
	}

	cli := &Cli{
		endpoint:   url,
		opts:       opts,
		cookieFile: fmt.Sprintf("%s/.jira.d/cookies.js", homedir),
		ua: &http.Client{
			Jar:       cookieJar,
			Transport: transport,
		},
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
		mkdir(path.Dir(c.cookieFile))
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
		log.Errorf("Failed to open %s: %s", c.cookieFile, err)
		os.Exit(1)
	}
	cookies := make([]*http.Cookie, 0)
	err = json.Unmarshal(bytes, &cookies)
	if err != nil {
		log.Errorf("Failed to parse json from file %s: %s", c.cookieFile, err)
	}
	log.Debugf("Loading Cookies: %s", cookies)
	return cookies
}

func (c *Cli) post(uri string, content string) (*http.Response, error) {
	return c.makeRequestWithContent("POST", uri, content)
}

func (c *Cli) put(uri string, content string) (*http.Response, error) {
	return c.makeRequestWithContent("PUT", uri, content)
}

func (c *Cli) delete(uri string) (*http.Response, error) {
	method := "DELETE"
	req, _ := http.NewRequest(method, uri, nil)
	log.Infof("%s %s", req.Method, req.URL.String())
	if resp, err := c.makeRequest(req); err != nil {
		return nil, err
	} else {
		if resp.StatusCode == 401 {
			if err := c.CmdLogin(); err != nil {
				return nil, err
			}
			req, _ = http.NewRequest(method, uri, nil)
			return c.makeRequest(req)
		}
		return resp, err
	}
}

func (c *Cli) makeRequestWithContent(method string, uri string, content string) (*http.Response, error) {
	buffer := bytes.NewBufferString(content)
	req, _ := http.NewRequest(method, uri, buffer)

	log.Infof("%s %s", req.Method, req.URL.String())
	if log.IsEnabledFor(logging.DEBUG) {
		logBuffer := bytes.NewBuffer(make([]byte, 0, len(content)))
		req.Write(logBuffer)
		log.Debugf("%s", logBuffer)
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
	log.Infof("%s %s", req.Method, req.URL.String())
	if log.IsEnabledFor(logging.DEBUG) {
		logBuffer := bytes.NewBuffer(make([]byte, 0))
		req.Write(logBuffer)
		log.Debugf("%s", logBuffer)
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
		log.Errorf("Failed to %s %s: %s", req.Method, req.URL.String(), err)
		return nil, err
	} else {
		if resp.StatusCode < 200 || resp.StatusCode >= 300 && resp.StatusCode != 401 {
			log.Errorf("response status: %s", resp.Status)
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

func (c *Cli) GetTemplate(name string) string {
	return c.getTemplate(name)
}

func (c *Cli) getTemplate(name string) string {
	if override, ok := c.opts["template"].(string); ok {
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

type NoChangesFound struct{}

func (f NoChangesFound) Error() string {
	return "No changes found, aborting"
}

func (c *Cli) editTemplate(template string, tmpFilePrefix string, templateData map[string]interface{}, templateProcessor func(string) error) error {

	tmpdir := fmt.Sprintf("%s/.jira.d/tmp", os.Getenv("HOME"))
	if err := mkdir(tmpdir); err != nil {
		return err
	}

	fh, err := ioutil.TempFile(tmpdir, tmpFilePrefix)
	if err != nil {
		log.Errorf("Failed to make temp file in %s: %s", tmpdir, err)
		return err
	}
	defer fh.Close()

	tmpFileName := fmt.Sprintf("%s.yml", fh.Name())
	if err := os.Rename(fh.Name(), tmpFileName); err != nil {
		log.Errorf("Failed to rename %s to %s: %s", fh.Name(), fmt.Sprintf("%s.yml", fh.Name()), err)
		return err
	}
	defer func() {
		os.Remove(tmpFileName)
	}()

	err = runTemplate(template, templateData, fh)
	if err != nil {
		return err
	}

	fh.Close()

	editor, ok := c.opts["editor"].(string)
	if !ok {
		editor = os.Getenv("JIRA_EDITOR")
		if editor == "" {
			editor = os.Getenv("EDITOR")
			if editor == "" {
				editor = "vim"
			}
		}
	}

	editing := c.getOptBool("edit", true)

	tmpFileNameOrig := fmt.Sprintf("%s.orig", tmpFileName)
	copyFile(tmpFileName, tmpFileNameOrig)
	defer func() {
		os.Remove(tmpFileNameOrig)
	}()

	for true {
		if editing {
			shell, _ := shellquote.Split(editor)
			shell = append(shell, tmpFileName)
			log.Debugf("Running: %#v", shell)
			cmd := exec.Command(shell[0], shell[1:]...)
			cmd.Stdout, cmd.Stderr, cmd.Stdin = os.Stdout, os.Stderr, os.Stdin
			if err := cmd.Run(); err != nil {
				log.Errorf("Failed to edit template with %s: %s", editor, err)
				if promptYN("edit again?", true) {
					continue
				}
				return err
			}

			diff := exec.Command("diff", "-q", tmpFileNameOrig, tmpFileName)
			// if err == nil then diff found no changes
			if err := diff.Run(); err == nil {
				return NoChangesFound{}
			}
		}

		edited := make(map[string]interface{})
		if fh, err := ioutil.ReadFile(tmpFileName); err != nil {
			log.Errorf("Failed to read tmpfile %s: %s", tmpFileName, err)
			if editing && promptYN("edit again?", true) {
				continue
			}
			return err
		} else {
			if err := yaml.Unmarshal(fh, &edited); err != nil {
				log.Errorf("Failed to parse YAML: %s", err)
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

		// if you want to abort editing a jira issue then
		// you can add the "abort: true" flag to the document
		// and we will abort now
		if val, ok := edited["abort"].(bool); ok && val {
			log.Infof("abort flag found in template, quiting")
			return fmt.Errorf("abort flag found in template, quiting")
		}

		if _, ok := templateData["meta"]; ok {
			mf := templateData["meta"].(map[string]interface{})["fields"]
			if f, ok := edited["fields"].(map[string]interface{}); ok {
				for k := range f {
					if _, ok := mf.(map[string]interface{})[k]; !ok {
						err := fmt.Errorf("Field %s is not editable", k)
						log.Errorf("%s", err)
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
			log.Errorf("%s", err)
			if editing && promptYN("edit again?", true) {
				continue
			}
		}
		return nil
	}
	return nil
}

func (c *Cli) Browse(issue string) error {
	if val, ok := c.opts["browse"].(bool); ok && val {
		if runtime.GOOS == "darwin" {
			return exec.Command("open", fmt.Sprintf("%s/browse/%s", c.endpoint, issue)).Run()
		} else if runtime.GOOS == "linux" {
			return exec.Command("xdg-open", fmt.Sprintf("%s/browse/%s", c.endpoint, issue)).Run()
		}
	}
	return nil
}

func (c *Cli) SaveData(data interface{}) error {
	if val, ok := c.opts["saveFile"].(string); ok && val != "" {
		yamlWrite(val, data)
	}
	return nil
}

func (c *Cli) ViewIssue(issue string) (interface{}, error) {
	uri := fmt.Sprintf("%s/rest/api/2/issue/%s", c.endpoint, issue)
	if x := c.expansions(); len(x) > 0 {
		uri = fmt.Sprintf("%s?expand=%s", uri, strings.Join(x, ","))
	}

	data, err := responseToJson(c.get(uri))
	if err != nil {
		return nil, err
	} else {
		return data, nil
	}
}

func (c *Cli) FindIssues() (interface{}, error) {
	var query string
	var ok bool
	// project = BAKERY and status not in (Resolved, Closed)
	if query, ok = c.opts["query"].(string); !ok {
		qbuff := bytes.NewBufferString("resolution = unresolved")
		if project, ok := c.opts["project"]; !ok {
			err := fmt.Errorf("Missing required arguments, either 'query' or 'project' are required")
			log.Errorf("%s", err)
			return nil, err
		} else {
			qbuff.WriteString(fmt.Sprintf(" AND project = '%s'", project))
		}

		if component, ok := c.opts["component"]; ok {
			qbuff.WriteString(fmt.Sprintf(" AND component = '%s'", component))
		}

		if assignee, ok := c.opts["assignee"]; ok {
			qbuff.WriteString(fmt.Sprintf(" AND assignee = '%s'", assignee))
		}

		if issuetype, ok := c.opts["issuetype"]; ok {
			qbuff.WriteString(fmt.Sprintf(" AND issuetype = '%s'", issuetype))
		}

		if watcher, ok := c.opts["watcher"]; ok {
			qbuff.WriteString(fmt.Sprintf(" AND watcher = '%s'", watcher))
		}

		if reporter, ok := c.opts["reporter"]; ok {
			qbuff.WriteString(fmt.Sprintf(" AND reporter = '%s'", reporter))
		}

		if sort, ok := c.opts["sort"]; ok && sort != "" {
			qbuff.WriteString(fmt.Sprintf(" ORDER BY %s", sort))
		}

		query = qbuff.String()
	}

	fields := make([]string, 0)
	if qf, ok := c.opts["queryfields"].(string); ok {
		fields = strings.Split(qf, ",")
	} else {
		fields = append(fields, "summary")
	}

	json, err := jsonEncode(map[string]interface{}{
		"jql":        query,
		"startAt":    "0",
		"maxResults": c.opts["max_results"],
		"fields":     fields,
		"expand":     c.expansions(),
	})
	if err != nil {
		return nil, err
	}

	uri := fmt.Sprintf("%s/rest/api/2/search", c.endpoint)
	if data, err := responseToJson(c.post(uri, json)); err != nil {
		return nil, err
	} else {
		return data, nil
	}
}

func (c *Cli) GetOptString(optName string, dflt string) string {
	return c.getOptString(optName, dflt)
}

func (c *Cli) getOptString(optName string, dflt string) string {
	if val, ok := c.opts[optName].(string); ok {
		return val
	} else {
		return dflt
	}
}

func (c *Cli) GetOptBool(optName string, dflt bool) bool {
	return c.getOptBool(optName, dflt)
}

func (c *Cli) getOptBool(optName string, dflt bool) bool {
	if val, ok := c.opts[optName].(bool); ok {
		return val
	} else {
		return dflt
	}
}

// expansions returns a comma-separated list of values for field expansion
func (c *Cli) expansions() []string {
	expansions := make([]string, 0)
	if x, ok := c.opts["expand"].(string); ok {
		expansions = strings.Split(x, ",")
	}
	return expansions
}