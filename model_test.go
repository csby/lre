package lre

import (
	"testing"
)

func TestModel_GetPath(t *testing.T) {
	model := &Model{}
	path := model.GetPath(model)
	t.Log("path:", path)
	if path != "github.com/csby/lre/Model" {
		t.Error(path)
	}
	path = model.GetPath(Model{})
	if path != "github.com/csby/lre/Model" {
		t.Error(path)
	}

	path = model.GetPath(t)
	if path != "testing/T" {
		t.Error(path)
	}

	path = model.GetPath(1)
	if path != "int" {
		t.Error(path)
	}
	path = model.GetPath(uint64(1))
	if path != "uint64" {
		t.Error(path)
	}
	path = model.GetPath("")
	if path != "string" {
		t.Error(path)
	}

	path = model.GetPath(nil)
	if path != "nil" {
		t.Error(path)
	}
}
