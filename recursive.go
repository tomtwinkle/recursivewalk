package recursivewalk

import (
	"fmt"
	"log"
	"reflect"
	"strings"
)

type Recursive interface {
	Recursive(
		x interface{},
		walkFunc WalkFunc,
	) error
}

type recursiveInfo struct {
	structPath  string
	structField reflect.StructField
	parentStack stackField
}

func NewRecursive() Recursive {
	return &recursiveInfo{}
}

func (r recursiveInfo) Recursive(
	x interface{},
	walkFunc WalkFunc,
) error {
	return r.recursive(x, walkFunc)
}

func (r recursiveInfo) recursive(
	x interface{},
	walkFunc WalkFunc,
) error {
	orgValue := reflect.ValueOf(x)
	reflectValue := reflect.Indirect(orgValue)

	if x == nil {
		return nil
	}

	switch reflectValue.Kind() {
	case reflect.Struct:
		return r.recursiveStruct(reflectValue, walkFunc)
	case reflect.Slice:
		return r.recursiveSlice(reflectValue, walkFunc)
	case reflect.Map:
		return r.recursiveMap(reflectValue, walkFunc)
	case reflect.Ptr:
		return r.recursive(reflectValue, walkFunc)
	default:
		if orgValue.CanInterface() {
			walkFunc(WalkMeta{
				FieldPath: fmt.Sprintf("%s.%s", r.structPath, r.structField.Name),
				FieldName: r.structField.Name,
				Type:      orgValue.Type(),
				Value:     orgValue.Interface(),
			})
		}
	}
	return nil
}

func (r recursiveInfo) recursiveStruct(
	reflectValue reflect.Value,
	walkFunc WalkFunc,
) error {
	defer r.clearMeta()
	reflectType := reflectValue.Type()
	for i := 0; i < reflectType.NumField(); i++ {
		structField := reflectType.Field(i)
		value := reflectValue.Field(i)
		if value.CanInterface() {
			r.setMeta(structField, value)
			if err := r.recursive(value.Interface(), walkFunc); err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *recursiveInfo) setMeta(filed reflect.StructField, value reflect.Value) {
	if s, err := r.parentStack.Peek(); err == nil {
		if r.structPath == "" {
			r.structPath = s.Name
		} else {
			depth := len(strings.Split(r.structPath, "."))
			if depth < r.parentStack.Size() {
				r.structPath = fmt.Sprintf("%s.%s", r.structPath, s.Name)
			}
		}
	}
	r.structField = filed
	if reflect.Indirect(value).Kind() == reflect.Struct {
		r.parentStack.Push(filed)
	}
}

func (r *recursiveInfo) clearMeta() {
	r.structPath = ""
	r.structField = reflect.StructField{}
	if _, err := r.parentStack.Pop(); err != nil {
		log.Fatal(err)
	}
}

func (r recursiveInfo) recursiveSlice(
	reflectValue reflect.Value,
	walkFunc WalkFunc,
) error {
	for j := 0; j < reflectValue.Len(); j++ {
		sliceStruct := reflectValue.Index(j)
		if sliceStruct.CanInterface() {
			if err := r.recursive(sliceStruct.Interface(), walkFunc); err != nil {
				return err
			}
		}
	}
	return nil
}

func (r recursiveInfo) recursiveMap(
	reflectValue reflect.Value,
	walkFunc WalkFunc,
) error {
	for _, key := range reflectValue.MapKeys() {
		mapStruct := reflectValue.MapIndex(key)
		if mapStruct.CanInterface() {
			if err := r.recursive(mapStruct.Interface(), walkFunc); err != nil {
				return err
			}
		}
	}
	return nil
}
