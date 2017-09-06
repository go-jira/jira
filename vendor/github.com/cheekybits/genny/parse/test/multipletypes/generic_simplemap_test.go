package multipletypes

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

func TestCustomTypesMap(t *testing.T) {

	key1 := new(MyType1)
	key2 := new(MyType1)
	value1 := new(MyOtherType)
	m := make(MyType1MyOtherTypeMap)

	assert.Equal(t, m, m.Set(key1, value1))
	assert.True(t, m.Has(key1))
	assert.False(t, m.Has(key2))
	assert.Equal(t, value1, m.Get(key1))

}
