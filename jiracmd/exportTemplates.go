package jiracmd

import (
	"fmt"
	"os"
	"path"

	"github.com/coryb/figtree"
	"github.com/eroshan/oreo"
	"github.com/go-jira/jira/jiracli"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type ExportTemplatesOptions struct {
	Template string `yaml:"template,omitempty" json:"template,omitempty"`
	Dir      string `yaml:"dir,omitempty" json:"dir,omitempty"`
}

func CmdExportTemplatesRegistry() *jiracli.CommandRegistryEntry {
	opts := ExportTemplatesOptions{}

	return &jiracli.CommandRegistryEntry{
		"Export templates for customizations",
		func(fig *figtree.FigTree, cmd *kingpin.CmdClause) error {
			jiracli.LoadConfigs(cmd, fig, &opts)
			return CmdExportTemplatesUsage(cmd, &opts)
		},
		func(o *oreo.Client, globals *jiracli.GlobalOptions) error {
			if opts.Dir == "" {
				opts.Dir = fmt.Sprintf("%s/.jira.d/templates", jiracli.Homedir())
			}
			return CmdExportTemplates(globals, &opts)
		},
	}
}

func CmdExportTemplatesUsage(cmd *kingpin.CmdClause, opts *ExportTemplatesOptions) error {
	cmd.Flag("template", "Template to export").Short('t').StringVar(&opts.Template)
	cmd.Flag("dir", "directory to write tempates to").Short('d').StringVar(&opts.Dir)

	return nil
}

// CmdExportTemplates will export templates to directory
func CmdExportTemplates(globals *jiracli.GlobalOptions, opts *ExportTemplatesOptions) error {
	if err := os.MkdirAll(opts.Dir, 0755); err != nil {
		return err
	}

	for name, template := range jiracli.AllTemplates {
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
		if !globals.Quiet.Value {
			log.Noticef("Creating %s", templateFile)
		}
		fh.Write([]byte(template))
	}
	return nil
}
