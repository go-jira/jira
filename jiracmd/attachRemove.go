package jiracmd

import (
	"fmt"

	"github.com/coryb/figtree"
	"github.com/eroshan/oreo"
	jira "github.com/go-jira/jira"
	"github.com/go-jira/jira/jiracli"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type AttachRemoveOptions struct {
	AttachmentID string `yaml:"attachment-id,omitempty" json:"attachment-id,omitempty"`
}

func CmdAttachRemoveRegistry() *jiracli.CommandRegistryEntry {
	opts := AttachRemoveOptions{}

	return &jiracli.CommandRegistryEntry{
		"Delete attachment",
		func(fig *figtree.FigTree, cmd *kingpin.CmdClause) error {
			jiracli.LoadConfigs(cmd, fig, &opts)
			return CmdAttachRemoveUsage(cmd, &opts)
		},
		func(o *oreo.Client, globals *jiracli.GlobalOptions) error {
			return CmdAttachRemove(o, globals, &opts)
		},
	}
}

func CmdAttachRemoveUsage(cmd *kingpin.CmdClause, opts *AttachRemoveOptions) error {
	cmd.Arg("ATTACHMENT-ID", "Attachment id to fetch").StringVar(&opts.AttachmentID)
	return nil
}

func CmdAttachRemove(o *oreo.Client, globals *jiracli.GlobalOptions, opts *AttachRemoveOptions) error {
	if err := jira.RemoveAttachment(o, globals.Endpoint.Value, opts.AttachmentID); err != nil {
		return err
	}

	if !globals.Quiet.Value {
		fmt.Printf("OK Deleted Attachment %s\n", opts.AttachmentID)
	}
	return nil
}
