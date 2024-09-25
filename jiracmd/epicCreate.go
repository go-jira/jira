package jiracmd

import (
	"github.com/coryb/figtree"
	"github.com/eroshan/oreo"

	"github.com/go-jira/jira/jiracli"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func CmdEpicCreateRegistry() *jiracli.CommandRegistryEntry {
	opts := CreateOptions{
		CommonOptions: jiracli.CommonOptions{
			Template: figtree.NewStringOption("epic-create"),
		},
		Overrides: map[string]string{},
	}

	return &jiracli.CommandRegistryEntry{
		"Create Epic",
		func(fig *figtree.FigTree, cmd *kingpin.CmdClause) error {
			jiracli.LoadConfigs(cmd, fig, &opts)
			return CmdEpicCreateUsage(cmd, &opts)
		},
		func(o *oreo.Client, globals *jiracli.GlobalOptions) error {
			return CmdCreate(o, globals, &opts)
		},
	}
}

func CmdEpicCreateUsage(cmd *kingpin.CmdClause, opts *CreateOptions) error {
	jiracli.BrowseUsage(cmd, &opts.CommonOptions)
	jiracli.EditorUsage(cmd, &opts.CommonOptions)
	jiracli.TemplateUsage(cmd, &opts.CommonOptions)
	cmd.Flag("noedit", "Disable opening the editor").SetValue(&opts.SkipEditing)
	cmd.Flag("project", "project to create epic in").Short('p').StringVar(&opts.Project)
	cmd.Flag("epic-name", "Epic Name").Short('n').PreAction(func(ctx *kingpin.ParseContext) error {
		opts.Overrides["epic-name"] = jiracli.FlagValue(ctx, "epic-name")
		return nil
	}).String()
	cmd.Flag("comment", "Comment message for epic").Short('m').PreAction(func(ctx *kingpin.ParseContext) error {
		opts.Overrides["comment"] = jiracli.FlagValue(ctx, "comment")
		return nil
	}).String()
	cmd.Flag("override", "Set epic property").Short('o').StringMapVar(&opts.Overrides)
	cmd.Flag("saveFile", "Write epic as yaml to file").StringVar(&opts.SaveFile)
	return nil
}
