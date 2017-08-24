package jiracli

import (
	"fmt"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type WatchAction int

const (
	WatcherAdd WatchAction = iota
	WatcherRemove
)

type WatchOptions struct {
	GlobalOptions
	Issue   string
	Watcher string
	Action  WatchAction
}

func (jc *JiraCli) CmdWatchRegistry() *CommandRegistryEntry {
	opts := WatchOptions{
		GlobalOptions: GlobalOptions{},
		Action:        WatcherAdd,
	}

	return &CommandRegistryEntry{
		"Add/Remove watcher to issue",
		func() error {
			return jc.CmdWatch(&opts)
		},
		func(cmd *kingpin.CmdClause) error {
			return jc.CmdWatchUsage(cmd, &opts)
		},
	}
}

func (jc *JiraCli) CmdWatchUsage(cmd *kingpin.CmdClause, opts *WatchOptions) error {
	if err := jc.GlobalUsage(cmd, &opts.GlobalOptions); err != nil {
		return err
	}
	jc.BrowseUsage(cmd, &opts.GlobalOptions)
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
func (jc *JiraCli) CmdWatch(opts *WatchOptions) error {
	if opts.Watcher == "" {
		opts.Watcher = opts.User.Value
	}
	if opts.Action == WatcherAdd {
		if err := jc.IssueAddWatcher(opts.Issue, opts.Watcher); err != nil {
			return err
		}
	} else {
		if err := jc.IssueRemoveWatcher(opts.Issue, opts.Watcher); err != nil {
			return err
		}
	}

	fmt.Printf("OK %s %s/browse/%s\n", opts.Issue, jc.Endpoint, opts.Issue)

	if opts.Browse.Value {
		return jc.CmdBrowse(&BrowseOptions{opts.GlobalOptions, opts.Issue})
	}

	return nil
}
