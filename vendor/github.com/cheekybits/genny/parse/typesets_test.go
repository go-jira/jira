package parse_test

import (
	"testing"

	"github.com/cheekybits/genny/parse"
	"github.com/stretchr/testify/assert"
)

func TestArgsToTypeset(t *testing.T) {

	args := "Person=man,woman Animal=dog,cat Place=london,paris"
	ts, err := parse.TypeSet(args)

	if assert.NoError(t, err) {
		if assert.Equal(t, 8, len(ts)) {

			assert.Equal(t, ts[0]["Person"], "man")
			assert.Equal(t, ts[0]["Animal"], "dog")
			assert.Equal(t, ts[0]["Place"], "london")

			assert.Equal(t, ts[1]["Person"], "man")
			assert.Equal(t, ts[1]["Animal"], "dog")
			assert.Equal(t, ts[1]["Place"], "paris")

			assert.Equal(t, ts[2]["Person"], "man")
			assert.Equal(t, ts[2]["Animal"], "cat")
			assert.Equal(t, ts[2]["Place"], "london")

			assert.Equal(t, ts[3]["Person"], "man")
			assert.Equal(t, ts[3]["Animal"], "cat")
			assert.Equal(t, ts[3]["Place"], "paris")

			assert.Equal(t, ts[4]["Person"], "woman")
			assert.Equal(t, ts[4]["Animal"], "dog")
			assert.Equal(t, ts[4]["Place"], "london")

			assert.Equal(t, ts[5]["Person"], "woman")
			assert.Equal(t, ts[5]["Animal"], "dog")
			assert.Equal(t, ts[5]["Place"], "paris")

			assert.Equal(t, ts[6]["Person"], "woman")
			assert.Equal(t, ts[6]["Animal"], "cat")
			assert.Equal(t, ts[6]["Place"], "london")

			assert.Equal(t, ts[7]["Person"], "woman")
			assert.Equal(t, ts[7]["Animal"], "cat")
			assert.Equal(t, ts[7]["Place"], "paris")

		}
	}

	ts, err = parse.TypeSet("Person=man Animal=dog Place=london")
	if assert.NoError(t, err) {
		assert.Equal(t, 1, len(ts))
	}
	ts, err = parse.TypeSet("Person=1,2,3,4,5 Animal=1,2,3,4,5 Place=1,2,3,4,5")
	if assert.NoError(t, err) {
		assert.Equal(t, 125, len(ts))
	}
	ts, err = parse.TypeSet("Person=1 Animal=1,2,3,4,5 Place=1,2")
	if assert.NoError(t, err) {
		assert.Equal(t, 10, len(ts))
	}

	ts, err = parse.TypeSet("Person=interface{} Animal=interface{} Place=interface{}")
	if assert.NoError(t, err) {
		assert.Equal(t, 1, len(ts))
		assert.Equal(t, ts[0]["Animal"], "interface{}")
		assert.Equal(t, ts[0]["Person"], "interface{}")
		assert.Equal(t, ts[0]["Place"], "interface{}")
	}

}
