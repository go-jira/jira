package figtree

import "regexp"

type Option interface {
	IsDefined() bool
	GetValue() interface{}
	SetValue(interface{}) error
	SetSource(string)
}

var StringifyValue = true

// used in option parsing for map types Set routines
var stringMapRegex = regexp.MustCompile("[:=]")

// IsBoolFlag is required by kingpin interface to determine if
// this variable requires a value
func (b BoolOption) IsBoolFlag() bool {
	return true
}
