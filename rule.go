package lre

/*
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
type Rule struct {
	Name      string   `json:"name"`
	Index     int      `json:"index"`
	Condition string   `json:"condition"`
	Actions   []string `json:"actions"`
}

type Rules []*Rule

func (s Rules) Len() int {
	return len(s)
}

func (s Rules) Less(i, j int) bool {
	return s[i].Index > s[j].Index
}

func (s Rules) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
