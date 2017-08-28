package jiracli

import (
	"fmt"
	"os"
	"path"

	"github.com/coryb/figtree"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type ExportTemplatesOptions struct {
	Template string `yaml:"template,omitempty" json:"template,omitempty"`
	Dir      string `yaml:"dir,omitempty" json:"dir,omitempty"`
}

func CmdExportTemplatesRegistry(fig *figtree.FigTree) *CommandRegistryEntry {
	opts := ExportTemplatesOptions{}

	return &CommandRegistryEntry{
		"Export templates for customizations",
		func() error {
			return CmdExportTemplates(&opts)
		},
		func(cmd *kingpin.CmdClause) error {
			LoadConfigs(cmd, fig, &opts)
			if opts.Dir == "" {
				opts.Dir = fmt.Sprintf("%s/.jira.d/templates", Homedir())
			}
			return CmdExportTemplatesUsage(cmd, &opts)
		},
	}
}

func CmdExportTemplatesUsage(cmd *kingpin.CmdClause, opts *ExportTemplatesOptions) error {
	cmd.Flag("template", "Template to export").Short('t').StringVar(&opts.Template)
	cmd.Flag("dir", "directory to write tempates to").Short('d').StringVar(&opts.Dir)

	return nil
}

// CmdExportTemplates will export templates to directory
func CmdExportTemplates(opts *ExportTemplatesOptions) error {
	if err := os.MkdirAll(opts.Dir, 0755); err != nil {
		return err
	}

	for name, template := range allTemplates {
		if opts.Template != "" && opts.Template != name {
			continue
		}
		templateFile := path.Join(opts.Dir, name)
		if _, err := os.Stat(templateFile); err == nil {
			log.Warning("Skipping %s, already exists", templateFile)
			continue
		}
		fh, err := os.OpenFile(templateFile, os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			log.Errorf("Failed to open %s for writing: %s", templateFile, err)
			return err
		}
		defer fh.Close()
		log.Noticef("Creating %s", templateFile)
		fh.Write([]byte(template))
	}
	return nil
}
