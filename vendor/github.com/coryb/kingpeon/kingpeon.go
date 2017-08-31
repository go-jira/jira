package kingpeon

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"text/template"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

// type so we can mock out how scripts are executed for testing
// it defaults to syscall.Exec
type runner func(string, []string, []string) error

func runDynamicCommand(run runner, dynamiccommand *DynamicCommand, t *template.Template) error {
	buf := bytes.NewBufferString("")
	t, err := t.Parse(dynamiccommand.Script)
	if err != nil {
		return err
	}
	err = t.Execute(buf, nil)
	if err != nil {
		return err
	}

	bin, err := exec.LookPath("sh")
	if err != nil {
		return err
	}
	cmd := []string{"sh", "-c", buf.String()}
	return run(bin, cmd, os.Environ())
}

// either kingpin.Application or kingpin.CmdClause fit this interface
type kingpinAppOrCommand interface {
	Command(string, string) *kingpin.CmdClause
	GetCommand(string) *kingpin.CmdClause
}

func lookupCommand(app *kingpin.Application, command *DynamicCommand) *kingpin.CmdClause {
	commandWords := strings.Fields(command.Name)
	var appOrCmd kingpinAppOrCommand = app
	if len(commandWords) > 1 {
		for _, name := range commandWords[0 : len(commandWords)-1] {
			tmp := appOrCmd.GetCommand(name)
			if tmp == nil {
				tmp = appOrCmd.Command(name, "")
			}
			appOrCmd = tmp
		}
	}

	return appOrCmd.Command(commandWords[len(commandWords)-1], command.Help)
}

func RegisterDynamicCommands(app *kingpin.Application, commands DynamicCommands, t *template.Template) error {
	return doRegisterDynamicCommands(syscall.Exec, app, commands, t)
}

func doRegisterDynamicCommands(run runner, app *kingpin.Application, commands DynamicCommands, t *template.Template) error {
	args := map[string]interface{}{}
	opts := map[string]interface{}{}

	t = t.Funcs(map[string]interface{}{
		"args": func() map[string]interface{} {
			return args
		},
		"options": func() map[string]interface{} {
			return opts
		},
	})

	for _, command := range commands {
		cmd := lookupCommand(app, &command)
		for _, alt := range command.Aliases {
			cmd = cmd.Alias(alt)
		}

		if command.Default {
			cmd = cmd.Default()
		}

		if command.Hidden {
			cmd = cmd.Hidden()
		}

		for _, opt := range command.Options {
			cmdFlag := cmd.Flag(opt.Name, opt.Help)
			if opt.Short != "" {
				cmdFlag.Short(rune(opt.Short[0]))
			}
			if opt.Required {
				cmdFlag.Required()
			}

			if opt.Default != "" {
				opts[opt.Name] = opt.Default
			}

			if opt.Hidden {
				cmdFlag.Hidden()
			}

			switch opt.Type {
			case BOOL:
				if opt.Repeat {
					var val []bool
					cmdFlag.PreAction(func(_ *kingpin.ParseContext) error {
						opts[opt.Name] = val
						return nil
					}).BoolListVar(&val)
				} else {
					var val bool
					cmdFlag.PreAction(func(_ *kingpin.ParseContext) error {
						opts[opt.Name] = val
						return nil
					}).BoolVar(&val)
				}
			case COUNTER:
				if opt.Repeat {
					return fmt.Errorf("`type: COUNTER` and `repeat: true` not supported for %s", opt.Name)
				} else {
					var val int
					cmdFlag.PreAction(func(_ *kingpin.ParseContext) error {
						opts[opt.Name] = val
						return nil
					}).CounterVar(&val)
				}
			case ENUM:
				if opt.Repeat {
					var val []string
					cmdFlag.PreAction(func(_ *kingpin.ParseContext) error {
						opts[opt.Name] = val
						return nil
					}).EnumsVar(&val, opt.Enum...)
				} else {
					var val string
					cmdFlag.PreAction(func(_ *kingpin.ParseContext) error {
						opts[opt.Name] = val
						return nil
					}).EnumVar(&val, opt.Enum...)
				}
			case FLOAT32:
				if opt.Repeat {
					var val []float32
					cmdFlag.PreAction(func(_ *kingpin.ParseContext) error {
						opts[opt.Name] = val
						return nil
					}).Float32ListVar(&val)
				} else {
					var val float32
					cmdFlag.PreAction(func(_ *kingpin.ParseContext) error {
						opts[opt.Name] = val
						return nil
					}).Float32Var(&val)
				}
			case FLOAT64:
				if opt.Repeat {
					var val []float64
					cmdFlag.PreAction(func(_ *kingpin.ParseContext) error {
						opts[opt.Name] = val
						return nil
					}).Float64ListVar(&val)
				} else {
					var val float64
					cmdFlag.PreAction(func(_ *kingpin.ParseContext) error {
						opts[opt.Name] = val
						return nil
					}).Float64Var(&val)
				}
			case INT8:
				if opt.Repeat {
					var val []int8
					cmdFlag.PreAction(func(_ *kingpin.ParseContext) error {
						opts[opt.Name] = val
						return nil
					}).Int8ListVar(&val)
				} else {
					var val int8
					cmdFlag.PreAction(func(_ *kingpin.ParseContext) error {
						opts[opt.Name] = val
						return nil
					}).Int8Var(&val)
				}
			case INT16:
				if opt.Repeat {
					var val []int16
					cmdFlag.PreAction(func(_ *kingpin.ParseContext) error {
						opts[opt.Name] = val
						return nil
					}).Int16ListVar(&val)
				} else {
					var val int16
					cmdFlag.PreAction(func(_ *kingpin.ParseContext) error {
						opts[opt.Name] = val
						return nil
					}).Int16Var(&val)
				}
			case INT32:
				if opt.Repeat {
					var val []int32
					cmdFlag.PreAction(func(_ *kingpin.ParseContext) error {
						opts[opt.Name] = val
						return nil
					}).Int32ListVar(&val)
				} else {
					var val int32
					cmdFlag.PreAction(func(_ *kingpin.ParseContext) error {
						opts[opt.Name] = val
						return nil
					}).Int32Var(&val)
				}
			case INT64:
				if opt.Repeat {
					var val []int64
					cmdFlag.PreAction(func(_ *kingpin.ParseContext) error {
						opts[opt.Name] = val
						return nil
					}).Int64ListVar(&val)
				} else {
					var val int64
					cmdFlag.PreAction(func(_ *kingpin.ParseContext) error {
						opts[opt.Name] = val
						return nil
					}).Int64Var(&val)
				}
			case INT:
				if opt.Repeat {
					var val []int
					cmdFlag.PreAction(func(_ *kingpin.ParseContext) error {
						opts[opt.Name] = val
						return nil
					}).IntsVar(&val)
				} else {
					var val int
					cmdFlag.PreAction(func(_ *kingpin.ParseContext) error {
						opts[opt.Name] = val
						return nil
					}).IntVar(&val)
				}
			case DEFAULT, STRING:
				if opt.Repeat {
					var val []string
					cmdFlag.PreAction(func(_ *kingpin.ParseContext) error {
						opts[opt.Name] = val
						return nil
					}).StringsVar(&val)
				} else {
					var val string
					cmdFlag.PreAction(func(_ *kingpin.ParseContext) error {
						opts[opt.Name] = val
						return nil
					}).StringVar(&val)
				}
			case STRINGMAP:
				if opt.Repeat {
					return fmt.Errorf("`type: STRINGMAP` and `repeat: true` not supported for %s", opt.Name)
				} else {
					val := make(map[string]string)
					cmdFlag.PreAction(func(_ *kingpin.ParseContext) error {
						opts[opt.Name] = val
						return nil
					}).StringMapVar(&val)
				}
			case UINT8:
				if opt.Repeat {
					var val []uint8
					cmdFlag.PreAction(func(_ *kingpin.ParseContext) error {
						opts[opt.Name] = val
						return nil
					}).Uint8ListVar(&val)
				} else {
					var val uint8
					cmdFlag.PreAction(func(_ *kingpin.ParseContext) error {
						opts[opt.Name] = val
						return nil
					}).Uint8Var(&val)
				}
			case UINT16:
				if opt.Repeat {
					var val []uint16
					cmdFlag.PreAction(func(_ *kingpin.ParseContext) error {
						opts[opt.Name] = val
						return nil
					}).Uint16ListVar(&val)
				} else {
					var val uint16
					cmdFlag.PreAction(func(_ *kingpin.ParseContext) error {
						opts[opt.Name] = val
						return nil
					}).Uint16Var(&val)
				}
			case UINT32:
				if opt.Repeat {
					var val []uint32
					cmdFlag.PreAction(func(_ *kingpin.ParseContext) error {
						opts[opt.Name] = val
						return nil
					}).Uint32ListVar(&val)
				} else {
					var val uint32
					cmdFlag.PreAction(func(_ *kingpin.ParseContext) error {
						opts[opt.Name] = val
						return nil
					}).Uint32Var(&val)
				}
			case UINT64:
				if opt.Repeat {
					var val []uint64
					cmdFlag.PreAction(func(_ *kingpin.ParseContext) error {
						opts[opt.Name] = val
						return nil
					}).Uint64ListVar(&val)
				} else {
					var val uint64
					cmdFlag.PreAction(func(_ *kingpin.ParseContext) error {
						opts[opt.Name] = val
						return nil
					}).Uint64Var(&val)
				}
			case UINT:
				if opt.Repeat {
					var val []uint
					cmdFlag.PreAction(func(_ *kingpin.ParseContext) error {
						opts[opt.Name] = val
						return nil
					}).UintsVar(&val)
				} else {
					var val uint
					cmdFlag.PreAction(func(_ *kingpin.ParseContext) error {
						opts[opt.Name] = val
						return nil
					}).UintVar(&val)
				}
			}
		}

		for _, arg := range command.Args {
			cmdArg := cmd.Arg(arg.Name, arg.Help)
			if arg.Required {
				cmdArg.Required()
			}
			if arg.Default != "" {
				args[arg.Name] = arg.Default
			}

			switch arg.Type {
			case BOOL:
				if arg.Repeat {
					var val []bool
					cmdArg.PreAction(func(_ *kingpin.ParseContext) error {
						args[arg.Name] = val
						return nil
					}).BoolListVar(&val)
				} else {
					var val bool
					cmdArg.PreAction(func(_ *kingpin.ParseContext) error {
						args[arg.Name] = val
						return nil
					}).BoolVar(&val)
				}
			case COUNTER:
				if arg.Repeat {
					return fmt.Errorf("`type: COUNTER` and `repeat: true` not supported for %s", arg.Name)
				} else {
					var val int
					cmdArg.PreAction(func(_ *kingpin.ParseContext) error {
						args[arg.Name] = val
						return nil
					}).CounterVar(&val)
				}
			case ENUM:
				if arg.Repeat {
					var val []string
					cmdArg.PreAction(func(_ *kingpin.ParseContext) error {
						args[arg.Name] = val
						return nil
					}).EnumsVar(&val, arg.Enum...)
				} else {
					var val string
					cmdArg.PreAction(func(_ *kingpin.ParseContext) error {
						args[arg.Name] = val
						return nil
					}).EnumVar(&val, arg.Enum...)
				}
			case FLOAT32:
				if arg.Repeat {
					var val []float32
					cmdArg.PreAction(func(_ *kingpin.ParseContext) error {
						args[arg.Name] = val
						return nil
					}).Float32ListVar(&val)
				} else {
					var val float32
					cmdArg.PreAction(func(_ *kingpin.ParseContext) error {
						args[arg.Name] = val
						return nil
					}).Float32Var(&val)
				}
			case FLOAT64:
				if arg.Repeat {
					var val []float64
					cmdArg.PreAction(func(_ *kingpin.ParseContext) error {
						args[arg.Name] = val
						return nil
					}).Float64ListVar(&val)
				} else {
					var val float64
					cmdArg.PreAction(func(_ *kingpin.ParseContext) error {
						args[arg.Name] = val
						return nil
					}).Float64Var(&val)
				}
			case INT8:
				if arg.Repeat {
					var val []int8
					cmdArg.PreAction(func(_ *kingpin.ParseContext) error {
						args[arg.Name] = val
						return nil
					}).Int8ListVar(&val)
				} else {
					var val int8
					cmdArg.PreAction(func(_ *kingpin.ParseContext) error {
						args[arg.Name] = val
						return nil
					}).Int8Var(&val)
				}
			case INT16:
				if arg.Repeat {
					var val []int16
					cmdArg.PreAction(func(_ *kingpin.ParseContext) error {
						args[arg.Name] = val
						return nil
					}).Int16ListVar(&val)
				} else {
					var val int16
					cmdArg.PreAction(func(_ *kingpin.ParseContext) error {
						args[arg.Name] = val
						return nil
					}).Int16Var(&val)
				}
			case INT32:
				if arg.Repeat {
					var val []int32
					cmdArg.PreAction(func(_ *kingpin.ParseContext) error {
						args[arg.Name] = val
						return nil
					}).Int32ListVar(&val)
				} else {
					var val int32
					cmdArg.PreAction(func(_ *kingpin.ParseContext) error {
						args[arg.Name] = val
						return nil
					}).Int32Var(&val)
				}
			case INT64:
				if arg.Repeat {
					var val []int64
					cmdArg.PreAction(func(_ *kingpin.ParseContext) error {
						args[arg.Name] = val
						return nil
					}).Int64ListVar(&val)
				} else {
					var val int64
					cmdArg.PreAction(func(_ *kingpin.ParseContext) error {
						args[arg.Name] = val
						return nil
					}).Int64Var(&val)
				}
			case INT:
				if arg.Repeat {
					var val []int
					cmdArg.PreAction(func(_ *kingpin.ParseContext) error {
						args[arg.Name] = val
						return nil
					}).IntsVar(&val)
				} else {
					var val int
					cmdArg.PreAction(func(_ *kingpin.ParseContext) error {
						args[arg.Name] = val
						return nil
					}).IntVar(&val)
				}
			case DEFAULT, STRING:
				if arg.Repeat {
					var val []string
					cmdArg.PreAction(func(_ *kingpin.ParseContext) error {
						args[arg.Name] = val
						return nil
					}).StringsVar(&val)
				} else {
					var val string
					cmdArg.PreAction(func(_ *kingpin.ParseContext) error {
						args[arg.Name] = val
						return nil
					}).StringVar(&val)
				}
			case STRINGMAP:
				if arg.Repeat {
					return fmt.Errorf("`type: STRINGMAP` and `repeat: true` not supported for %s", arg.Name)
				} else {
					val := make(map[string]string)
					cmdArg.PreAction(func(_ *kingpin.ParseContext) error {
						args[arg.Name] = val
						return nil
					}).StringMapVar(&val)
				}
			case UINT8:
				if arg.Repeat {
					var val []uint8
					cmdArg.PreAction(func(_ *kingpin.ParseContext) error {
						args[arg.Name] = val
						return nil
					}).Uint8ListVar(&val)
				} else {
					var val uint8
					cmdArg.PreAction(func(_ *kingpin.ParseContext) error {
						args[arg.Name] = val
						return nil
					}).Uint8Var(&val)
				}
			case UINT16:
				if arg.Repeat {
					var val []uint16
					cmdArg.PreAction(func(_ *kingpin.ParseContext) error {
						args[arg.Name] = val
						return nil
					}).Uint16ListVar(&val)
				} else {
					var val uint16
					cmdArg.PreAction(func(_ *kingpin.ParseContext) error {
						args[arg.Name] = val
						return nil
					}).Uint16Var(&val)
				}
			case UINT32:
				if arg.Repeat {
					var val []uint32
					cmdArg.PreAction(func(_ *kingpin.ParseContext) error {
						args[arg.Name] = val
						return nil
					}).Uint32ListVar(&val)
				} else {
					var val uint32
					cmdArg.PreAction(func(_ *kingpin.ParseContext) error {
						args[arg.Name] = val
						return nil
					}).Uint32Var(&val)
				}
			case UINT64:
				if arg.Repeat {
					var val []uint64
					cmdArg.PreAction(func(_ *kingpin.ParseContext) error {
						args[arg.Name] = val
						return nil
					}).Uint64ListVar(&val)
				} else {
					var val uint64
					cmdArg.PreAction(func(_ *kingpin.ParseContext) error {
						args[arg.Name] = val
						return nil
					}).Uint64Var(&val)
				}
			case UINT:
				if arg.Repeat {
					var val []uint
					cmdArg.PreAction(func(_ *kingpin.ParseContext) error {
						args[arg.Name] = val
						return nil
					}).UintsVar(&val)
				} else {
					var val uint
					cmdArg.PreAction(func(_ *kingpin.ParseContext) error {
						args[arg.Name] = val
						return nil
					}).UintVar(&val)
				}
			}
		}
		copy := command
		cmd.Action(func(_ *kingpin.ParseContext) error {
			return runDynamicCommand(run, &copy, t)
		})
	}
	return nil
}
