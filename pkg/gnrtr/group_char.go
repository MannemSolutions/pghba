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
	return append(g.allStartChars(), g.allEndChars()...) // TODO why use a variadic argument here?
}

func (g GroupChars) allEndChars() (all []string) {
	for _, end := range g { // ask for key/value pairs, discard the keys.
		all = append(all, string(end))
	}
	return all
}

func (g GroupChars) allStartChars() (all []string) {
	for start := range g {
		// ask only for the keys. How does this work? (as in: how can you get away with not accepting the value into a var here.)
		// Answer: yes, you can. See https://bitfieldconsulting.com/golang/map-iteration
		all = append(all, string(start))
	}
	return all
}

// I fail to see why these aren't constants instead of vars?
// This maps group start characters to group end characters.
// And why don't you need to 'make' this map?
var (
	groupStartToEnd = GroupChars{
		curlyStart:  curlyEnd,
		roundStart:  roundEnd,
		squareStart: squareEnd,
	}
)

// This function walks over string 's' and figures out the groupings to cut it up correctly.
// The way this works is that it creates regular expressions on-the-fly using the somewhat
// difficult to read 'named capture' construct: (?P<label>re) to separate out the characters
// before the 'groupStart' character into <prefix>, everything from the 'groupStart' character
func (groupStart GroupChar) Parts(s string) (prefix string, comprehension string, postfix string, err error) {
	var exists bool
	groupEnd, exists := groupStartToEnd[groupStart] // groupStartToEnd is the variable defind above!
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
