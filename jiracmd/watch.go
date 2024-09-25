package jiracmd

import (
	"fmt"
	"strings"

	"github.com/coryb/figtree"
	"github.com/eroshan/oreo"

	"github.com/go-jira/jira"
	"github.com/go-jira/jira/jiracli"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type WatchAction int

const (
	WatcherAdd WatchAction = iota
	WatcherRemove
)

type WatchOptions struct {
	jiracli.CommonOptions `yaml:",inline" json:",inline" figtree:",inline"`
	Project               string      `yaml:"project,omitempty" json:"project,omitempty"`
	Issue                 string      `yaml:"issue,omitempty" json:"issue,omitempty"`
	Watcher               string      `yaml:"watcher,omitempty" json:"watcher,omitempty"`
	Action                WatchAction `yaml:"-" json:"-"`
}

func CmdWatchRegistry() *jiracli.CommandRegistryEntry {
	opts := WatchOptions{
		CommonOptions: jiracli.CommonOptions{},
		Action:        WatcherAdd,
	}

	return &jiracli.CommandRegistryEntry{
		"Add/Remove watcher to issue",
		func(fig *figtree.FigTree, cmd *kingpin.CmdClause) error {
			jiracli.LoadConfigs(cmd, fig, &opts)
			return CmdWatchUsage(cmd, &opts)
		},
		func(o *oreo.Client, globals *jiracli.GlobalOptions) error {
			opts.Issue = jiracli.FormatIssue(opts.Issue, opts.Project)
			return CmdWatch(o, globals, &opts)
		},
	}
}

func CmdWatchUsage(cmd *kingpin.CmdClause, opts *WatchOptions) error {
	jiracli.BrowseUsage(cmd, &opts.CommonOptions)
	cmd.Flag("remove", "remove watcher from issue").Short('r').PreAction(func(ctx *kingpin.ParseContext) error {
		opts.Action = WatcherRemove
		return nil
	}).Bool()
	cmd.Arg("ISSUE", "issue to add watcher").Required().StringVar(&opts.Issue)
	cmd.Arg("WATCHER", "email or display name of watcher to add to issue").StringVar(&opts.Watcher)
	return nil
}

// CmdWatch will add the given watcher to the issue (or remove the watcher
// with the 'remove' flag)
func CmdWatch(o *oreo.Client, globals *jiracli.GlobalOptions, opts *WatchOptions) error {
	if opts.Watcher == "" {
		opts.Watcher = globals.Login.Value
	}

	if globals.JiraDeploymentType.Value == "" {
		serverInfo, err := jira.ServerInfo(o, globals.Endpoint.Value)
		if err != nil {
			return err
		}
		globals.JiraDeploymentType.Value = strings.ToLower(serverInfo.DeploymentType)
	}

	if globals.JiraDeploymentType.Value == jiracli.CloudDeploymentType {
		users, err := jira.UserSearch(o, globals.Endpoint.Value, &jira.UserSearchOptions{
			Query: opts.Watcher,
		})
		if err != nil {
			return err
		}
		if len(users) > 1 {
			return fmt.Errorf("Found %d accounts for users with username %q", len(users), opts.Watcher)
		} else if len(users) == 1 {
			opts.Watcher = users[0].AccountID
		}
	}

	if opts.Action == WatcherAdd {
		if err := jira.IssueAddWatcher(o, globals.Endpoint.Value, opts.Issue, opts.Watcher); err != nil {
			return err
		}
	} else {
		if err := jira.IssueRemoveWatcher(o, globals.Endpoint.Value, opts.Issue, opts.Watcher); err != nil {
			return err
		}
	}

	if !globals.Quiet.Value {
		fmt.Printf("OK %s %s\n", opts.Issue, jira.URLJoin(globals.Endpoint.Value, "browse", opts.Issue))
	}

	if opts.Browse.Value {
		return CmdBrowse(globals, opts.Issue)
	}

	return nil
}
