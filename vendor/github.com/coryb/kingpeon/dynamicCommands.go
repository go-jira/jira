package kingpeon

import (
	"fmt"
	"strings"
)

type DynamicCommandOpt struct {
	Name     string                  `yaml:"name,omitempty" json:"name,omitempty"`
	Type     DynamicCommandValueType `yaml:"type,omitempty" json:"type,omitempty"`
	Help     string                  `yaml:"help,omitempty" json:"help,omitempty"`
	Short    string                  `yaml:"short,omitempty" json:"short,omitempty"`
	Required bool                    `yaml:"required,omitempty" json:"required,omitempty"`
	Default  interface{}             `yaml:"default,omitempty" json:"default,omitempty"`
	Hidden   bool                    `yaml:"hidden,omitempty" json:"hidden,omitempty"`
	Repeat   bool                    `yaml:"repeat,omitempty" json:"repeat,omitempty"`
	Enum     []string                `yaml:"enum,omitempty" json:"enum,omitempty"`
}

type DynamicCommandArg struct {
	Name     string                  `yaml:"name,omitempty" json:"name,omitempty"`
	Help     string                  `yaml:"help,omitempty" json:"help,omitempty"`
	Type     DynamicCommandValueType `yaml:"type,omitempty" json:"type,omitempty"`
	Required bool                    `yaml:"required,omitempty" json:"required,omitempty"`
	Default  interface{}             `yaml:"default,omitempty" json:"default,omitempty"`
	Repeat   bool                    `yaml:"repeat,omitempty" json:"repeat,omitempty"`
	Enum     []string                `yaml:"enum,omitempty" json:"enum,omitempty"`
}

type DynamicCommand struct {
	Name    string              `yaml:"name,omitempty" json:"name,omitempty"`
	Options []DynamicCommandOpt `yaml:"options,omitempty" json:"options,omitempty"`
	Args    []DynamicCommandArg `yaml:"args,omitempty" json:"args,omitempty"`
	Script  string              `yaml:"script,omitempty" json:"script,omitempty"`
	Help    string              `yaml:"help,omitempty" json:"help,omitempty"`
	Default bool                `yaml:"default,omitempty" json:"default,omitempty"`
	Hidden  bool                `yaml:"hidden,omitempty" json:"hidden,omitempty"`
	Aliases []string            `yaml:"aliases,omitempty" json:"aliases,omitempty"`
}

type DynamicCommands []DynamicCommand

type DynamicCommandValueType int

const (
	DEFAULT DynamicCommandValueType = iota
	BOOL
	COUNTER
	ENUM
	FLOAT32
	FLOAT64
	INT8
	INT16
	INT32
	INT64
	INT
	STRING
	STRINGMAP
	UINT8
	UINT16
	UINT32
	UINT64
	UINT
)

func (o *DynamicCommandValueType) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var optType string
	if err := unmarshal(&optType); err != nil {
		return err
	}
	switch strings.ToUpper(optType) {
	case "BOOL":
		*o = BOOL
	case "COUNTER":
		*o = COUNTER
	case "ENUM":
		*o = ENUM
	case "FLOAT32":
		*o = FLOAT32
	case "FLOAT64":
		*o = FLOAT64
	case "INT8":
		*o = INT8
	case "INT16":
		*o = INT16
	case "INT32":
		*o = INT32
	case "INT64":
		*o = INT64
	case "INT":
		*o = INT
	case "STRING":
		*o = STRING
	case "STRINGMAP":
		*o = STRINGMAP
	case "UINT8":
		*o = UINT8
	case "UINT16":
		*o = UINT16
	case "UINT32":
		*o = UINT32
	case "UINT64":
		*o = UINT64
	case "UINT":
		*o = UINT
	default:
		return fmt.Errorf("Unknown option type: %s", optType)
	}
	return nil
}
