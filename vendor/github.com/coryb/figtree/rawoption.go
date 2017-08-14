//go:generate genny -in=$GOFILE -out=gen-$GOFILE gen "RawType=BUILTINS"

package figtree

import (
	"encoding/json"
	"fmt"

	"github.com/cheekybits/genny/generic"
)

type RawType generic.Type

type RawTypeOption struct {
	Source  string
	Defined bool
	Value   RawType
}

func NewRawTypeOption(dflt RawType) RawTypeOption {
	return RawTypeOption{
		Source:  "default",
		Defined: true,
		Value:   dflt,
	}
}

func (o RawTypeOption) IsDefined() bool {
	return o.Defined
}

func (o *RawTypeOption) SetSource(source string) {
	o.Source = source
}

func (o RawTypeOption) GetValue() interface{} {
	return o.Value
}

// This is useful with kingpin option parser
func (o *RawTypeOption) Set(s string) error {
	err := convertString(s, &o.Value)
	if err != nil {
		return err
	}
	o.Source = "override"
	o.Defined = true
	return nil
}

// This is useful with survey prompting library
func (o *RawTypeOption) WriteAnswer(name string, value interface{}) error {
	if v, ok := value.(RawType); ok {
		o.Value = v
		o.Defined = true
		o.Source = "prompt"
		return nil
	}
	return fmt.Errorf("Got %T expected %T type: %v", value, o.Value, value)
}

func (o *RawTypeOption) SetValue(v interface{}) error {
	if val, ok := v.(RawType); ok {
		o.Value = val
		o.Defined = true
		return nil
	}
	return fmt.Errorf("Got %T expected %T type: %v", v, o.Value, v)
}

func (o *RawTypeOption) UnmarshalYAML(unmarshal func(interface{}) error) error {
	if err := unmarshal(&o.Value); err != nil {
		return err
	}
	o.Defined = true
	return nil
}

func (o *RawTypeOption) UnmarshalJSON(b []byte) error {
	if err := json.Unmarshal(b, &o.Value); err != nil {
		return err
	}
	o.Defined = true
	return nil
}

func (o RawTypeOption) MarshalYAML() (interface{}, error) {
	if StringifyValue {
		return o.Value, nil
	}
	// need a copy of this struct without the MarshalYAML interface attached
	return struct {
		Value   RawType
		Source  string
		Defined bool
	}{
		Value:   o.Value,
		Source:  o.Source,
		Defined: o.Defined,
	}, nil
}

func (o RawTypeOption) MarshalJSON() ([]byte, error) {
	if StringifyValue {
		return json.Marshal(o.Value)
	}
	// need a copy of this struct without the MarshalJSON interface attached
	return json.Marshal(struct {
		Value   RawType
		Source  string
		Defined bool
	}{
		Value:   o.Value,
		Source:  o.Source,
		Defined: o.Defined,
	})
}

// String is required for kingpin to generate usage with this datatype
func (o RawTypeOption) String() string {
	if StringifyValue {
		return fmt.Sprintf("%v", o.Value)
	}
	return fmt.Sprintf("{Source:%s Defined:%t Value:%v}", o.Source, o.Defined, o.Value)
}

type MapRawTypeOption map[string]RawTypeOption

// Set is required for kingpin interfaces to allow command line params
// to be set to our map datatype
func (o *MapRawTypeOption) Set(value string) error {
	parts := stringMapRegex.Split(value, 2)
	if len(parts) != 2 {
		return fmt.Errorf("expected KEY=VALUE got '%s'", value)
	}
	val := RawTypeOption{}
	val.Set(parts[1])
	(*o)[parts[0]] = val
	return nil
}

// IsCumulative is required for kingpin interfaces to allow multiple values
// to be set on the data structure.
func (o MapRawTypeOption) IsCumulative() bool {
	return true
}

// String is required for kingpin to generate usage with this datatype
func (o MapRawTypeOption) String() string {
	return fmt.Sprintf("%v", map[string]RawTypeOption(o))
}

func (o MapRawTypeOption) Map() map[string]RawType {
	tmp := map[string]RawType{}
	for k, v := range o {
		tmp[k] = v.Value
	}
	return tmp
}

// This is useful with survey prompting library
func (o *MapRawTypeOption) WriteAnswer(name string, value interface{}) error {
	tmp := RawTypeOption{}
	if v, ok := value.(RawType); ok {
		tmp.Value = v
		tmp.Defined = true
		tmp.Source = "prompt"
		(*o)[name] = tmp
		return nil
	}
	return fmt.Errorf("Got %T expected %T type: %v", value, tmp.Value, value)
}

func (o MapRawTypeOption) IsDefined() bool {
	// true if the map has any keys
	if len(o) > 0 {
		return true
	}
	return false
}

type ListRawTypeOption []RawTypeOption

// Set is required for kingpin interfaces to allow command line params
// to be set to our map datatype
func (o *ListRawTypeOption) Set(value string) error {
	val := RawTypeOption{}
	val.Set(value)
	*o = append(*o, val)
	return nil
}

// This is useful with survey prompting library
func (o *ListRawTypeOption) WriteAnswer(name string, value interface{}) error {
	tmp := RawTypeOption{}
	if v, ok := value.(RawType); ok {
		tmp.Value = v
		tmp.Defined = true
		tmp.Source = "prompt"
		*o = append(*o, tmp)
		return nil
	}
	return fmt.Errorf("Got %T expected %T type: %v", value, tmp.Value, value)
}

// IsCumulative is required for kingpin interfaces to allow multiple values
// to be set on the data structure.
func (o ListRawTypeOption) IsCumulative() bool {
	return true
}

// String is required for kingpin to generate usage with this datatype
func (o ListRawTypeOption) String() string {
	return fmt.Sprintf("%v", []RawTypeOption(o))
}

func (o ListRawTypeOption) Slice() []RawType {
	tmp := []RawType{}
	for _, elem := range o {
		tmp = append(tmp, elem.Value)
	}
	return tmp
}

func (o ListRawTypeOption) IsDefined() bool {
	// true if the list is not empty
	if len(o) > 0 {
		return true
	}
	return false
}
