package jiracmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/coryb/figtree"
	"github.com/eroshan/oreo"
	"github.com/go-jira/jira/jiracli"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func CmdUnexportTemplatesRegistry() *jiracli.CommandRegistryEntry {
	opts := ExportTemplatesOptions{}

	return &jiracli.CommandRegistryEntry{
		"Remove unmodified exported templates",
		func(fig *figtree.FigTree, cmd *kingpin.CmdClause) error {
			jiracli.LoadConfigs(cmd, fig, &opts)
			return CmdExportTemplatesUsage(cmd, &opts)
		},
		func(o *oreo.Client, globals *jiracli.GlobalOptions) error {
			if opts.Dir == "" {
				opts.Dir = fmt.Sprintf("%s/.jira.d/templates", jiracli.Homedir())
			}
			return CmdUnexportTemplates(globals, &opts)
		},
	}
}

// CmdUnexportTemplates will remove unmodified templates from export directory
func CmdUnexportTemplates(globals *jiracli.GlobalOptions, opts *ExportTemplatesOptions) error {
	for name, template := range jiracli.AllTemplates {
		if opts.Template != "" && opts.Template != name {
			continue
		}
		templateFile := path.Join(opts.Dir, name)
		if _, err := os.Stat(templateFile); err != nil {
			log.Warning("Skipping %s, not found", templateFile)
			continue
		}
		// open, read, compare
		contents, err := ioutil.ReadFile(templateFile)
		if err != nil {
			return err
		}
		if bytes.Equal([]byte(template), contents) {
			if !globals.Quiet.Value {
				log.Notice("Removing %s, template identical to default", templateFile)
			}
			os.Remove(templateFile)
		} else if !globals.Quiet.Value {
			log.Notice("Skipping %s, found customizations to template", templateFile)
		}
	}
	return nil
}
