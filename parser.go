package lre

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
	"unicode"
)

const (
	commentSharp = "#"
	commentSlash = "//"

	tagModel         = "Model"
	tagRuleStart     = "Rule"
	tagRuleIndex     = "Index"
	tagRuleCondition = "When"
	tagRuleAction    = "Then"
	tagRuleEnd       = "End"
)

type Parser struct {
}

/*
Model "data model full path"

Rule "rule name"
	Index 0
	When
		condition
	Then
		action 1
		action 2
		...
		action N
	End
*/
func (s *Parser) Parse(reader io.Reader, groups *Groups) error {
	if reader == nil || groups == nil {
		return fmt.Errorf("parameter invalid: nil")
	}

	var group *Group = nil
	var rule *Rule = nil
	var ruleCondition *strings.Builder = nil
	var ruleActions []string = nil
	bufReader := bufio.NewReader(reader)
	for {
		line, err := s.readLine(bufReader)
		if err == io.EOF {
			if len(line) <= 0 {
				break
			}
		}
		if len(line) <= 0 {
			continue
		}

		if strings.HasPrefix(line, tagModel) {
			modelName := s.getTagValue(line, tagModel)
			group = groups.Add(modelName)
			continue
		}
		if group == nil {
			continue
		}

		if strings.HasPrefix(line, tagRuleStart) {
			ruleName := s.getTagValue(line, tagRuleStart)
			rule = &Rule{
				Name:    ruleName,
				Actions: make([]string, 0),
			}
			continue
		}
		if rule == nil {
			continue
		}
		if strings.HasPrefix(line, tagRuleEnd) {
			group.Rules = append(group.Rules, rule)
			if ruleActions != nil {
				rule.Actions = ruleActions
			}
			rule = nil
			ruleCondition = nil
			ruleActions = nil
			continue
		}

		// rule content
		if strings.HasPrefix(line, tagRuleIndex) {
			ruleIndex := s.getTagValue(line, tagRuleIndex)
			index, err := strconv.Atoi(ruleIndex)
			if err == nil {
				rule.Index = index
			}
			continue
		}

		if strings.HasPrefix(line, tagRuleCondition) {
			ruleCondition = &strings.Builder{}
			continue
		}
		if strings.HasPrefix(line, tagRuleAction) {
			if ruleCondition != nil {
				rule.Condition = ruleCondition.String()
				ruleCondition = nil
			}
			ruleActions = make([]string, 0)
			continue
		}
		if ruleCondition != nil {
			ruleCondition.WriteString(line)
			continue
		}
		if ruleActions != nil {
			ruleActions = append(ruleActions, line)
			continue
		}
	}

	return nil
}

func (s *Parser) readLine(reader *bufio.Reader) (string, error) {
	line, err := reader.ReadString('\n')
	if err != nil {
		return line, err
	}
	if len(line) < 1 {
		return line, err
	}

	trimLine := s.trimLeft(line)
	if strings.HasPrefix(trimLine, commentSharp) || strings.HasPrefix(trimLine, commentSlash) {
		return "", nil
	}

	return s.trimRight(trimLine), nil
}

func (s *Parser) trimLeft(v string) string {
	return strings.TrimLeftFunc(v, unicode.IsSpace)
}

func (s *Parser) trimRight(v string) string {
	return strings.TrimRightFunc(v, func(r rune) bool {
		if uint32(r) <= unicode.MaxLatin1 {
			switch r {
			case '\n', '\r':
				return true
			}
			return false
		}

		return false
	})
}

func (s *Parser) getTagValue(line, tag string) string {
	val := strings.TrimLeft(line, tag)
	val = strings.TrimSpace(val)
	val = strings.TrimLeft(val, "\"")
	val = strings.TrimRight(val, "\"")

	return val
}
