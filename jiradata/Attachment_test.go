package jiradata

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAttachmentID(t *testing.T) {
	// this is because schema is wrong, defaults to type `int`, so we manually change it
	// to `string`.  If the jiradata is regenerated we need to manually make the change
	// again to include:
	// ID         IntOrString            `json:"id,omitempty" yaml:"id,omitempty"`
	assert.IsType(t, IntOrString(0), Attachment{}.ID)
}
