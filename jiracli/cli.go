package jiracli

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"reflect"
	"strings"

	"github.com/coryb/figtree"
	"github.com/coryb/oreo"
	"github.com/jinzhu/copier"
	shellquote "github.com/kballard/go-shellquote"
	"github.com/tidwall/gjson"
	"gopkg.in/AlecAivazis/survey.v1"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
	yaml "gopkg.in/coryb/yaml.v2"
)

type Exit struct {
	Code int
}

type GlobalOptions struct {
	AuthenticationMethod figtree.StringOption `yaml:"authentication-method,omitempty" json:"authentication-method,omitempty"`
	Endpoint             figtree.StringOption `yaml:"endpoint,omitempty" json:"endpoint,omitempty"`
	Insecure             figtree.BoolOption   `yaml:"insecure,omitempty" json:"insecure,omitempty"`
	Login                figtree.StringOption `yaml:"login,omitempty" json:"login,omitempty"`
	PasswordSource       figtree.StringOption `yaml:"password-source,omitempty" json:"password-source,omitempty"`
	PasswordDirectory    figtree.StringOption `yaml:"password-directory,omitempty" json:"password-directory,omitempty"`
	Quiet                figtree.BoolOption   `yaml:"quiet,omitempty" json:"quiet,omitempty"`
	SocksProxy           figtree.StringOption `yaml:"socksproxy,omitempty" json:"socksproxy,omitempty"`
	UnixProxy            figtree.StringOption `yaml:"unixproxy,omitempty" json:"unixproxy,omitempty"`
	User                 figtree.StringOption `yaml:"user,omitempty" json:"user,omitempty"`
}

type CommonOptions struct {
	Browse      figtree.BoolOption   `yaml:"browse,omitempty" json:"browse,omitempty"`
	Editor      figtree.StringOption `yaml:"editor,omitempty" json:"editor,omitempty"`
	GJsonQuery  figtree.StringOption `yaml:"gjq,omitempty" json:"gjq,omitempty"`
	SkipEditing figtree.BoolOption   `yaml:"noedit,omitempty" json:"noedit,omitempty"`
	Template    figtree.StringOption `yaml:"template,omitempty" json:"template,omitempty"`
}

type CommandRegistryEntry struct {
	Help        string
	UsageFunc   func(*figtree.FigTree, *kingpin.CmdClause) error
	ExecuteFunc func(*oreo.Client, *GlobalOptions) error
}

type CommandRegistry struct {
	Command string
	Aliases []string
	Entry   *CommandRegistryEntry
	Default bool
}

// either kingpin.Application or kingpin.CmdClause fit this interface
type kingpinAppOrCommand interface {
	Command(string, string) *kingpin.CmdClause
	GetCommand(string) *kingpin.CmdClause
}

var globalCommandRegistry = []CommandRegistry{}

func RegisterCommand(regEntry CommandRegistry) {
	globalCommandRegistry = append(globalCommandRegistry, regEntry)
}

func (o *GlobalOptions) AuthMethod() string {
	if strings.Contains(o.Endpoint.Value, ".atlassian.net") && o.AuthenticationMethod.Source == "default" {
		return "api-token"
	}
	return o.AuthenticationMethod.Value
}

func register(app *kingpin.Application, o *oreo.Client, fig *figtree.FigTree) {
	globals := GlobalOptions{
		User:                 figtree.NewStringOption(os.Getenv("USER")),
		AuthenticationMethod: figtree.NewStringOption("session"),
	}
	app.Flag("endpoint", "Base URI to use for Jira").Short('e').SetValue(&globals.Endpoint)
	app.Flag("insecure", "Disable TLS certificate verification").Short('k').SetValue(&globals.Insecure)
	app.Flag("quiet", "Suppress output to console").Short('Q').SetValue(&globals.Quiet)
	app.Flag("unixproxy", "Path for a unix-socket proxy").SetValue(&globals.UnixProxy)
	app.Flag("socksproxy", "Address for a socks proxy").SetValue(&globals.SocksProxy)
	app.Flag("user", "user name used within the Jira service").Short('u').SetValue(&globals.User)
	app.Flag("login", "login name that corresponds to the user used for authentication").SetValue(&globals.Login)

	o = o.WithPreCallback(
		func(req *http.Request) (*http.Request, error) {
			if globals.AuthMethod() == "api-token" {
				// need to set basic auth header with user@domain:api-token
				token := globals.GetPass()
				authHeader := fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", globals.Login.Value, token))))
				req.Header.Add("Authorization", authHeader)
			}
			return req, nil
		},
	)

	o = o.WithPostCallback(
		func(req *http.Request, resp *http.Response) (*http.Response, error) {
			if globals.AuthMethod() == "session" {
				authUser := resp.Header.Get("X-Ausername")
				if authUser == "" || authUser == "anonymous" {
					// preserve the --quiet value, we need to temporarily disable it so
					// the normal login output is surpressed
					defer func(quiet bool) {
						globals.Quiet.Value = quiet
					}(globals.Quiet.Value)
					globals.Quiet.Value = true

					// we are not logged in, so force login now by running the "login" command
					app.Parse([]string{"login"})

					// rerun the original request
					return o.Do(req)
				}
			}
			return resp, nil
		},
	)

	for _, command := range globalCommandRegistry {
		copy := command
		commandFields := strings.Fields(copy.Command)
		var appOrCmd kingpinAppOrCommand = app
		if len(commandFields) > 1 {
			for _, name := range commandFields[0 : len(commandFields)-1] {
				tmp := appOrCmd.GetCommand(name)
				if tmp == nil {
					tmp = appOrCmd.Command(name, "")
				}
				appOrCmd = tmp
			}
		}

		cmd := appOrCmd.Command(commandFields[len(commandFields)-1], copy.Entry.Help)
		LoadConfigs(cmd, fig, &globals)
		cmd.PreAction(func(_ *kingpin.ParseContext) error {
			if globals.Insecure.Value {
				transport := &http.Transport{
					Proxy: http.ProxyFromEnvironment,
					TLSClientConfig: &tls.Config{
						InsecureSkipVerify: true,
					},
				}
				o = o.WithTransport(transport)
			}
			if globals.UnixProxy.Value != "" {
				o = o.WithTransport(unixProxy(globals.UnixProxy.Value))
			} else if globals.SocksProxy.Value != "" {
				o = o.WithTransport(socksProxy(globals.SocksProxy.Value))
			}
			if globals.AuthMethod() == "api-token" {
				o = o.WithCookieFile("")
			}
			if globals.Login.Value == "" {
				globals.Login = globals.User
			}
			return nil
		})

		for _, alias := range copy.Aliases {
			cmd = cmd.Alias(alias)
		}
		if copy.Default {
			cmd = cmd.Default()
		}
		if copy.Entry.UsageFunc != nil {
			copy.Entry.UsageFunc(fig, cmd)
		}

		cmd.Action(
			func(_ *kingpin.ParseContext) error {
				return copy.Entry.ExecuteFunc(o, &globals)
			},
		)
	}
}

func LoadConfigs(cmd *kingpin.CmdClause, fig *figtree.FigTree, opts interface{}) {
	cmd.PreAction(func(_ *kingpin.ParseContext) error {
		os.Setenv("JIRA_OPERATION", cmd.FullCommand())
		// load command specific configs first
		if err := fig.LoadAllConfigs(strings.Join(strings.Fields(cmd.FullCommand()), "_")+".yml", opts); err != nil {
			return err
		}
		// then load generic configs if not already populated above
		return fig.LoadAllConfigs("config.yml", opts)
	})
}

func BrowseUsage(cmd *kingpin.CmdClause, opts *CommonOptions) {
	cmd.Flag("browse", "Open issue(s) in browser after operation").Short('b').SetValue(&opts.Browse)
}

func EditorUsage(cmd *kingpin.CmdClause, opts *CommonOptions) {
	cmd.Flag("editor", "Editor to use").SetValue(&opts.Editor)
}

func TemplateUsage(cmd *kingpin.CmdClause, opts *CommonOptions) {
	cmd.Flag("template", "Template to use for output").Short('t').SetValue(&opts.Template)
}

func GJsonQueryUsage(cmd *kingpin.CmdClause, opts *CommonOptions) {
	cmd.Flag("gjq", "GJSON Query to filter output, see https://goo.gl/iaYwJ5").SetValue(&opts.GJsonQuery)
}

func (o *CommonOptions) PrintTemplate(data interface{}) error {
	if o.GJsonQuery.Value != "" {
		buf := bytes.NewBufferString("")
		RunTemplate("json", data, buf)
		results := gjson.GetBytes(buf.Bytes(), o.GJsonQuery.Value)
		_, err := os.Stdout.Write([]byte(results.String()))
		os.Stdout.Write([]byte{'\n'})
		return err
	}
	return RunTemplate(o.Template.Value, data, nil)
}

func (o *CommonOptions) editFile(fileName string) (changes bool, err error) {
	var editor string
	for _, ed := range []string{o.Editor.Value, os.Getenv("JIRA_EDITOR"), os.Getenv("EDITOR"), "vim"} {
		if ed != "" {
			editor = ed
			break
		}
	}

	if o.SkipEditing.Value {
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

var EditLoopAbort = fmt.Errorf("Edit Loop aborted by request")

func EditLoop(opts *CommonOptions, input interface{}, output interface{}, submit func() error) error {
	tmpFile, err := tmpTemplate(opts.Template.Value, input)
	if err != nil {
		return err
	}

	confirm := func(dflt bool, msg string) (answer bool) {
		survey.AskOne(
			&survey.Confirm{Message: msg, Default: dflt},
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
		if !opts.SkipEditing.Value {
			changes, err := opts.editFile(tmpFile)
			if err != nil {
				log.Error(err.Error())
				if confirm(true, "Editor reported an error, edit again?") {
					continue
				}
				return EditLoopAbort
			}
			if !changes {
				if !confirm(false, "No changes detected, submit anyway?") {
					return EditLoopAbort
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
		// interface, then we serialize to a new yaml document ... then is
		// parsed as the original document to populate the output struct.  Phew.
		var raw interface{}
		if err := yaml.Unmarshal(data, &raw); err != nil {
			log.Error(err.Error())
			if confirm(true, "Invalid YAML syntax, edit again?") {
				continue
			}
			return EditLoopAbort
		}
		yamlFixup(&raw)
		fixedYAML, err := yaml.Marshal(&raw)
		if err != nil {
			log.Error(err.Error())
			if confirm(true, "Invalid YAML syntax, edit again?") {
				continue
			}
			return EditLoopAbort
		}

		if err := yaml.Unmarshal(fixedYAML, output); err != nil {
			log.Error(err.Error())
			if confirm(true, "Invalid YAML syntax, edit again?") {
				continue
			}
			return EditLoopAbort
		}
		// submit template
		if err := submit(); err != nil {
			log.Error(err.Error())
			if confirm(true, "Jira reported an error, edit again?") {
				continue
			}
			return EditLoopAbort
		}
		break
	}
	return nil
}
