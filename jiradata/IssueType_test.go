package jiradata

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIssueTypeFields(t *testing.T) {
	// this is because schema is wrong, missing the 'Fields' arguments, so we manually add it.
	// If the jiradata is regenerated we need to manually make the change again to include:
	// Fields      FieldMetaMap `json:"fields,omitempty" yaml:"fields,omitempty"`
	assert.IsType(t, FieldMetaMap{}, IssueType{}.Fields)
}
