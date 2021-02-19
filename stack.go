package recursivewalk

import (
	"fmt"
	"reflect"
)

type stackField []reflect.StructField

func (s *stackField) Push(v reflect.StructField) {
	*s = append(*s, v)
}

func (s *stackField) Pop() (reflect.StructField, error) {
	if s.Empty() {
		return reflect.StructField{}, fmt.Errorf("stack is empty")
	}
	v := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	return v, nil
}

func (s *stackField) Peek() (reflect.StructField, error) {
	if s.Empty() {
		return reflect.StructField{}, fmt.Errorf("stack is empty")
	}
	return (*s)[len(*s)-1], nil
}

func (s *stackField) Size() int {
	return len(*s)
}

func (s *stackField) Empty() bool {
	return s.Size() == 0
}
