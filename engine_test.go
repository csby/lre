package lre

import (
	"bytes"
	"path/filepath"
	"runtime"
	"testing"
)

func TestEngine_getFiles(t *testing.T) {
	folder := sourceFolder()
	t.Log("folder:", folder)
	files := make(map[string]string)
	eg := &engine{}
	eg.getFiles(folder, files, ".go")
	t.Log("file count:", len(files))
	if len(files) < 1 {
		t.Fatal(len(files))
	}
	for k, v := range files {
		t.Log(v, ":", k)
	}
}

func TestEngine_Execute(t *testing.T) {
	engine := NewEngine()
	rules, err := engine.Execute(nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(rules) > 0 {
		t.Fatal(len(rules))
	}

	reader := &bytes.Buffer{}
	_, err = reader.WriteString(`
Model "github.com/csby/lre/TestDataModel"

Rule "rule 1"
	When
		A > 1 && A < 10
	Then
		SetB(11)
		SetC("r1")
	End


Rule "rule 2"
	Index 20
	When
		A >= 10 && A < 20 && X != nil
	Then
		SetB(21)
		SetC("r2")
		X.SetB(22)
		Y.SetA(23)
	End
`)
	if err != nil {
		t.Fatal(err)
	}
	err = engine.Parse(reader)
	if err != nil {
		t.Fatal(err)
	}

	dm := &TestDataModel{X: &TestDataModel{}}
	rules, err = engine.Execute(dm)
	if err != nil {
		t.Fatal(err)
	}
	if len(rules) > 0 {
		t.Fatal(len(rules))
	}
	if dm.A != 0 {
		t.Fatal(dm.A)
	}
	if dm.B != 0 {
		t.Fatal(dm.B)
	}
	if dm.C != "" {
		t.Fatal(dm.C)
	}
	if dm.X.B != 0 {
		t.Fatal(dm.X.B)
	}
	if dm.Y.A != 0 {
		t.Fatal(dm.Y.A)
	}

	dm.A = 5
	rules, err = engine.Execute(dm)
	if err != nil {
		t.Fatal(err)
	}
	if len(rules) != 1 {
		t.Fatal(len(rules))
	}
	if rules[0].Name != "rule 1" {
		t.Fatal(rules[0].Name)
	}
	if dm.A != 5 {
		t.Fatal(dm.A)
	}
	if dm.B != 11 {
		t.Fatal(dm.B)
	}
	if dm.C != "r1" {
		t.Fatal(dm.C)
	}
	if dm.X.B != 0 {
		t.Fatal(dm.X.B)
	}
	if dm.Y.A != 0 {
		t.Fatal(dm.Y.A)
	}

	dm.A = 100
	rules, err = engine.Execute(dm)
	if err != nil {
		t.Fatal(err)
	}
	if len(rules) != 0 {
		t.Fatal(len(rules))
	}
	if dm.A != 100 {
		t.Fatal(dm.A)
	}
	if dm.B != 11 {
		t.Fatal(dm.B)
	}
	if dm.C != "r1" {
		t.Fatal(dm.C)
	}
	if dm.X.B != 0 {
		t.Fatal(dm.X.B)
	}
	if dm.Y.A != 0 {
		t.Fatal(dm.Y.A)
	}

	dm.A = 15
	rules, err = engine.Execute(dm)
	if err != nil {
		t.Fatal(err)
	}
	if len(rules) != 1 {
		t.Fatal(len(rules))
	}
	if rules[0].Name != "rule 2" {
		t.Fatal(rules[0].Name)
	}
	if dm.A != 15 {
		t.Fatal(dm.A)
	}
	if dm.B != 21 {
		t.Fatal(dm.B)
	}
	if dm.C != "r2" {
		t.Fatal(dm.C)
	}
	if dm.X.B != 22 {
		t.Fatal(dm.X.B)
	}
	if dm.Y.A != 0 {
		t.Fatal(dm.Y.A)
	}
}

func sourceFolder() string {
	_, file, _, _ := runtime.Caller(0)

	return filepath.Join(filepath.Dir(file))
}

type TestDataModel struct {
	A int
	B int
	C string

	X *TestDataModel
	Y TestSubDataModel
}

func (s *TestDataModel) SetB(v float64) {
	s.B = int(v)
}

func (s *TestDataModel) SetC(v string) {
	s.C = v
}

type TestSubDataModel struct {
	A int
}

func (s TestSubDataModel) SetA(v float64) {
	s.A = int(v)
}
