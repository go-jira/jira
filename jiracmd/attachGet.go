package jiracmd

import (
	"fmt"
	"io"
	"os"

	"github.com/coryb/figtree"
	"github.com/eroshan/oreo"
	jira "github.com/go-jira/jira"
	"github.com/go-jira/jira/jiracli"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type AttachGetOptions struct {
	AttachmentID string `yaml:"attachment-id,omitempty" json:"attachment-id,omitempty"`
	OutputFile   string `yaml:"output-file,omitempty" json:"output-file,omitempty"`
}

func CmdAttachGetRegistry() *jiracli.CommandRegistryEntry {
	opts := AttachGetOptions{}

	return &jiracli.CommandRegistryEntry{
		"Fetch attachment",
		func(fig *figtree.FigTree, cmd *kingpin.CmdClause) error {
			jiracli.LoadConfigs(cmd, fig, &opts)
			return CmdAttachGetUsage(cmd, &opts)
		},
		func(o *oreo.Client, globals *jiracli.GlobalOptions) error {
			return CmdAttachGet(o, globals, &opts)
		},
	}
}

func CmdAttachGetUsage(cmd *kingpin.CmdClause, opts *AttachGetOptions) error {
	cmd.Flag("output", "Write attachment to specified file name, '-' for stdout").Short('o').StringVar(&opts.OutputFile)
	cmd.Arg("ATTACHMENT-ID", "Attachment id to fetch").StringVar(&opts.AttachmentID)
	return nil
}

func CmdAttachGet(o *oreo.Client, globals *jiracli.GlobalOptions, opts *AttachGetOptions) error {
	attachment, err := jira.GetAttachment(o, globals.Endpoint.Value, opts.AttachmentID)
	if err != nil {
		return err
	}

	resp, err := o.Get(attachment.Content)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var output *os.File
	if opts.OutputFile == "-" {
		output = os.Stdout
	} else if opts.OutputFile != "" {
		output, err = os.Create(opts.OutputFile)
		if err != nil {
			return err
		}
		defer output.Close()
	} else {
		output, err = os.Create(attachment.Filename)
		if err != nil {
			return err
		}
		defer output.Close()
	}

	_, err = io.Copy(output, resp.Body)
	if err != nil {
		return err
	}
	output.Close()
	if opts.OutputFile != "-" && !globals.Quiet.Value {
		fmt.Printf("OK Wrote %s\n", output.Name())
	}
	return nil
}
