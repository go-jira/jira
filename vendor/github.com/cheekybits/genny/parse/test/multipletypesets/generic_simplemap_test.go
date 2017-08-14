package multipletypesets

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimpleMap(t *testing.T) {

	key1 := new(KeyType)
	key2 := new(KeyType)
	value1 := new(ValueType)
	m := make(KeyTypeValueTypeMap)

	assert.Equal(t, m, m.Set(key1, value1))
	assert.True(t, m.Has(key1))
	assert.False(t, m.Has(key2))
	assert.Equal(t, value1, m.Get(key1))

}
