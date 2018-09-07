package jiracmd

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/coryb/figtree"
	"github.com/coryb/oreo"
	"gopkg.in/Netflix-Skunkworks/go-jira.v1"
	"gopkg.in/Netflix-Skunkworks/go-jira.v1/jiracli"
	"gopkg.in/Netflix-Skunkworks/go-jira.v1/jiradata"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
	"strings"
)

const (
	epicTitleStage         = "epicTitleStage"
	epicSummaryStage       = "epic_name"
	epicDescriptionStage   = "epic_description"
	ticketTitleStage       = "ticketTitleStage"
	ticketDescriptionStage = "ticketDescriptionStage"
)

type IssueResult struct {
	key       string
	link     string
}

type MarkdownOptions struct {
	MarkdownFile string `yaml:"mdfile,omitempty" json:"mdfile,omitempty"`
	jiracli.CommonOptions `yaml:",inline" json:",inline" figtree:",inline"`
	jiradata.IssueUpdate  `yaml:",inline" json:",inline" figtree:",inline"`
	Project               string            `yaml:"project,omitempty" json:"project,omitempty"`
	IssueType             string            `yaml:"issuetype,omitempty" json:"issuetype,omitempty"`
	Overrides             map[string]string `yaml:"overrides,omitempty" json:"overrides,omitempty"`
}

func CmdMarkdownParseRegistry() *jiracli.CommandRegistryEntry {
	opts := MarkdownOptions{}

	return &jiracli.CommandRegistryEntry{
		"Parse Markdown",
		func(fig *figtree.FigTree, cmd *kingpin.CmdClause) error {
			jiracli.LoadConfigs(cmd, fig, &opts)
			return CmdMdCreateUsage(cmd, &opts)
		},
		func(o *oreo.Client, globals *jiracli.GlobalOptions) error {
			return CmdMarkdownParse(o, globals, &opts)
		},
	}
}

func CmdMdCreateUsage(cmd *kingpin.CmdClause, opts *MarkdownOptions) error {
	jiracli.BrowseUsage(cmd, &opts.CommonOptions)
	jiracli.EditorUsage(cmd, &opts.CommonOptions)
	jiracli.TemplateUsage(cmd, &opts.CommonOptions)
	cmd.Flag("md", "Markdown file to parse").Short('f').StringVar(&opts.MarkdownFile)
	cmd.Flag("project", "project to create issue in").Short('p').StringVar(&opts.Project)
	return nil
}

func Map(vs []IssueResult, f func(IssueResult) string) []string {
	vsm := make([]string, len(vs))
	for i, v := range vs {
		vsm[i] = f(v)
	}
	return vsm
}

func CmdMarkdownParse(o *oreo.Client, globals *jiracli.GlobalOptions, opts *MarkdownOptions) error {
	plan, err := parseFile(opts.MarkdownFile)
	if err != nil {
		return err
	}
	for i := 0; i < len(plan.Epics); i++ {
		epic := plan.Epics[i]
		epicOpts := CreateOptions{
			CommonOptions: jiracli.CommonOptions{
				Template: figtree.NewStringOption("epic-create"),
			},
			//add overrides
			Overrides: map[string]string{
				"epic-name": strings.Replace(epic.Title, ":", "-", -1),
				"description": strings.Replace(epic.Description, ":", "-", -1),
				"summary": strings.Replace(epic.Summary, ":", "-", -1),
			},
			Project: opts.Project,
		}
		epicOpts.SkipEditing = figtree.NewBoolOption(true)

		epicResult, err := createWithResult(o, globals, &epicOpts)
		if err != nil {
			return err
		}
		tickets := []IssueResult{}
		for j := 0; j < len(epic.Tickets); j++ {

			ticket := epic.Tickets[j]
			ticketOpts := CreateOptions{
				CommonOptions: jiracli.CommonOptions{
					Template: figtree.NewStringOption("create"),
				},
				Overrides: map[string]string{
					"summary": strings.Replace(ticket.Title, ":", "-", -1),
					"description": strings.Replace(ticket.Description, ":", "-", -1),
				},
				IssueType: "Task",
				Project: opts.Project,
			}
			ticketOpts.SkipEditing = figtree.NewBoolOption(true)

			ticketResult, err := createWithResult(o, globals, &ticketOpts)
			if err != nil {
				return err
			}
			tickets = append(tickets, *ticketResult)
		}

		issues := Map(tickets, func(result IssueResult) string {
			return result.key
		})
		addOpts := EpicAddOptions{jiradata.EpicIssues{issues},epicResult.key}
		err = CmdEpicAdd(o, globals, &addOpts)
		if err != nil {
			return err
		}
		fmt.Printf("\nEpic Done %s \nTickets Added %+v", epicResult.key, issues)
	}

	return nil
}

func createWithResult(o *oreo.Client, globals *jiracli.GlobalOptions, opts *CreateOptions) (*IssueResult, error) {
	type templateInput struct {
		Meta      *jiradata.IssueType `yaml:"meta" json:"meta"`
		Overrides map[string]string   `yaml:"overrides" json:"overrides"`
	}

	if err := defaultIssueType(o, globals.Endpoint.Value, &opts.Project, &opts.IssueType); err != nil {
		return nil, err
	}
	createMeta, err := jira.GetIssueCreateMetaIssueType(o, globals.Endpoint.Value, opts.Project, opts.IssueType)
	if err != nil {
		return nil, err
	}

	issueUpdate := jiradata.IssueUpdate{}
	input := templateInput{
		Meta:      createMeta,
		Overrides: opts.Overrides,
	}
	input.Overrides["project"] = opts.Project
	input.Overrides["issuetype"] = opts.IssueType
	input.Overrides["user"] = globals.User.Value

	var issueResp *jiradata.IssueCreateResponse
	err = jiracli.EditLoop(&opts.CommonOptions, &input, &issueUpdate, func() error {
		issueResp, err = jira.CreateIssue(o, globals.Endpoint.Value, &issueUpdate)
		return err
	})
	if err != nil {
		return nil, err
	}

	browseLink := jira.URLJoin(globals.Endpoint.Value, "browse", issueResp.Key)

	return &IssueResult{issueResp.Key, browseLink}, nil
}

func parseFile(filePath string) (*jiradata.Plan, error) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	stage := epicTitleStage

	var plan jiradata.Plan
	var epic jiradata.Epic
	var epicPointer *jiradata.Epic = nil

	var epicTitle string
	var epicSummary string
	var epicDescription string

	var ticket jiradata.Ticket
	var ticketPointer *jiradata.Ticket = nil
	var ticketTitle string
	var ticketDescription string

	for scanner.Scan() {
		line := scanner.Text()

		if jiradata.IsEpicTitle(line) {
			if ticketPointer != nil {
				epic.AddTicket(ticket)
				ticketPointer = nil
			}

			if epicPointer != nil {
				plan.AddEpic(epic)
				epicPointer = nil
			}

			stage = epicTitleStage
			epicTitle = line[2:]
		} else if stage == epicTitleStage && jiradata.IsEpicSummary(line) {
			stage = epicSummaryStage
			epicSummary = line[2:len(line) - 2]
		} else if (stage == epicSummaryStage || stage == epicTitleStage) && jiradata.IsDescriptionLine(line) {
			stage = epicDescriptionStage
			epicDescription = line
		} else if stage == epicDescriptionStage && jiradata.IsDescriptionLine(line) {
			//sepator := " "
			//
			//if jiradata.IsItem(line) {
			//	sepator = "\n"
			//}

			epicDescription += "\n" + line
		} else if jiradata.IsTicketTitle(line) {
			if ticketPointer != nil {
				epic.AddTicket(ticket)
				ticketPointer = nil
			}

			stage = ticketTitleStage
			ticketTitle = line[3:]
		} else if stage == ticketTitleStage && jiradata.IsDescriptionLine(line) {
			stage = ticketDescriptionStage
			ticketDescription = line
		} else if stage == ticketDescriptionStage && jiradata.IsDescriptionLine(line) {
			//sepator := "\n"
			//
			//if jiradata.IsItem(line) {
			//	sepator = "\n"
			//}
			ticketDescription += "\n" + line
		} else if stage == epicDescriptionStage && jiradata.IsSeparator(line) {
			epic = jiradata.NewEpic(epicTitle, epicSummary, epicDescription)
			epicPointer = &epic
		} else if stage == ticketDescriptionStage && jiradata.IsSeparator(line) {
			ticket = jiradata.NewTicket(ticketTitle, ticketDescription)
			ticketPointer = &ticket
		} else if jiradata.IsSeparator(line) {
		} else {
			return nil, errors.New("Unexpected state")
		}
	}

	if ticketPointer != nil {
		epic.AddTicket(ticket)
		ticketPointer = nil
	}
	plan.AddEpic(epic)

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return &plan, nil
}