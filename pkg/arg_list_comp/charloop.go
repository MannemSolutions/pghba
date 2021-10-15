package arg_list_comp

import (
	"fmt"
	"strings"
)

type charLoop struct {
	prefix string
	begin byte
	index byte
	end byte
	suffix string
	subIterator ALC
}

func newAlcCharLoop(s string) (l *charLoop, err error) {
	prefix, comprehension, suffix, err := groupChar("{").Parts(s)
	if err != nil {
		return nil, err
	}
	l = &charLoop {
		prefix: prefix,
		suffix: suffix,
	}
	csplit := strings.Split(comprehension, "..")
	if len (csplit) != 2 {
		return nil, fmt.Errorf("invalid format for character loop comprehension (should be {cStart..cEnd}, is %s)", comprehension)
	}
	if len(csplit[0]) != 1 {
		return nil, fmt.Errorf("invalid format for character loop comprehension (cStart should be 1 character), is %s)", csplit[0])
	}
	if len(csplit[1]) != 1 {
		return nil, fmt.Errorf("invalid format for character loop comprehension (cEnd should be 1 character), is %s)", csplit[1])
	}
	l.begin = []byte(csplit[0])[0]
	l.end = []byte(csplit[1])[0]
	if l.begin > l.end {
		return nil, fmt.Errorf("invalid format for character loop comprehension (cStart should be before cEbd), is %s)", comprehension)
	}
	l.index = l.begin
	return l, nil
}

func (l *charLoop) Next() (next string, done bool) {
	if l.subIterator != nil {
		next, done := l.subIterator.Next()
		if done {
			l.subIterator = nil
		} else {
			return next, false
		}
	}
	if l.index > l.end {
		return "", true
	}
	next = fmt.Sprintf("%s%s%s", l.prefix, string(l.index), l.suffix)
	l.index += 1
	l.subIterator = NewALC(next)
	if l.subIterator != nil {
		// Let s call the method again, just to let the top part handle this
		return l.Next()
	}
	return next, false
}

func (l *charLoop) Reset() {
	l.index = l.begin
}

func (l charLoop) ToArray() (a array) {
	a = array{
		prefix: l.prefix,
		suffix: l.suffix,
		index: int(l.index-l.begin),
	}
	for c:=l.begin;c<=l.end;c++ {
		a.list = append(a.list, string(c))
	}
	return a
}

func (cl charLoop) String() (s string) {
	return fmt.Sprintf("%s{%s..%s}%s", cl.prefix, string(cl.begin), string(cl.end), cl.suffix)
}

func (cl charLoop) ToSortedArray() array {
	return alcToSortedArray(&cl)
}

func (cl charLoop) ToList() []string {
	return alcToList(&cl)
}