package jiracli

import (
	"fmt"

	"github.com/coryb/figtree"
	"github.com/coryb/oreo"

	jira "gopkg.in/Netflix-Skunkworks/go-jira.v1"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type WatchAction int

const (
	WatcherAdd WatchAction = iota
	WatcherRemove
)

type WatchOptions struct {
	GlobalOptions `yaml:",inline" figtree:",inline"`
	Issue         string
	Watcher       string
	Action        WatchAction
}

func CmdWatchRegistry(fig *figtree.FigTree, o *oreo.Client) *CommandRegistryEntry {
	opts := WatchOptions{
		GlobalOptions: GlobalOptions{},
		Action:        WatcherAdd,
	}

	return &CommandRegistryEntry{
		"Add/Remove watcher to issue",
		func() error {
			return CmdWatch(o, &opts)
		},
		func(cmd *kingpin.CmdClause) error {
			LoadConfigs(cmd, fig, &opts)
			return CmdWatchUsage(cmd, &opts)
		},
	}
}

func CmdWatchUsage(cmd *kingpin.CmdClause, opts *WatchOptions) error {
	if err := GlobalUsage(cmd, &opts.GlobalOptions); err != nil {
		return err
	}
	BrowseUsage(cmd, &opts.GlobalOptions)
	cmd.Flag("remove", "remove watcher from issue").Short('r').PreAction(func(ctx *kingpin.ParseContext) error {
		opts.Action = WatcherRemove
		return nil
	}).Bool()
	cmd.Arg("ISSUE", "issue to add watcher").Required().StringVar(&opts.Issue)
	cmd.Arg("WATCHER", "username of watcher to add to issue").StringVar(&opts.Watcher)
	return nil
}

// CmdWatch will add the given watcher to the issue (or remove the watcher
// with the 'remove' flag)
func CmdWatch(o *oreo.Client, opts *WatchOptions) error {
	if opts.Watcher == "" {
		opts.Watcher = opts.User.Value
	}
	if opts.Action == WatcherAdd {
		if err := jira.IssueAddWatcher(o, opts.Endpoint.Value, opts.Issue, opts.Watcher); err != nil {
			return err
		}
	} else {
		if err := jira.IssueRemoveWatcher(o, opts.Endpoint.Value, opts.Issue, opts.Watcher); err != nil {
			return err
		}
	}

	fmt.Printf("OK %s %s/browse/%s\n", opts.Issue, opts.Endpoint.Value, opts.Issue)

	if opts.Browse.Value {
		return CmdBrowse(&BrowseOptions{opts.GlobalOptions, opts.Issue})
	}

	return nil
}
