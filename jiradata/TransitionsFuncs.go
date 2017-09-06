package jiradata

import (
	"strings"
)

// Find will search the transitions for one that matches
// the given name.  It will return a valid trantion that matches
// or nil
func (t Transitions) Find(name string) *Transition {
	name = strings.ToLower(name)
	for _, trans := range t {
		if strings.Contains(strings.ToLower(trans.Name), name) {
			return trans
		}
	}
	return nil
}
