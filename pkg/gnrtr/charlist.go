package gnrtr

import (
	"fmt"
	"regexp"
	"strings"
)

type charList struct {
	list   []byte
	index  int
}

func newCharList(s string) (cl *charList, err error) {
	if err != nil {
		return nil, err
	}
	if strings.HasPrefix(s, "^") {
		return nil, fmt.Errorf("cannot make an iterator of a negative character list (starting with ^)")
	}
	cl = &charList{
		index:  0,
	}
	if strings.HasSuffix(s, "-") {
		cl.list = []byte("-")
	}
	if strings.HasPrefix(s, "-") {
		cl.list = []byte("-")
	}
	re := regexp.MustCompile(`(-)?([^-]|([^-][-][^-]))*(-)?`)
	matches := re.FindAllIndex([]byte(s), -1)
	if matches == nil {
		return &charList{}, fmt.Errorf("could not parse charList %s", s)
	}
	for _, match := range matches {
		start := s[match[0]]
		if match[1]-match[0] == 1 {
			cl.list = append(cl.list, start)
		} else if match[1]-match[0] == 3 {
			end := s[match[0]+match[1]-1]
			for c := start; c <= end; c++ {
				cl.list = append(cl.list, c)
			}
		} else {
			return &charList{}, fmt.Errorf("could not parse the %s part of %s", s[match[0]:match[1]],
				s)
		}
	}
	return cl, nil
}

func (cl charList) Current() string {
	if cl.index > len(cl.list) {
		return ""
	}
	return string(cl.list[cl.index-1])
}

func (cl *charList) Next() (next string, done bool) {
	cl.index += 1
	next = cl.Current()
	if next == "" {
		done = true
	}
	return next, done
}

func (cl *charList) Reset() {
	cl.index = 0
}

func (cl *charList) ToArray() (a array) {
	a = array{
		index:  cl.index,
	}
	for i := 0; i < len(cl.list); i++ {
		a.list = append(a.list, string([]byte{cl.list[i]}))
	}
	return a
}

func (cl *charList) String() (s string) {
	return fmt.Sprintf("[%s]", string(cl.list))
}

func (cl charList) ToList() []string {
	return subGnrtrToList(&cl)
}
