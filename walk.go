package recursivewalk

import "reflect"

type WalkFunc func(meta WalkMeta)

type WalkMeta struct {
	FieldPath string
	FieldName string
	Type      reflect.Type
	Value     interface{}
}

func Walk(x interface{}, walkFunc WalkFunc) error {
	r := NewRecursive()
	return r.Recursive(x, walkFunc)
}
