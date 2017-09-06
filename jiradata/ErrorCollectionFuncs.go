package jiradata

import "strings"

// Error is needed to make ErrorCollection implement the error interface
func (e ErrorCollection) Error() string {
	if len(e.ErrorMessages) > 0 {
		return strings.Join(e.ErrorMessages, ". ")
	}
	out := ""
	for k, v := range e.Errors {
		if len(out) > 0 {
			out += ". "
		}
		out += k + ": " + v
	}
	return out
}
