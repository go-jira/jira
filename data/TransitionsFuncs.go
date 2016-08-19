package jiradata

import (
	"strings"
)

func (t Transitions) Find(name string) *Transition {
	name = strings.ToLower(name)
	for _, trans := range t {
		if strings.Contains(strings.ToLower(trans.Name), name) {
			return trans
		}
	}
	return nil
}
