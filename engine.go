package lre

import (
	"fmt"
	"github.com/antonmedv/expr"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
)

type Engine interface {
	ParseFolder(folder string, ext ...string) error
	ParseFiles(files ...string) error
	Parse(readers ...io.Reader) error
	Execute(model interface{}) ([]*Rule, error)
}

func NewEngine() Engine {
	return &engine{
		groups: make(Groups),
	}
}

type engine struct {
	sync.RWMutex

	groups Groups
}

func (s *engine) ParseFolder(folder string, ext ...string) error {
	kv := make(map[string]string)
	s.getFiles(folder, kv, ext...)
	files := make([]string, 0)
	for k := range kv {
		files = append(files, k)
	}

	return s.ParseFiles(files...)
}

func (s *engine) ParseFiles(files ...string) error {
	s.Lock()
	defer s.Unlock()

	s.groups = make(Groups)
	defer s.groups.Sort()

	count := len(files)
	for i := 0; i < count; i++ {
		err := s.parseFile(files[0])
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *engine) Parse(readers ...io.Reader) error {
	s.Lock()
	defer s.Unlock()

	s.groups = make(Groups)
	defer s.groups.Sort()

	count := len(readers)
	for i := 0; i < count; i++ {
		err := s.parse(readers[0])
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *engine) Execute(model interface{}) ([]*Rule, error) {
	s.RLock()
	defer s.RUnlock()

	rules := make([]*Rule, 0)
	path := getModelPath(model)
	group, ok := s.groups[path]
	if ok {
		ruleCount := len(group.Rules)
		for ruleIndex := 0; ruleIndex < ruleCount; ruleIndex++ {
			rule := group.Rules[ruleIndex]
			result, err := expr.Eval(rule.Condition, model)
			if err != nil {
				return rules, fmt.Errorf("%s: %v", rule.Name, err)
			}
			if result == nil {
				continue
			}
			if !reflect.DeepEqual(result, true) {
				continue
			}

			rules = append(rules, rule)
			actionCount := len(rule.Actions)
			for actionIndex := 0; actionIndex < actionCount; actionIndex++ {
				_, err = expr.Eval(rule.Actions[actionIndex], model)
				if err != nil {
					return rules, fmt.Errorf("%s: %v", rule.Name, err)
				}
			}
		}
	}

	return rules, nil
}

func (s *engine) parseFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	return s.parse(file)
}

func (s *engine) parse(reader io.Reader) error {
	if s.groups == nil {
		s.groups = make(Groups)
	}

	parser := &Parser{}
	return parser.Parse(reader, &s.groups)
}

func (s *engine) getFiles(folder string, files map[string]string, ext ...string) {
	if files == nil {
		return
	}
	extLen := len(ext)
	infos, err := ioutil.ReadDir(folder)
	if err != nil {
		return
	}
	for _, info := range infos {
		path := filepath.Join(folder, info.Name())
		if info.IsDir() {
			s.getFiles(path, files, ext...)
		} else {
			include := false
			if extLen > 0 {
				for i := 0; i < extLen; i++ {
					if strings.HasSuffix(strings.ToLower(path), strings.ToLower(ext[i])) {
						include = true
						break
					}
				}
			} else {
				include = true
			}

			if include {
				_, ok := files[path]
				if !ok {
					files[path] = info.Name()
				}
			}

		}
	}
}
