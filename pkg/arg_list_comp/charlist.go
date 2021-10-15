package arg_list_comp

import (
	"fmt"
	"regexp"
	"strings"
)

type charList struct {
	prefix string
	list []byte
	suffix string
	index int
	subIterator ALC
}

func newAlcCharList(s string) (cl *charList, err error) {
	prefix, comprehension, suffix, err := groupChar("[").Parts(s)
	if err != nil {
		return nil, err
	}
	if strings.HasPrefix(comprehension, "^") {
		return nil, fmt.Errorf("cannot make an iterator of a negative character list (starting with ^)")
	}
	cl = &charList{
		prefix: prefix,
		suffix: suffix,
		index: 0,
	}
	if strings.HasSuffix(comprehension, "-" ) {
		cl.list = []byte("-")
	}
	if strings.HasPrefix(comprehension, "-" ) {
		cl.list = []byte("-")
	}
	re := regexp.MustCompile(`(-)?([^-]|([^-][-][^-]))*(-)?`)
	matches := re.FindAllIndex([]byte(comprehension), -1)
	if matches == nil {
		return &charList{}, fmt.Errorf("could not parse charList %s", comprehension)
	}
	for _, match := range matches {
		start := comprehension[match[0]]
		if match[1] - match[0] == 1 {
			cl.list = append(cl.list, start)
		} else if match[1] - match[0] == 3 {
			end := comprehension[match[0]+match[1]-1]
			for c:= start; c <= end; c++ {
				cl.list = append(cl.list, c)
			}
		} else {
			return &charList{}, fmt.Errorf("could not parse the %s part of %s", comprehension[match[0]:match[1]],
				comprehension)
		}
	}
	return cl, nil
}

func (cl *charList) Next() (next string, done bool) {
	if cl.subIterator != nil {
		next, done := cl.subIterator.Next()
		if done {
			cl.subIterator = nil
		} else {
			return next, false
		}
	}
	if cl.index >= len(cl.list) {
		return "", true
	}
	next = fmt.Sprintf("%s%s%s", cl.prefix, string(cl.list[cl.index]), cl.suffix)
	cl.index += 1
	cl.subIterator = NewALC(next)
	if cl.subIterator != nil {
		// Let s call the method again, just to let the top part handle this
		return cl.Next()
	}
	return next, false
}

func (cl *charList) Reset() {
	cl.index = 0
}

func (cl *charList) ToArray() (a array) {
	a = array{
		prefix: cl.prefix,
		suffix: cl.suffix,
		index: cl.index,
	}
	for i:=0;i<len(cl.list); i++ {
		a.list = append(a.list, string([]byte{cl.list[i]}))
	}
	return a
}

func (cl *charList) String() (s string) {
	return fmt.Sprintf("%s[%s]%s", cl.prefix, string(cl.list), cl.suffix)
}

func (cl charList) Unique() ALC {
	return uniqueAlc(&cl)
}
func (cl charList) ToList() []string {
	return alcToList(&cl)
}
