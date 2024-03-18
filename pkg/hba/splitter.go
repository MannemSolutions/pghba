package hba

import (
	"fmt"
	"regexp"
)

type splitter struct {
	rest string
}

func newSplitter(str string) splitter {
	return splitter{rest: str}
}

func (s *splitter) Next() (string, error) {
	re := regexp.MustCompile(`^\s*(?P<part>([^"\t ]+|"[^"\t ]*"))\s*(?P<rest>.*)$`)
	matches := re.FindStringSubmatch(s.rest)
	if matches == nil {
		return "", fmt.Errorf("could not find next part in %s", s.rest)
	}
	fields := make(map[string]string)
	for id, name := range re.SubexpNames() {
		fields[name] = matches[id]
	}
	part, exists := fields["part"]
	if !exists {
		return "", fmt.Errorf("next part seems empty in %s", s.rest)
	}
	s.rest = fields["rest"]
	return part, nil
}

func (s splitter) Rest() string {
	return s.rest
}
