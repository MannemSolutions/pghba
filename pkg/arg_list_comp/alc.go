package arg_list_comp

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	partsIsDone = fmt.Errorf("no more parts to split")
)

func parts (s string) (prefix string, comprehension string, postfix string, err error) {
	var exists bool
	re := regexp.MustCompile(`(?P<prefix>.*)(?P<comprehension>[[(][^])]*[])])(?P<postfix>.*)`)
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
	return
}

type ALC interface {
	Next()      (string, bool)
}

func NewALC (s string) (alc ALC, err error){
	prefix, comprehension, suffix, err := parts(s)
	if err != nil {
		return nil, err
	}
	if strings.HasPrefix(comprehension, "[") {
		return NewAlcCharList(prefix, comprehension, suffix)
	}
	if strings.Contains(comprehension, "..") {
		return NewAlcLoop(prefix, comprehension, suffix)
	}
	return NewAlcArray(prefix, comprehension, suffix)
}