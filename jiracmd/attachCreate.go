package jiracmd

import (
	"fmt"
	"os"
	"sort"

	"golang.org/x/crypto/ssh/terminal"

	"github.com/coryb/figtree"
	"github.com/eroshan/oreo"
	jira "github.com/go-jira/jira"
	"github.com/go-jira/jira/jiracli"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
	yaml "gopkg.in/coryb/yaml.v2"
)

type AttachCreateOptions struct {
	jiracli.CommonOptions `yaml:",inline" json:",inline" figtree:",inline"`
	Project               string `yaml:"project,omitempty" json:"project,omitempty"`
	Issue                 string `yaml:"issue,omitempty" json:"issue,omitempty"`
	Attachment            string `yaml:"attachment,omitempty" json:"attachment,omitempty"`
	Filename              string `yaml:"filename,omitempty" json:"filename,omitempty"`
	SaveFile              string `yaml:"savefile,omitempty" json:"savefile,omitempty"`
}

func CmdAttachCreateRegistry() *jiracli.CommandRegistryEntry {
	opts := AttachCreateOptions{}

	return &jiracli.CommandRegistryEntry{
		"Attach file to issue",
		func(fig *figtree.FigTree, cmd *kingpin.CmdClause) error {
			jiracli.LoadConfigs(cmd, fig, &opts)
			return CmdAttachCreateUsage(cmd, &opts)
		},
		func(o *oreo.Client, globals *jiracli.GlobalOptions) error {
			opts.Issue = jiracli.FormatIssue(opts.Issue, opts.Project)
			return CmdAttachCreate(o, globals, &opts)
		},
	}
}

func CmdAttachCreateUsage(cmd *kingpin.CmdClause, opts *AttachCreateOptions) error {
	jiracli.BrowseUsage(cmd, &opts.CommonOptions)
	cmd.Flag("saveFile", "Write attachment information as yaml to file").StringVar(&opts.SaveFile)
	cmd.Flag("filename", "Filename to use for attachment").Short('f').StringVar(&opts.Filename)
	cmd.Arg("ISSUE", "issue to assign").Required().StringVar(&opts.Issue)
	cmd.Arg("ATTACHMENT", "File to attach to issue, if not provided read from stdin").StringVar(&opts.Attachment)
	return nil
}

func CmdAttachCreate(o *oreo.Client, globals *jiracli.GlobalOptions, opts *AttachCreateOptions) error {
	var contents *os.File
	var err error
	if opts.Attachment == "" {
		if terminal.IsTerminal(int(os.Stdin.Fd())) {
			return fmt.Errorf("ATTACHMENT argument required or redirect from STDIN")
		}
		contents = os.Stdin
		if opts.Filename == "" {
			return fmt.Errorf("--filename required when reading from stdin")
		}
	} else {
		contents, err = os.Open(opts.Attachment)
		if err != nil {
			return err
		}
		if opts.Filename == "" {
			opts.Filename = opts.Attachment
		}
	}
	attachments, err := jira.IssueAttachFile(o, globals.Endpoint.Value, opts.Issue, opts.Filename, contents)
	if err != nil {
		return err
	}

	sort.Sort(sort.Reverse(attachments))

	if opts.SaveFile != "" {
		fh, err := os.Create(opts.SaveFile)
		if err != nil {
			return err
		}
		defer fh.Close()
		out, err := yaml.Marshal((*attachments)[0])
		if err != nil {
			return err
		}
		fh.Write(out)
	}

	if !globals.Quiet.Value {
		fmt.Printf("OK %d %s\n", (*attachments)[0].ID, (*attachments)[0].Content)
	}

	if opts.Browse.Value {
		return CmdBrowse(globals, opts.Issue)
	}

	return nil
}
