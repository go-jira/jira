package jiracli

import (
	"fmt"
	"os"
	"path"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type ExportTemplatesOptions struct {
	Template string
	Dir      string
}

func (jc *JiraCli) CmdExportTemplatesRegistry() *CommandRegistryEntry {
	opts := ExportTemplatesOptions{
		Dir: fmt.Sprintf("%s/.jira.d/templates", homedir()),
	}

	return &CommandRegistryEntry{
		"Export templates for customizations",
		func() error {
			return jc.CmdExportTemplates(&opts)
		},
		func(cmd *kingpin.CmdClause) error {
			return jc.CmdExportTemplatesUsage(cmd, &opts)
		},
	}
}

func (jc *JiraCli) CmdExportTemplatesUsage(cmd *kingpin.CmdClause, opts *ExportTemplatesOptions) error {
	cmd.Flag("template", "Template to export").Short('t').StringVar(&opts.Template)
	cmd.Flag("dir", "directory to write tempates to").Short('d').StringVar(&opts.Dir)

	return nil
}

// CmdExportTemplates will export templates to directory
func (jc *JiraCli) CmdExportTemplates(opts *ExportTemplatesOptions) error {
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
