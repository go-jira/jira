package slice

import (
	"log"
	"reflect"

	"github.com/cheekybits/genny/generic"
)

type MyType generic.Type

func EnsureSlice(objectOrSlice interface{}) []MyType {
	log.Printf("%v", reflect.TypeOf(objectOrSlice))
	switch obj := objectOrSlice.(type) {
	case []MyType:
		log.Println("  returning it untouched")
		return obj
	case MyType, *MyType:
		log.Println("  wrapping in slice")
		return []MyType{obj}
	default:
		panic("ensure slice needs MyType or []MyType")
	}
}
