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
	"net/http/httputil"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

var (
	log = logging.MustGetLogger("jira")
	// VERSION is the go-jira library version
	VERSION string
)

// Cli is go-jira client object
type Cli struct {
	endpoint   *url.URL
	opts       map[string]interface{}
	cookieFile string
	ua         *http.Client
}

// New creates go-jira client object
func New(opts map[string]interface{}) *Cli {
	homedir := homedir()
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
		cookieFile: filepath.Join(homedir, ".jira.d", "cookies.js"),
		ua: &http.Client{
			Jar:       cookieJar,
			Transport: transport,
		},
	}

	cli.ua.Jar.SetCookies(url, cli.loadCookies())

	return cli
}

func (c *Cli) saveCookies(resp *http.Response) {
	if _, ok := resp.Header["Set-Cookie"]; !ok {
		return
	}

	cookies := resp.Cookies()
	for _, cookie := range cookies {
		if cookie.Domain == "" {
			// if it is host:port then we need to split off port
			parts := strings.Split(resp.Request.URL.Host, ":")
			host := parts[0]
			log.Debugf("Setting DOMAIN to %s for Cookie: %s", host, cookie)
			cookie.Domain = host
		}
	}

	// expiry in one week from now
	expiry := time.Now().Add(24 * 7 * time.Hour)
	for _, cookie := range cookies {
		cookie.Expires = expiry
	}

	if currentCookies := c.loadCookies(); currentCookies != nil {
		currentCookiesByName := make(map[string]*http.Cookie)
		for _, cookie := range currentCookies {
			currentCookiesByName[cookie.Name+cookie.Domain] = cookie
		}

		for _, cookie := range cookies {
			currentCookiesByName[cookie.Name+cookie.Domain] = cookie
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
		panic(err)
	}
	cookies := []*http.Cookie{}
	err = json.Unmarshal(bytes, &cookies)
	if err != nil {
		log.Errorf("Failed to parse json from file %s: %s", c.cookieFile, err)
	}

	if os.Getenv("LOG_TRACE") != "" && log.IsEnabledFor(logging.DEBUG) {
		log.Debugf("Loading Cookies: %s", cookies)
	}
	return cookies
}

func (c *Cli) post(uri string, content string) (*http.Response, error) {
	return c.makeRequestWithContent("POST", uri, content)
}

func (c *Cli) put(uri string, content string) (*http.Response, error) {
	return c.makeRequestWithContent("PUT", uri, content)
}

func (c *Cli) delete(uri string) (resp *http.Response, err error) {
	method := "DELETE"
	req, _ := http.NewRequest(method, uri, nil)
	log.Infof("%s %s", req.Method, req.URL.String())
	if resp, err = c.makeRequest(req); err != nil {
		return nil, err
	}
	if resp.StatusCode == 401 {
		if err = c.CmdLogin(); err != nil {
			return nil, err
		}
		req, _ = http.NewRequest(method, uri, nil)
		return c.makeRequest(req)
	}
	return resp, err
}

func (c *Cli) makeRequestWithContent(method string, uri string, content string) (resp *http.Response, err error) {
	buffer := bytes.NewBufferString(content)
	req, _ := http.NewRequest(method, uri, buffer)

	log.Infof("%s %s", req.Method, req.URL.String())
	if resp, err = c.makeRequest(req); err != nil {
		return nil, err
	}
	if resp.StatusCode == 401 {
		if err = c.CmdLogin(); err != nil {
			return nil, err
		}
		req, _ = http.NewRequest(method, uri, bytes.NewBufferString(content))
		return c.makeRequest(req)
	}
	return resp, err
}

func (c *Cli) get(uri string) (resp *http.Response, err error) {
	req, _ := http.NewRequest("GET", uri, nil)
	log.Infof("%s %s", req.Method, req.URL.String())
	if log.IsEnabledFor(logging.DEBUG) {
		logBuffer := bytes.NewBuffer(make([]byte, 0))
		req.Write(logBuffer)
		log.Debugf("%s", logBuffer)
	}

	if resp, err = c.makeRequest(req); err != nil {
		return nil, err
	}
	if resp.StatusCode == 401 {
		if err := c.CmdLogin(); err != nil {
			return nil, err
		}
		return c.makeRequest(req)
	}
	return resp, err
}

func (c *Cli) makeRequest(req *http.Request) (resp *http.Response, err error) {
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	// this is actually done in http.send but doing it
	// here so we can log it in DumpRequest for debugging
	for _, cookie := range c.ua.Jar.Cookies(req.URL) {
		req.AddCookie(cookie)
	}

	if log.IsEnabledFor(logging.DEBUG) {
		out, _ := httputil.DumpRequest(req, true)
		log.Debugf("Request: %s", out)
	}

	if resp, err = c.ua.Do(req); err != nil {
		log.Errorf("Failed to %s %s: %s", req.Method, req.URL.String(), err)
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 && resp.StatusCode != 401 {
		log.Errorf("response status: %s", resp.Status)
	}

	runtime.SetFinalizer(resp, func(r *http.Response) {
		r.Body.Close()
	})

	if _, ok := resp.Header["Set-Cookie"]; ok {
		c.saveCookies(resp)
	}
	if log.IsEnabledFor(logging.DEBUG) {
		out, _ := httputil.DumpResponse(resp, true)
		log.Debugf("Response: %s", out)
	}
	return resp, nil
}

// GetTemplate will return the text/template for the given command name
func (c *Cli) GetTemplate(name string) string {
	return c.getTemplate(name)
}

func getLookedUpTemplate(name string, dflt string) string {
	if file, err := FindClosestParentPath(filepath.Join(".jira.d", "templates", name)); err == nil {
		return readFile(file)
	}
	if _, err := os.Stat(fmt.Sprintf("/etc/go-jira/templates/%s", name)); err == nil {
		file := fmt.Sprintf("/etc/go-jira/templates/%s", name)
		return readFile(file)
	}
	return dflt
}

func (c *Cli) getTemplate(name string) string {
	if override, ok := c.opts["template"].(string); ok {
		if _, err := os.Stat(override); err == nil {
			return readFile(override)
		}
		if t := getLookedUpTemplate(override, allTemplates[override]); t != "" {
			return t
		}
	}
	// create-bug etc are special, if we dont find it in the path
	// then just return the create template
	if strings.HasPrefix(name, "create-") {
		return getLookedUpTemplate(name, c.getTemplate("create"))
	}
	return getLookedUpTemplate(name, allTemplates[name])
}

// NoChangesFound is an error returned from when editing templates
// and no modifications were made while editing
type NoChangesFound struct{}

func (f NoChangesFound) Error() string {
	return "No changes found, aborting"
}

func (c *Cli) editTemplate(template string, tmpFilePrefix string, templateData map[string]interface{}, templateProcessor func(string) error) error {

	tmpdir := filepath.Join(homedir(), ".jira.d", "tmp")
	if err := mkdir(tmpdir); err != nil {
		return err
	}

	fh, err := ioutil.TempFile(tmpdir, tmpFilePrefix)
	if err != nil {
		log.Errorf("Failed to make temp file in %s: %s", tmpdir, err)
		return err
	}

	oldFileName := fh.Name()
	tmpFileName := fmt.Sprintf("%s.yml", oldFileName)

	// close tmpfile so we can rename on windows
	fh.Close()

	if err := os.Rename(oldFileName, tmpFileName); err != nil {
		log.Errorf("Failed to rename %s to %s: %s", oldFileName, tmpFileName, err)
		return err
	}

	fh, err = os.OpenFile(tmpFileName, os.O_RDWR|os.O_EXCL, 0600)
	if err != nil {
		log.Errorf("Failed to reopen temp file file in %s: %s", tmpFileName, err)
		return err
	}

	defer fh.Close()
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
		var data []byte
		if data, err = ioutil.ReadFile(tmpFileName); err != nil {
			log.Errorf("Failed to read tmpfile %s: %s", tmpFileName, err)
			if editing && promptYN("edit again?", true) {
				continue
			}
			return err
		}
		if err := yaml.Unmarshal(data, &edited); err != nil {
			log.Errorf("Failed to parse YAML: %s", err)
			if editing && promptYN("edit again?", true) {
				continue
			}
			return err
		}

		var fixed interface{}
		if fixed, err = yamlFixup(edited); err != nil {
			return err
		}
		edited = fixed.(map[string]interface{})

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

// Browse will open up your default browser to the provided issue
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

// SaveData will write out the yaml formated --saveFile file with provided data
func (c *Cli) SaveData(data interface{}) error {
	if val, ok := c.opts["saveFile"].(string); ok && val != "" {
		yamlWrite(val, data)
	}
	return nil
}

// ViewIssueWorkLogs gets the worklog data for the given issue
func (c *Cli) ViewIssueWorkLogs(issue string) (interface{}, error) {
	uri := fmt.Sprintf("%s/rest/api/2/issue/%s/worklog", c.endpoint, issue)
	data, err := responseToJSON(c.get(uri))
	if err != nil {
		return nil, err
	}
	return data, nil
}

// ViewIssue will return the details for the given issue id
func (c *Cli) ViewIssue(issue string) (interface{}, error) {
	uri := fmt.Sprintf("%s/rest/api/2/issue/%s", c.endpoint, issue)
	if x := c.expansions(); len(x) > 0 {
		uri = fmt.Sprintf("%s?expand=%s", uri, strings.Join(x, ","))
	}

	data, err := responseToJSON(c.get(uri))
	if err != nil {
		return nil, err
	}
	return data, nil
}

// FindIssues will return a list of issues that match the given options.
// If the "query" option is undefined it will generate a JQL query
// using any/all of the provide options: project, component, assignee,
// issuetype, watcher, reporter, sort
// Further it will restrict the fields being extracted from the jira
// response with the 'queryfields' option
func (c *Cli) FindIssues() (interface{}, error) {
	var query string
	var ok bool
	// project = BAKERY and status not in (Resolved, Closed)
	if query, ok = c.opts["query"].(string); !ok {
		qbuff := bytes.NewBufferString("resolution = unresolved")
		var project string
		if project, ok = c.opts["project"].(string); !ok {
			err := fmt.Errorf("Missing required arguments, either 'query' or 'project' are required")
			log.Errorf("%s", err)
			return nil, err
		}
		qbuff.WriteString(fmt.Sprintf(" AND project = '%s'", project))

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

	fields := []string{"summary"}
	if qf, ok := c.opts["queryfields"].(string); ok {
		fields = strings.Split(qf, ",")
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
	var data interface{}
	if data, err = responseToJSON(c.post(uri, json)); err != nil {
		return nil, err
	}
	return data, nil
}

type RankOrder int

const (
	RANKBEFORE RankOrder = iota
	RANKAFTER  RankOrder = iota
)

func (c *Cli) RankIssue(issue, target string, order RankOrder) error {
	type RankRequest struct {
		Issues []string `json:"issues"`
		Before string   `json:"rankBeforeIssue,omitempty"`
		After  string   `json:"rankAfterIssue,omitempty"`
	}
	req := &RankRequest{
		Issues: []string{
			issue,
		},
	}
	if order == RANKBEFORE {
		req.Before = target
	} else {
		req.After = target
	}

	json, err := jsonEncode(req)
	if err != nil {
		return err
	}

	uri := fmt.Sprintf("%s/rest/agile/1.0/issue/rank", c.endpoint)
	if c.getOptBool("dryrun", false) {
		log.Debugf("PUT: %s", json)
		log.Debugf("Dryrun mode, skipping PUT")
		return nil
	}
	resp, err := c.put(uri, json)
	if err != nil {
		return err
	}
	if resp.StatusCode != 204 {
		return fmt.Errorf("failed to modify issue rank: %s", resp.Status)
	}
	return nil
}

// GetOptString will extract the string from the Cli object options
// otherwise return the provided default
func (c *Cli) GetOptString(optName string, dflt string) string {
	return c.getOptString(optName, dflt)
}

func (c *Cli) getOptString(optName string, dflt string) string {
	if val, ok := c.opts[optName].(string); ok {
		return val
	}
	return dflt
}

// GetOptBool will extract the boolean value from the Client object options
// otherwise return the provided default\
func (c *Cli) GetOptBool(optName string, dflt bool) bool {
	return c.getOptBool(optName, dflt)
}

func (c *Cli) getOptBool(optName string, dflt bool) bool {
	if val, ok := c.opts[optName].(bool); ok {
		return val
	}
	return dflt
}

// expansions returns a comma-separated list of values for field expansion
func (c *Cli) expansions() []string {
	var expansions []string
	if x, ok := c.opts["expand"].(string); ok {
		expansions = strings.Split(x, ",")
	}
	return expansions
}
