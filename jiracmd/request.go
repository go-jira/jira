package jiracmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"strings"

	"github.com/coryb/figtree"
	"github.com/coryb/oreo"
	"github.com/go-jira/jira/jiracli"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type RequestOptions struct {
	jiracli.CommonOptions `yaml:",inline" json:",inline" figtree:",inline"`
	Method                string `yaml:"method,omitempty" json:"method,omitempty"`
	URI                   string `yaml:"uri,omitempty" json:"uri,omitempty"`
	Data                  string `yaml:"data,omitempty" json:"data,omitempty"`
}

func CmdRequestRegistry() *jiracli.CommandRegistryEntry {
	opts := RequestOptions{
		CommonOptions: jiracli.CommonOptions{
			Template: figtree.NewStringOption("request"),
		},
	}

	return &jiracli.CommandRegistryEntry{
		"Open issue in requestr",
		func(fig *figtree.FigTree, cmd *kingpin.CmdClause) error {
			jiracli.LoadConfigs(cmd, fig, &opts)
			jiracli.TemplateUsage(cmd, &opts.CommonOptions)
			jiracli.GJsonQueryUsage(cmd, &opts.CommonOptions)
			return CmdRequestUsage(cmd, &opts)
		},
		func(o *oreo.Client, globals *jiracli.GlobalOptions) error {
			if opts.Method == "" {
				opts.Method = "GET"
			}
			return CmdRequest(o, globals, &opts)
		},
	}
}

func CmdRequestUsage(cmd *kingpin.CmdClause, opts *RequestOptions) error {
	cmd.Flag("method", "HTTP request method to use").Short('M').EnumVar(&opts.Method, "GET", "PUT", "POST", "DELETE")
	cmd.Arg("API", "Path to Jira API (ie: /rest/api/2/issue)").Required().StringVar(&opts.URI)
	cmd.Arg("JSON", "JSON Content to send to API").StringVar(&opts.Data)

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
	if resp.Body == nil {
		return fmt.Errorf("Empty Response Body")
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Response Body read Error: %v", err)
	}
	if len(bodyBytes) == 0 {
		log.Info("Empty response for status %d", resp.StatusCode)
		return nil
	}

	var data interface{}
	if err := json.Unmarshal(bodyBytes, &data); err != nil {
		return fmt.Errorf("JSON Parse Error: %v", err)
	}
	return opts.PrintTemplate(&data)
}
