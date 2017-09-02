package jiracmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"strings"

	"github.com/coryb/figtree"
	"github.com/coryb/oreo"

	"gopkg.in/Netflix-Skunkworks/go-jira.v1/jiracli"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type RequestOptions struct {
	jiracli.CommonOptions `yaml:",inline" json:",inline" figtree:",inline"`
	Method                string `yaml:"method,omitempty" json:"method,omitempty"`
	URI                   string `yaml:"uri,omitempty" json:"uri,omitempty"`
	Data                  string `yaml:"data,omitempty" json:"data,omitempty"`
}

func CmdRequestRegistry(o *oreo.Client) *jiracli.CommandRegistryEntry {
	opts := RequestOptions{
		CommonOptions: jiracli.CommonOptions{
			Template: figtree.NewStringOption("request"),
		},
	}

	return &jiracli.CommandRegistryEntry{
		"Open issue in requestr",
		func(fig *figtree.FigTree, cmd *kingpin.CmdClause) error {
			jiracli.LoadConfigs(cmd, fig, &opts)
			if opts.Method == "" {
				opts.Method = "GET"
			}
			return CmdRequestUsage(cmd, &opts)
		},
		func(globals *jiracli.GlobalOptions) error {
			return CmdRequest(o, globals, &opts)
		},
	}
}

func CmdRequestUsage(cmd *kingpin.CmdClause, opts *RequestOptions) error {
	cmd.Flag("method", "HTTP request method to use").Short('m').EnumVar(&opts.Method, "GET", "PUT", "POST", "DELETE")
	cmd.Arg("API", "Path to Jira API (ie: /rest/api/2/issue)").Required().StringVar(&opts.URI)
	cmd.Arg("JSON", "JSON Content to send to API").Required().StringVar(&opts.Data)

	return nil
}

// CmdRequest open the default system requestr to the provided issue
func CmdRequest(o *oreo.Client, globals *jiracli.GlobalOptions, opts *RequestOptions) error {
	uri := opts.URI
	if !strings.HasPrefix(uri, "http") {
		uri = globals.Endpoint.Value + uri
	}

	parsedURI, err := url.Parse(uri)
	if err != nil {
		return err
	}
	builder := oreo.RequestBuilder(parsedURI).WithMethod(opts.Method)
	if opts.Data != "" {
		builder = builder.WithJSON(opts.Data)
	}

	resp, err := o.Do(builder.Build())
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if len(content) == 0 {
		fmt.Println("No Content")
		return nil
	}
	var data interface{}
	err = json.Unmarshal(content, &data)
	if err != nil {
		return fmt.Errorf("JSON Parse Error: %s from %q", err, content)
	}

	return jiracli.RunTemplate(opts.Template.Value, &data, nil)
}
