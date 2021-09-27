package arg_list_comp

import (
	"fmt"
	"regexp"
	"strings"
)

type charlist struct {
	prefix string
	list []byte
	suffix string
	index int
}

func NewAlcCharList(prefix string, comprehension string, suffix string) (cl charlist, err error) {
	if ! (strings.HasPrefix(comprehension, "[") && strings.HasSuffix(comprehension, "]")) {
		return charlist{}, fmt.Errorf("missing round braces around array comprehension")
	}
	cl = charlist{
		prefix: prefix,
		suffix: suffix,
		index: 0,
	}
	charDefinition := comprehension[1:len(comprehension)-1]
	if strings.HasSuffix(charDefinition, "-" ) {
		cl.list = []byte("-")
	}
	if strings.HasPrefix(charDefinition, "-" ) {
		cl.list = []byte("-")
	}
	re := regexp.MustCompile(`(-)?([^-]|([^-][-][^-]))*(-)?`)
	matches := re.FindAllIndex([]byte(charDefinition), -1)
	if matches == nil {
		return charlist{}, fmt.Errorf("could not parse charlist %s", charDefinition)
	}
	for _, match := range matches {
		start := charDefinition[match[0]]
		if match[1] - match[0] == 1 {
			cl.list = append(cl.list, start)
		} else if match[1] - match[0] == 3 {
			end := charDefinition[match[0]]
			for c:= start; c <= end; c++ {
				cl.list = append(cl.list, c)
			}
		} else {
			return charlist{}, fmt.Errorf("could not parse the %s part of %s", charDefinition[match[0]:match[1]],
				comprehension)
		}
	}
	return cl, nil
}

func (cl charlist) Next() (next string, done bool) {
	if cl.index >= len(cl.list) {
		return "", true
	}
	next = fmt.Sprintf("%s%s%s", cl.prefix, string(cl.list[cl.index]), cl.suffix)
	cl.index += 1
	return next, false
}
