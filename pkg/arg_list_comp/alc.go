package arg_list_comp

import (
	"fmt"
	"regexp"
	"strings"
)

type groupChar string

const (
	squareStart groupChar = "["
	roundStart groupChar = "("
	curlyStart groupChar = "{"
	squareEnd groupChar = "]"
	roundEnd groupChar = ")"
	curlyEnd groupChar = "}"
)

type groupChars map[groupChar]groupChar

func (g groupChars) allChars() []string {
	return append(g.allStartChars(), g.allEndChars()...)
}

func (g groupChars) allEndChars() (all []string) {
	for _, end := range g {
		all = append(all, string(end))
	}
	return all
}

func (g groupChars) allStartChars() (all []string) {
	for start := range g {
		all = append(all, string(start))
	}
	return all
}

var (
	partsIsDone = fmt.Errorf("no more parts to split")
	groupStartToEnd = groupChars{
		curlyStart: curlyEnd,
		roundStart: roundEnd,
		squareStart: squareEnd,
	}
)

func parts (s string, groupStart groupChar) (prefix string, comprehension string, postfix string, err error) {
	var exists bool
	groupEnd, exists := groupStartToEnd[groupStart]
	if ! exists {
		return "","", "", fmt.Errorf("invalid group start")
	}
	re := regexp.MustCompile(fmt.Sprintf(`(?P<prefix>.*)(?P<comprehension>\%s[^\%s]*\%s)(?P<postfix>.*)`,
		groupStart, strings.Join(groupStartToEnd.allChars(),"\\"), groupEnd))
	matches := re.FindStringSubmatch(s)
	if matches == nil {
		err = partsIsDone
		return
	}
	fields := make(map[string]string)
	for id, name := range re.SubexpNames() {
		fields[name] = matches[id]
	}
	prefix, exists = fields["prefix"]
	if ! exists {
		err = fmt.Errorf("there is no prefix")
		return
	}
	comprehension, exists = fields["comprehension"]
	if ! exists {
		err = fmt.Errorf("there is no comprehension part")
		return
	}
	postfix, exists = fields["postfix"]
	if ! exists {
		err = fmt.Errorf("there is no postfix")
		return
	}
	comprehension = comprehension[1:len(comprehension)-1]
	return
}

type ALC interface {
	Next()      (string, bool)
	ToArray()   (a array)
	Reset()
	String()    (s string)
}

func NewALC (s string) (alc ALC){
	if cl, err := newAlcCharList(s); err == nil {
		return cl
	}
	if a, err := newAlcArray(s); err == nil {
		return a
	}
	if l, err := newAlcLoop(s); err == nil {
		return l
	}
	if l, err := newAlcCharLoop(s); err == nil {
		return l
	}
	return nil
}