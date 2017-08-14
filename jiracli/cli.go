package jiracli

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/AlecAivazis/survey"
	"github.com/coryb/figtree"
	"github.com/coryb/oreo"
	"github.com/jinzhu/copier"
	shellquote "github.com/kballard/go-shellquote"
	jira "gopkg.in/Netflix-Skunkworks/go-jira.v1"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
	yaml "gopkg.in/coryb/yaml.v2"
	logging "gopkg.in/op/go-logging.v1"
)

var log = logging.MustGetLogger("jira")

type JiraCli struct {
	jira.Jira `yaml:",inline"`
	ConfigDir string
	oreoAgent *oreo.Client
}

func New(configDir string) *JiraCli {
	agent := oreo.New().WithCookieFile(filepath.Join(homedir(), configDir, "cookies.js"))
	return &JiraCli{
		ConfigDir: configDir,
		Jira: jira.Jira{
			UA: agent,
		},
		oreoAgent: agent,
	}
}

type Exit struct {
	Code int
}

type GlobalOptions struct {
	Browse         string `json:"browse,omitempty" yaml:"browse,omitempty"`
	Editor         string `json:"editor,omitempty" yaml:"editor,omitempty"`
	SkipEditing    bool   `json:"noedit,omitempty" yaml:"noedit,omitempty"`
	PasswordSource string `json:"password-source,omitempty" yaml:"password-source,omitempty"`
	Template       string `json:"template,omitempty" yaml:"template,omitempty"`
	User           string `json:"user,omitempty", yaml:"user,omitempty"`
}

func (jc *JiraCli) GlobalUsage(cmd *kingpin.CmdClause, opts *GlobalOptions) error {
	jc.LoadConfigs(cmd, opts)
	cmd.PreAction(func(_ *kingpin.ParseContext) error {
		fig := figtree.NewFigTree()
		fig.EnvPrefix = "JIRA"
		// populate JiraCli fields if defined in configs (ie for Endpoint)
		if err := fig.LoadAllConfigs(path.Join(jc.ConfigDir, "config.yml"), jc); err != nil {
			return err
		}
		if opts.User == "" {
			opts.User = os.Getenv("USER")
		}
		return nil
	})
	cmd.Flag("endpoint", "URI to use for Jira").Short('e').StringVar(&jc.Endpoint)
	cmd.Flag("user", "Login mame used for authentication with Jira service").Short('u').StringVar(&opts.User)
	return nil
}

func (jc *JiraCli) LoadConfigs(cmd *kingpin.CmdClause, opts interface{}) {
	cmd.PreAction(func(_ *kingpin.ParseContext) error {
		fig := figtree.NewFigTree()
		fig.EnvPrefix = "JIRA"
		// load command specific configs first
		if err := fig.LoadAllConfigs(path.Join(jc.ConfigDir, strings.Join(strings.Fields(cmd.FullCommand()), "_")+".yml"), opts); err != nil {
			return err
		}
		// then load generic configs if not already populated above
		return fig.LoadAllConfigs(path.Join(jc.ConfigDir, "config.yml"), opts)
	})
}

func (jc *JiraCli) EditorUsage(cmd *kingpin.CmdClause, opts *GlobalOptions) {
	cmd.Flag("editor", "Editor to use").StringVar(&opts.Editor)
}

func (jc *JiraCli) TemplateUsage(cmd *kingpin.CmdClause, opts *GlobalOptions) {
	cmd.Flag("template", "Template to use for output").Short('t').StringVar(&opts.Template)
}

func (o *GlobalOptions) editFile(fileName string) (changes bool, err error) {
	var editor string
	for _, ed := range []string{o.Editor, os.Getenv("JIRA_EDITOR"), os.Getenv("EDITOR"), "vim"} {
		if ed != "" {
			editor = ed
			break
		}
	}

	if o.SkipEditing {
		return false, nil
	}

	tmpFileNameOrig := fmt.Sprintf("%s.orig", fileName)
	if err := copyFile(fileName, tmpFileNameOrig); err != nil {
		return false, err
	}

	defer func() {
		os.Remove(tmpFileNameOrig)
	}()

	shell, _ := shellquote.Split(editor)
	shell = append(shell, fileName)
	log.Debugf("Running: %#v", shell)
	cmd := exec.Command(shell[0], shell[1:]...)
	cmd.Stdout, cmd.Stderr, cmd.Stdin = os.Stdout, os.Stderr, os.Stdin
	if err := cmd.Run(); err != nil {
		return false, err
	}

	// now we just need to diff the files to see if there are any changes
	var oldHandle, newHandle *os.File
	var oldStat, newStat os.FileInfo
	if oldHandle, err = os.Open(tmpFileNameOrig); err == nil {
		if newHandle, err = os.Open(fileName); err == nil {
			if oldStat, err = oldHandle.Stat(); err == nil {
				if newStat, err = newHandle.Stat(); err == nil {
					// different sizes, so must have changes
					if oldStat.Size() != newStat.Size() {
						return true, err
					}
					oldBuf, newBuf := make([]byte, 1024), make([]byte, 1024)
					var oldCount, newCount int
					// loop though 1024 bytes at a time comparing the buffers for changes
					for err != io.EOF {
						oldCount, _ = oldHandle.Read(oldBuf)
						newCount, err = newHandle.Read(newBuf)
						if oldCount != newCount {
							return true, nil
						}
						if bytes.Compare(oldBuf[:oldCount], newBuf[:newCount]) != 0 {
							return true, nil
						}
					}
					return false, nil
				}
			}
		}
	}
	return false, err
}

func (jc *JiraCli) editLoop(opts *GlobalOptions, input interface{}, output interface{}, submit func() error) error {
	tmpFile, err := jc.tmpTemplate(opts.Template, input)
	if err != nil {
		return err
	}

	confirm := func(msg string) (answer bool) {
		survey.AskOne(
			&survey.Confirm{Message: msg, Default: true},
			&answer,
			nil,
		)
		return
	}

	// we need to copy the original output so that we can restore
	// it on retries in case we try to populate bogus fields that
	// are rejected by the jira service.
	dup := reflect.New(reflect.ValueOf(output).Elem().Type())
	err = copier.Copy(dup.Interface(), output)
	if err != nil {
		return err
	}

	for {
		if !opts.SkipEditing {
			changes, err := opts.editFile(tmpFile)
			if err != nil {
				log.Error(err.Error())
				if confirm("Editor reported an error, edit again?") {
					continue
				}
				panic(Exit{Code: 1})
			}
			if !changes {
				if !confirm("No changes detected, submit anyway?") {
					panic(Exit{Code: 1})
				}
			}
		}
		// parse template
		data, err := ioutil.ReadFile(tmpFile)
		if err != nil {
			return err
		}

		defer func(mapType, iface reflect.Type) {
			yaml.DefaultMapType = mapType
			yaml.IfaceType = iface
		}(yaml.DefaultMapType, yaml.IfaceType)
		yaml.DefaultMapType = reflect.TypeOf(map[string]interface{}{})
		yaml.IfaceType = yaml.DefaultMapType.Elem()

		// restore output incase of retry loop
		err = copier.Copy(output, dup.Interface())
		if err != nil {
			return err
		}

		// HACK HACK HACK we want to trim out all the yaml garbage that is not
		// poplulated, like empty arrays, string values with only a newline,
		// etc.  We need to do this because jira will reject json documents
		// with empty arrays, or empty strings typically.  So here we process
		// the data to a raw interface{} then we fixup the yaml parsed
		// inferface, then we serialize to a new yaml document ... then is
		// parsed as the original document to populate the output struct.  Phew.
		var raw interface{}
		if err := yaml.Unmarshal(data, &raw); err != nil {
			log.Error(err.Error())
			if confirm("Invalid YAML syntax, edit again?") {
				continue
			}
			panic(Exit{Code: 1})
		}
		yamlFixup(&raw)
		fixedYAML, err := yaml.Marshal(&raw)
		if err != nil {
			log.Error(err.Error())
			if confirm("Invalid YAML syntax, edit again?") {
				continue
			}
			panic(Exit{Code: 1})
		}

		if err := yaml.Unmarshal(fixedYAML, output); err != nil {
			log.Error(err.Error())
			if confirm("Invalid YAML syntax, edit again?") {
				continue
			}
			panic(Exit{Code: 1})
		}
		// submit template
		if err := submit(); err != nil {
			log.Error(err.Error())
			if confirm("Jira reported an error, edit again?") {
				continue
			}
			panic(Exit{Code: 1})
		}
		break
	}
	return nil
}

// // New creates go-jira client object
// func New(opts map[string]interface{}) *Cli {
// 	homedir := homedir()
// 	cookieJar, _ := cookiejar.New(nil)
// 	endpoint, _ := opts["endpoint"].(string)
// 	url, _ := url.Parse(strings.TrimRight(endpoint, "/"))

// 	if project, ok := opts["project"].(string); ok {
// 		opts["project"] = strings.ToUpper(project)
// 	}

// 	var ua *http.Client
// 	if unixProxyPath, ok := opts["unixproxy"].(string); ok {
// 		ua = &http.Client{
// 			Jar:       cookieJar,
// 			Transport: UnixProxy(unixProxyPath),
// 		}
// 	} else {
// 		transport := &http.Transport{
// 			Proxy:           http.ProxyFromEnvironment,
// 			TLSClientConfig: &tls.Config{},
// 		}
// 		if insecureSkipVerify, ok := opts["insecure"].(bool); ok {
// 			transport.TLSClientConfig.InsecureSkipVerify = insecureSkipVerify
// 		}

// 		ua = &http.Client{
// 			Jar:       cookieJar,
// 			Transport: transport,
// 		}
// 	}

// 	cli := &Cli{
// 		endpoint:   url,
// 		opts:       opts,
// 		cookieFile: filepath.Join(homedir, ".jira.d", "cookies.js"),
// 		ua:         ua,
// 	}

// 	cli.ua.Jar.SetCookies(url, cli.loadCookies())

// 	return cli
// }

// // NoChangesFound is an error returned from when editing templates
// // and no modifications were made while editing
// type NoChangesFound struct{}

// func (f NoChangesFound) Error() string {
// 	return "No changes found, aborting"
// }

// func (c *Cli) editTemplate(template string, tmpFilePrefix string, templateData map[string]interface{}, templateProcessor func(string) error) error {

// 	tmpdir := filepath.Join(homedir(), ".jira.d", "tmp")
// 	if err := mkdir(tmpdir); err != nil {
// 		return err
// 	}

// 	fh, err := ioutil.TempFile(tmpdir, tmpFilePrefix)
// 	if err != nil {
// 		log.Errorf("Failed to make temp file in %s: %s", tmpdir, err)
// 		return err
// 	}

// 	oldFileName := fh.Name()
// 	tmpFileName := fmt.Sprintf("%s.yml", oldFileName)

// 	// close tmpfile so we can rename on windows
// 	fh.Close()

// 	if err := os.Rename(oldFileName, tmpFileName); err != nil {
// 		log.Errorf("Failed to rename %s to %s: %s", oldFileName, tmpFileName, err)
// 		return err
// 	}

// 	fh, err = os.OpenFile(tmpFileName, os.O_RDWR|os.O_EXCL, 0600)
// 	if err != nil {
// 		log.Errorf("Failed to reopen temp file file in %s: %s", tmpFileName, err)
// 		return err
// 	}

// 	defer fh.Close()
// 	defer func() {
// 		os.Remove(tmpFileName)
// 	}()

// 	err = runTemplate(template, templateData, fh)
// 	if err != nil {
// 		return err
// 	}

// 	fh.Close()

// 	editor, ok := c.opts["editor"].(string)
// 	if !ok {
// 		editor = os.Getenv("JIRA_EDITOR")
// 		if editor == "" {
// 			editor = os.Getenv("EDITOR")
// 			if editor == "" {
// 				editor = "vim"
// 			}
// 		}
// 	}

// 	editing := c.getOptBool("edit", true)

// 	tmpFileNameOrig := fmt.Sprintf("%s.orig", tmpFileName)
// 	copyFile(tmpFileName, tmpFileNameOrig)
// 	defer func() {
// 		os.Remove(tmpFileNameOrig)
// 	}()

// 	for true {
// 		if editing {
// 			shell, _ := shellquote.Split(editor)
// 			shell = append(shell, tmpFileName)
// 			log.Debugf("Running: %#v", shell)
// 			cmd := exec.Command(shell[0], shell[1:]...)
// 			cmd.Stdout, cmd.Stderr, cmd.Stdin = os.Stdout, os.Stderr, os.Stdin
// 			if err := cmd.Run(); err != nil {
// 				log.Errorf("Failed to edit template with %s: %s", editor, err)
// 				if promptYN("edit again?", true) {
// 					continue
// 				}
// 				return err
// 			}

// 			diff := exec.Command("diff", "-q", tmpFileNameOrig, tmpFileName)
// 			// if err == nil then diff found no changes
// 			if err := diff.Run(); err == nil {
// 				return NoChangesFound{}
// 			}
// 		}

// 		edited := make(map[string]interface{})
// 		var data []byte
// 		if data, err = ioutil.ReadFile(tmpFileName); err != nil {
// 			log.Errorf("Failed to read tmpfile %s: %s", tmpFileName, err)
// 			if editing && promptYN("edit again?", true) {
// 				continue
// 			}
// 			return err
// 		}
// 		if err := yaml.Unmarshal(data, &edited); err != nil {
// 			log.Errorf("Failed to parse YAML: %s", err)
// 			if editing && promptYN("edit again?", true) {
// 				continue
// 			}
// 			return err
// 		}

// 		var fixed interface{}
// 		if fixed, err = yamlFixup(edited); err != nil {
// 			return err
// 		}
// 		edited = fixed.(map[string]interface{})

// 		// if you want to abort editing a jira issue then
// 		// you can add the "abort: true" flag to the document
// 		// and we will abort now
// 		if val, ok := edited["abort"].(bool); ok && val {
// 			log.Infof("abort flag found in template, quiting")
// 			return fmt.Errorf("abort flag found in template, quiting")
// 		}

// 		if _, ok := templateData["meta"]; ok {
// 			mf := templateData["meta"].(map[string]interface{})["fields"]
// 			if f, ok := edited["fields"].(map[string]interface{}); ok {
// 				for k := range f {
// 					if _, ok := mf.(map[string]interface{})[k]; !ok {
// 						err := fmt.Errorf("Field %s is not editable", k)
// 						log.Errorf("%s", err)
// 						if editing && promptYN("edit again?", true) {
// 							continue
// 						}
// 						return err
// 					}
// 				}
// 			}
// 		}

// 		json, err := jsonEncode(edited)
// 		if err != nil {
// 			return err
// 		}

// 		if err := templateProcessor(json); err != nil {
// 			log.Errorf("%s", err)
// 			if editing && promptYN("edit again?", true) {
// 				continue
// 			}
// 		}
// 		return nil
// 	}
// 	return nil
// }

// // Browse will open up your default browser to the provided issue
// func (c *Cli) Browse(issue string) error {
// 	if val, ok := c.opts["browse"].(bool); ok && val {
// 		if runtime.GOOS == "darwin" {
// 			return exec.Command("open", fmt.Sprintf("%s/browse/%s", c.endpoint, issue)).Run()
// 		} else if runtime.GOOS == "linux" {
// 			return exec.Command("xdg-open", fmt.Sprintf("%s/browse/%s", c.endpoint, issue)).Run()
// 		} else if runtime.GOOS == "windows" {
// 			return exec.Command("cmd", "/c", "start", fmt.Sprintf("%s/browse/%s", c.endpoint, issue)).Run()
// 		}
// 	}
// 	return nil
// }

// // SaveData will write out the yaml formated --saveFile file with provided data
// func (c *Cli) SaveData(data interface{}) error {
// 	if val, ok := c.opts["saveFile"].(string); ok && val != "" {
// 		yamlWrite(val, data)
// 	}
// 	return nil
// }

// // FindIssues will return a list of issues that match the given options.
// // If the "query" option is undefined it will generate a JQL query
// // using any/all of the provide options: project, component, assignee,
// // issuetype, watcher, reporter, sort
// // Further it will restrict the fields being extracted from the jira
// // response with the 'queryfields' option
// func (c *Cli) FindIssues() (interface{}, error) {
// 	var query string
// 	var ok bool
// 	// project = BAKERY and status not in (Resolved, Closed)
// 	if query, ok = c.opts["query"].(string); !ok {
// 		qbuff := bytes.NewBufferString("resolution = unresolved")
// 		var project string
// 		if project, ok = c.opts["project"].(string); !ok {
// 			err := fmt.Errorf("Missing required arguments, either 'query' or 'project' are required")
// 			log.Errorf("%s", err)
// 			return nil, err
// 		}
// 		qbuff.WriteString(fmt.Sprintf(" AND project = '%s'", project))

// 		if component, ok := c.opts["component"]; ok {
// 			qbuff.WriteString(fmt.Sprintf(" AND component = '%s'", component))
// 		}

// 		if assignee, ok := c.opts["assignee"]; ok {
// 			qbuff.WriteString(fmt.Sprintf(" AND assignee = '%s'", assignee))
// 		}

// 		if issuetype, ok := c.opts["issuetype"]; ok {
// 			qbuff.WriteString(fmt.Sprintf(" AND issuetype = '%s'", issuetype))
// 		}

// 		if watcher, ok := c.opts["watcher"]; ok {
// 			qbuff.WriteString(fmt.Sprintf(" AND watcher = '%s'", watcher))
// 		}

// 		if reporter, ok := c.opts["reporter"]; ok {
// 			qbuff.WriteString(fmt.Sprintf(" AND reporter = '%s'", reporter))
// 		}

// 		if sort, ok := c.opts["sort"]; ok && sort != "" {
// 			qbuff.WriteString(fmt.Sprintf(" ORDER BY %s", sort))
// 		}

// 		query = qbuff.String()
// 	}

// 	fields := []string{"summary"}
// 	if qf, ok := c.opts["queryfields"].(string); ok {
// 		fields = strings.Split(qf, ",")
// 	}

// 	json, err := jsonEncode(map[string]interface{}{
// 		"jql":        query,
// 		"startAt":    c.opts["start_at"],
// 		"maxResults": c.opts["max_results"],
// 		"fields":     fields,
// 		"expand":     c.expansions(),
// 	})
// 	if err != nil {
// 		return nil, err
// 	}

// 	uri := fmt.Sprintf("%s/rest/api/2/search", c.endpoint)
// 	var data interface{}
// 	if data, err = responseToJSON(c.post(uri, json)); err != nil {
// 		return nil, err
// 	}
// 	return data, nil
// }

// // GetOptString will extract the string from the Cli object options
// // otherwise return the provided default
// func (c *Cli) GetOptString(optName string, dflt string) string {
// 	return c.getOptString(optName, dflt)
// }

// func (c *Cli) getOptString(optName string, dflt string) string {
// 	if val, ok := c.opts[optName].(string); ok {
// 		return val
// 	}
// 	return dflt
// }

// // GetOptBool will extract the boolean value from the Client object options
// // otherwise return the provided default\
// func (c *Cli) GetOptBool(optName string, dflt bool) bool {
// 	return c.getOptBool(optName, dflt)
// }

// func (c *Cli) getOptBool(optName string, dflt bool) bool {
// 	if val, ok := c.opts[optName].(bool); ok {
// 		return val
// 	}
// 	return dflt
// }

// // expansions returns a comma-separated list of values for field expansion
// func (c *Cli) expansions() []string {
// 	var expansions []string
// 	if x, ok := c.opts["expand"].(string); ok {
// 		expansions = strings.Split(x, ",")
// 	}
// 	return expansions
// }
