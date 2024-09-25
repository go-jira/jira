package jiracli

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"syscall"

	"github.com/coryb/figtree"
	"github.com/coryb/kingpeon"
	"github.com/eroshan/oreo"
	jira "github.com/go-jira/jira"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var usage = `{{define "FormatCommand"}}\
{{if .FlagSummary}} {{.FlagSummary}}{{end}}\
{{range .Args}} {{if not .Required}}[{{end}}<{{.Name}}>{{if .Value|IsCumulative}}...{{end}}{{if not .Required}}]{{end}}{{end}}\
{{end}}\

{{define "FormatBriefCommands"}}\
{{range .FlattenedCommands}}\
{{if not .Hidden}}\
  {{ print .FullCommand ":" | printf "%-20s"}} {{.Help}}
{{end}}\
{{end}}\
{{end}}\

{{define "FormatCommands"}}\
{{range .FlattenedCommands}}\
{{if not .Hidden}}\
  {{.FullCommand}}{{if .Default}}*{{end}}{{template "FormatCommand" .}}
{{.Help|Wrap 4}}
{{with .Flags|FlagsToTwoColumns}}{{FormatTwoColumnsWithIndent . 4 2}}{{end}}
{{end}}\
{{end}}\
{{end}}\

{{define "FormatUsage"}}\
{{template "FormatCommand" .}}{{if .Commands}} <command> [<args> ...]{{end}}
{{if .Help}}
{{.Help|Wrap 0}}\
{{end}}\

{{end}}\

{{if .Context.SelectedCommand}}\
usage: {{.App.Name}} {{.Context.SelectedCommand}}{{template "FormatCommand" .Context.SelectedCommand}}
{{if .Context.SelectedCommand.Aliases }}\
{{range $top := .App.Commands}}\
{{if eq $top.FullCommand $.Context.SelectedCommand.FullCommand}}\
{{range $alias := $.Context.SelectedCommand.Aliases}}\
alias: {{$.App.Name}} {{$alias}}{{template "FormatCommand" $.Context.SelectedCommand}}
{{end}}\
{{else}}\
{{range $sub := $top.Commands}}\
{{if eq $sub.FullCommand $.Context.SelectedCommand.FullCommand}}\
{{range $alias := $.Context.SelectedCommand.Aliases}}\
alias: {{$.App.Name}} {{$top.Name}} {{$alias}}{{template "FormatCommand" $.Context.SelectedCommand}}
{{end}}\
{{end}}\
{{end}}\
{{end}}\
{{end}}\
{{end}}
{{if .Context.SelectedCommand.Help}}\
{{.Context.SelectedCommand.Help|Wrap 0}}
{{end}}\
{{else}}\
usage: {{.App.Name}}{{template "FormatUsage" .App}}
{{end}}\

{{if .App.Flags}}\
Global flags:
{{.App.Flags|FlagsToTwoColumns|FormatTwoColumns}}
{{end}}\
{{if .Context.SelectedCommand}}\
{{if and .Context.SelectedCommand.Flags|RequiredFlags}}\
Required flags:
{{.Context.SelectedCommand.Flags|RequiredFlags|FlagsToTwoColumns|FormatTwoColumns}}
{{end}}\
{{if .Context.SelectedCommand.Flags|OptionalFlags}}\
Optional flags:
{{.Context.SelectedCommand.Flags|OptionalFlags|FlagsToTwoColumns|FormatTwoColumns}}
{{end}}\
{{end}}\
{{if .Context.Args}}\
Args:
{{.Context.Args|ArgsToTwoColumns|FormatTwoColumns}}
{{end}}\
{{if .Context.SelectedCommand}}\
{{if .Context.SelectedCommand.Commands}}\
Subcommands:
{{template "FormatCommands" .Context.SelectedCommand}}
{{end}}\
{{else if .App.Commands}}\
Commands:
{{template "FormatBriefCommands" .App}}
{{end}}\
`

func CommandLine(fig *figtree.FigTree, o *oreo.Client) *kingpin.Application {
	app := kingpin.New("jira", "Jira Command Line Interface")
	app.HelpFlag.Short('h')
	app.UsageWriter(os.Stdout)
	app.ErrorWriter(os.Stderr)
	app.Command("version", "Prints version").PreAction(func(*kingpin.ParseContext) error {
		fmt.Println(jira.VERSION)
		panic(Exit{Code: 0})
	})
	app.UsageTemplate(usage)

	var verbosity int
	app.Flag("verbose", "Increase verbosity for debugging").Short('v').PreAction(func(_ *kingpin.ParseContext) error {
		os.Setenv("JIRA_DEBUG", fmt.Sprintf("%d", verbosity))
		IncreaseLogLevel(1)
		return nil
	}).CounterVar(&verbosity)

	app.Terminate(func(status int) {
		for _, arg := range os.Args {
			if arg == "-h" || arg == "--help" || len(os.Args) == 1 {
				panic(Exit{Code: 0})
			}
		}
		panic(Exit{Code: 1})
	})

	register(app, o, fig)

	// register custom commands
	data := struct {
		CustomCommands kingpeon.DynamicCommands `yaml:"custom-commands" json:"custom-commands"`
	}{}

	if err := fig.LoadAllConfigs("config.yml", &data); err != nil {
		log.Errorf("%s", err)
		panic(Exit{Code: 1})
	}

	if len(data.CustomCommands) > 0 {
		runner := syscall.Exec
		if runtime.GOOS == "windows" {
			runner = func(binary string, cmd []string, env []string) error {
				command := exec.Command(binary, cmd[1:]...)
				command.Stdin = os.Stdin
				command.Stdout = os.Stdout
				command.Stderr = os.Stderr
				command.Env = env
				return command.Run()
			}
		}

		tmp := map[string]interface{}{}
		fig.LoadAllConfigs("config.yml", &tmp)
		kingpeon.RegisterDynamicCommandsWithRunner(runner, app, data.CustomCommands, TemplateProcessor())
	}

	return app
}

func ParseCommandLine(app *kingpin.Application, args []string) {
	// checking for default usage of `jira ISSUE-123` but need to allow
	// for global options first like: `jira --user mothra ISSUE-123`
	ctx, err := app.ParseContext(args)
	if err != nil && ctx == nil {
		// This is an internal kingpin usage error, duplicate options/commands
		log.Fatalf("error: %s, ctx: %v", err, ctx)
	}

	if ctx != nil {
		if ctx.SelectedCommand == nil {
			next := ctx.Next()
			if next != nil {
				if ok, err := regexp.MatchString("^([A-Z]+-)?[0-9]+$", next.Value); err != nil {
					log.Errorf("Invalid Regex: %s", err)
				} else if ok {
					// insert "view" at i=1 (2nd position)
					os.Args = append(os.Args[:1], append([]string{"view"}, os.Args[1:]...)...)
				}
			}
		}
	}

	if _, err := app.Parse(os.Args[1:]); err != nil {
		if _, ok := err.(*Error); ok {
			log.Errorf("%s", err)
			panic(Exit{Code: 1})
		}
		ctx, _ := app.ParseContext(os.Args[1:])
		if ctx != nil {
			app.UsageForContext(ctx)
		}
		log.Errorf("Invalid Usage: %s", err)
		panic(Exit{Code: 1})
	}
}
