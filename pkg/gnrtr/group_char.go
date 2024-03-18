package gnrtr

import (
	"fmt"
	"regexp"
	"strings"
)

type GroupChar string

const (
	squareStart GroupChar = "["
	roundStart  GroupChar = "("
	curlyStart  GroupChar = "{"
	squareEnd   GroupChar = "]"
	roundEnd    GroupChar = ")"
	curlyEnd    GroupChar = "}"
)

type GroupChars map[GroupChar]GroupChar

func (g GroupChars) allChars() []string {
	return append(g.allStartChars(), g.allEndChars()...)
}

func (g GroupChars) allEndChars() (all []string) {
	for _, end := range g {
		all = append(all, string(end))
	}
	return all
}

func (g GroupChars) allStartChars() (all []string) {
	for start := range g {
		all = append(all, string(start))
	}
	return all
}

var (
	groupStartToEnd = GroupChars{
		curlyStart:  curlyEnd,
		roundStart:  roundEnd,
		squareStart: squareEnd,
	}
)

func (groupStart GroupChar) Parts(s string) (prefix string, comprehension string, postfix string, err error) {
	var exists bool
	groupEnd, exists := groupStartToEnd[groupStart]
	if !exists {
		return "", "", "", fmt.Errorf("invalid group start")
	}
	re := regexp.MustCompile(fmt.Sprintf(`(?P<prefix>.*)(?P<comprehension>\%s[^\%s]*\%s)(?P<postfix>.*)`,
		groupStart, strings.Join(groupStartToEnd.allChars(), "\\"), groupEnd))
	matches := re.FindStringSubmatch(s)
	if matches == nil {
		err = fmt.Errorf("no more parts to split")
		return
	}
	fields := make(map[string]string)
	for id, name := range re.SubexpNames() {
		fields[name] = matches[id]
	}
	prefix, exists = fields["prefix"]
	if !exists {
		err = fmt.Errorf("there is no prefix")
		return
	}
	comprehension, exists = fields["comprehension"]
	if !exists {
		err = fmt.Errorf("there is no comprehension part")
		return
	}
	postfix, exists = fields["postfix"]
	if !exists {
		err = fmt.Errorf("there is no postfix")
		return
	}
	comprehension = comprehension[1 : len(comprehension)-1]
	return
}
