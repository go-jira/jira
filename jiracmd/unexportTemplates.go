package jiracmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/coryb/figtree"
	"gopkg.in/Netflix-Skunkworks/go-jira.v1/jiracli"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func CmdUnexportTemplatesRegistry() *jiracli.CommandRegistryEntry {
	opts := ExportTemplatesOptions{}

	return &jiracli.CommandRegistryEntry{
		"Remove unmodified exported templates",
		func(fig *figtree.FigTree, cmd *kingpin.CmdClause) error {
			jiracli.LoadConfigs(cmd, fig, &opts)
			if opts.Dir != "" {
				opts.Dir = fmt.Sprintf("%s/.jira.d/templates", jiracli.Homedir())
			}

			return CmdExportTemplatesUsage(cmd, &opts)
		},
		func(globals *jiracli.GlobalOptions) error {
			return CmdUnexportTemplates(&opts)
		},
	}
}

// CmdUnexportTemplates will remove unmodified templates from export directory
func CmdUnexportTemplates(opts *ExportTemplatesOptions) error {
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
		if bytes.Compare([]byte(template), contents) == 0 {
			log.Warning("Removing %s, template identical to default", templateFile)
			os.Remove(templateFile)
		} else {
			log.Warning("Skipping %s, found customizations to template", templateFile)
		}
	}
	return nil
}
