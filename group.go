package lre

import (
	"sort"
)

/*
Model "data model full path"

Rule "rule name"
	When
		condition
	Then
		action 1
		action 2
		...
		action N
	End
*/
type Group struct {
	Model string `json:"model"` // full path of data model, example: github.com/csby/lre/Model
	Rules Rules  `json:"rules"` // all rules of the data model
}

type Groups map[string]*Group

func (s Groups) Sort() {
	for _, v := range s {
		sort.Sort(v.Rules)
	}
}

func (s Groups) Add(model string) *Group {
	group, ok := s[model]
	if !ok {
		group = &Group{
			Model: model,
			Rules: make(Rules, 0),
		}
		s[model] = group
	}

	return group
}
