package jiracli

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/coryb/figtree"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func CmdUnexportTemplatesRegistry(fig *figtree.FigTree) *CommandRegistryEntry {
	opts := ExportTemplatesOptions{
		Dir: fmt.Sprintf("%s/.jira.d/templates", Homedir()),
	}

	return &CommandRegistryEntry{
		"Remove unmodified exported templates",
		func() error {
			return CmdUnexportTemplates(&opts)
		},
		func(cmd *kingpin.CmdClause) error {
			LoadConfigs(cmd, fig, &opts)
			return CmdExportTemplatesUsage(cmd, &opts)
		},
	}
}

// CmdUnexportTemplates will remove unmodified templates from export directory
func CmdUnexportTemplates(opts *ExportTemplatesOptions) error {
	for name, template := range allTemplates {
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
