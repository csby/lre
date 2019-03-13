package lre

import (
	"bufio"
	"bytes"
	"io"
	"testing"
)

func TestParser_getTagValue(t *testing.T) {
	parser := &Parser{}
	val := parser.getTagValue("Model ", "Model")
	if val != "" {
		t.Fatal(val)
	}

	val = parser.getTagValue("Model \"github.com/csby/lre/Model\" ", "Model")
	if val != "github.com/csby/lre/Model" {
		t.Fatal(val)
	}

	val = parser.getTagValue("Rule \"rule 1\" ", "Rule")
	if val != "rule 1" {
		t.Fatal(val)
	}

	val = parser.getTagValue("Index 15 ", "Index")
	if val != "15" {
		t.Fatal(val)
	}
}

func TestParser_trimRight(t *testing.T) {
	parser := &Parser{}
	trim := parser.trimRight("")
	if trim != "" {
		t.Fatal(trim)
	}

	trim = parser.trimRight("\n")
	if trim != "" {
		t.Fatal(trim)
	}
	trim = parser.trimRight("\r")
	if trim != "" {
		t.Fatal(trim)
	}
	trim = parser.trimRight("\r\n")
	if trim != "" {
		t.Fatal(trim)
	}
	trim = parser.trimRight("\n\r")
	if trim != "" {
		t.Fatal(trim)
	}

	trim = parser.trimRight("	Test \r\n")
	if trim != "	Test " {
		t.Fatal(trim)
	}
}

func TestParser_trimLeft(t *testing.T) {
	parser := &Parser{}
	trim := parser.trimLeft("")
	if trim != "" {
		t.Fatal(trim)
	}

	trim = parser.trimLeft(" #")
	if trim != "#" {
		t.Fatal(trim)
	}
	trim = parser.trimLeft(" 	#")
	if trim != "#" {
		t.Fatal(trim)
	}
}

func TestParser_readLine(t *testing.T) {
	parser := &Parser{}
	buf := &bytes.Buffer{}
	line, err := parser.readLine(bufio.NewReader(buf))
	if err != io.EOF {
		t.Fatal(err)
	}
	if len(line) != 0 {
		t.Fatal(line)
	}

	buf.Reset()
	buf.WriteString("Rule \"rn\"")
	line, err = parser.readLine(bufio.NewReader(buf))
	if err != io.EOF {
		t.Fatal(err)
	}
	if line != "Rule \"rn\"" {
		t.Fatal(line)
	}
	line, err = parser.readLine(bufio.NewReader(buf))
	if err != io.EOF {
		t.Fatal(err)
	}
	if line != "" {
		t.Fatal(line)
	}

	buf.Reset()
	buf.WriteString(" Rule \"rn\" \r\n")
	line, err = parser.readLine(bufio.NewReader(buf))
	if err != nil {
		t.Fatal(err)
	}
	if line != "Rule \"rn\" " {
		t.Fatal(line)
	}
	line, err = parser.readLine(bufio.NewReader(buf))
	if err != io.EOF {
		t.Fatal(err)
	}
	if line != "" {
		t.Fatal(line)
	}

	buf.Reset()
	buf.WriteString("# Rule \"rn\" \r\n")
	line, err = parser.readLine(bufio.NewReader(buf))
	if err != nil {
		t.Fatal(err)
	}
	if line != "" {
		t.Fatal(line)
	}
	line, err = parser.readLine(bufio.NewReader(buf))
	if err != io.EOF {
		t.Fatal(err)
	}
	if line != "" {
		t.Fatal(line)
	}

	buf.Reset()
	buf.WriteString(" // Rule \"rn\" \r\n")
	line, err = parser.readLine(bufio.NewReader(buf))
	if err != nil {
		t.Fatal(err)
	}
	if line != "" {
		t.Fatal(line)
	}
	line, err = parser.readLine(bufio.NewReader(buf))
	if err != io.EOF {
		t.Fatal(err)
	}
	if line != "" {
		t.Fatal(line)
	}
}

func TestParser_Parse(t *testing.T) {
	parser := &Parser{}
	err := parser.Parse(nil, nil)
	if err == nil {
		t.Fatal(err)
	}

	groups := make(Groups)
	reader := &bytes.Buffer{}
	_, err = reader.WriteString(`
Model "Model1"

Rule "rule 1"
	When
		2 > 1 && len([1, 3, 7]) == 3
	Then
		action 1
		action 2
		...
		action N
	End


Rule "rule 2"
	Index 20
	When
		len("Test") > 0 
	Then
		 action 21	
		action 22
	End
`)
	if err != nil {
		t.Fatal(err)
	}
	err = parser.Parse(reader, &groups)
	if err != nil {
		t.Fatal(err)
	}
	if len(groups) != 1 {
		t.Fatal(len(groups))
	}
	group, ok := groups["Model1"]
	if !ok {
		t.Fatal(ok)
	}
	if group.Model != "Model1" {
		t.Fatal(group.Model)
	}
	if len(group.Rules) != 2 {
		t.Fatal(len(group.Rules))
	}
	rule := group.Rules[0]
	if rule.Index != 0 {
		t.Fatal(rule.Index)
	}
	if rule.Name != "rule 1" {
		t.Fatal(rule.Name)
	}
	if rule.Condition != "2 > 1 && len([1, 3, 7]) == 3" {
		t.Fatal(rule.Condition)
	}
	if len(rule.Actions) != 4 {
		t.Fatal(len(rule.Actions))
	}
	if rule.Actions[0] != "action 1" {
		t.Fatal(rule.Actions[0])
	}
	if rule.Actions[1] != "action 2" {
		t.Fatal(rule.Actions[1])
	}
	if rule.Actions[2] != "..." {
		t.Fatal(rule.Actions[2])
	}
	if rule.Actions[3] != "action N" {
		t.Fatal(rule.Actions[3])
	}

	rule = group.Rules[1]
	if rule.Index != 20 {
		t.Fatal(rule.Index)
	}
	if rule.Name != "rule 2" {
		t.Fatal(rule.Name)
	}
	if rule.Condition != "len(\"Test\") > 0 " {
		t.Fatal(rule.Condition)
	}
	if len(rule.Actions) != 2 {
		t.Fatal(len(rule.Actions))
	}
	if rule.Actions[0] != "action 21	" {
		t.Fatal(rule.Actions[0])
	}
	if rule.Actions[1] != "action 22" {
		t.Fatal(rule.Actions[1])
	}

	reader.Reset()
	_, err = reader.WriteString(`
Model "Model1"

Rule "rule 31"
	Index 31
	When
		x > 10
	Then
		action 31
	End
`)
	if err != nil {
		t.Fatal(err)
	}
	err = parser.Parse(reader, &groups)
	if err != nil {
		t.Fatal(err)
	}
	if len(groups) != 1 {
		t.Fatal(len(groups))
	}
	group, ok = groups["Model1"]
	if !ok {
		t.Fatal(ok)
	}
	if group.Model != "Model1" {
		t.Fatal(group.Model)
	}
	if len(group.Rules) != 3 {
		t.Fatal(len(group.Rules))
	}
	rule = group.Rules[2]
	if rule.Index != 31 {
		t.Fatal(rule.Index)
	}
	if rule.Name != "rule 31" {
		t.Fatal(rule.Name)
	}
	if rule.Condition != "x > 10" {
		t.Fatal(rule.Condition)
	}
	if len(rule.Actions) != 1 {
		t.Fatal(len(rule.Actions))
	}
	if rule.Actions[0] != "action 31" {
		t.Fatal(rule.Actions[0])
	}

	reader.Reset()
	_, err = reader.WriteString(`
Model "Model2"

Rule "rule 41"
	When
		x > 44
	Then
		action 45
	End
`)
	if err != nil {
		t.Fatal(err)
	}
	err = parser.Parse(reader, &groups)
	if err != nil {
		t.Fatal(err)
	}
	if len(groups) != 2 {
		t.Fatal(len(groups))
	}
	group, ok = groups["Model1"]
	if !ok {
		t.Fatal(ok)
	}
	if group.Model != "Model1" {
		t.Fatal(group.Model)
	}
	if len(group.Rules) != 3 {
		t.Fatal(len(group.Rules))
	}
	rule = group.Rules[2]
	if rule.Index != 31 {
		t.Fatal(rule.Index)
	}
	if rule.Name != "rule 31" {
		t.Fatal(rule.Name)
	}
	if rule.Condition != "x > 10" {
		t.Fatal(rule.Condition)
	}
	if len(rule.Actions) != 1 {
		t.Fatal(len(rule.Actions))
	}
	if rule.Actions[0] != "action 31" {
		t.Fatal(rule.Actions[0])
	}

	group, ok = groups["Model2"]
	if !ok {
		t.Fatal(ok)
	}
	if group.Model != "Model2" {
		t.Fatal(group.Model)
	}
	if len(group.Rules) != 1 {
		t.Fatal(len(group.Rules))
	}
	rule = group.Rules[0]
	if rule.Index != 0 {
		t.Fatal(rule.Index)
	}
	if rule.Name != "rule 41" {
		t.Fatal(rule.Name)
	}
	if rule.Condition != "x > 44" {
		t.Fatal(rule.Condition)
	}
	if len(rule.Actions) != 1 {
		t.Fatal(len(rule.Actions))
	}
	if rule.Actions[0] != "action 45" {
		t.Fatal(rule.Actions[0])
	}
}
