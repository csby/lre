package lre

import (
	"fmt"
	"reflect"
)

type Model struct {
}

func (s *Model) GetPath(v interface{}) string {
	return getModelPath(v)
}

func getModelPath(v interface{}) string {
	if v == nil {
		return "nil"
	}

	t := reflect.TypeOf(v)
	e := t
	if t.Kind() == reflect.Ptr {
		e = t.Elem()
	}

	p := e.PkgPath()
	if len(p) > 0 {
		return fmt.Sprintf("%s/%s", p, e.Name())
	} else {
		return e.Name()
	}
}
