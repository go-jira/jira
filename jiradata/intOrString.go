package jiradata

import (
	"encoding/json"
	"strconv"
)

// this is for some bad schemas like Attachments.ID where in some api's it is an `int` and some it is a `string`
type IntOrString int

func (i *IntOrString) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var tmp string
	if err := unmarshal(&tmp); err != nil {
		return unmarshal((*int)(i))
	}
	tmpInt, err := strconv.Atoi(tmp)
	*i = IntOrString(tmpInt)
	return err
}

func (i *IntOrString) UnmarshalJSON(b []byte) error {
	var tmp string
	if err := json.Unmarshal(b, &tmp); err != nil {
		return json.Unmarshal(b, (*int)(i))
	}
	tmpInt, err := strconv.Atoi(tmp)
	*i = IntOrString(tmpInt)
	return err
}
