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
	current string
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

func (cl charLoop) Current() string {
	return cl.current
}

func (cl *charLoop) Next() (next string, done bool) {
	if cl.subIterator != nil {
		next, done := cl.subIterator.Next()
		if done {
			cl.subIterator = nil
		} else {
			cl.current = next
			return cl.current, false
		}
	}
	if cl.index > cl.end {
		cl.current = ""
		return cl.current, true
	}
	next = fmt.Sprintf("%s%s%s", cl.prefix, string(cl.index), cl.suffix)
	cl.index += 1
	cl.subIterator = NewALC(next)
	if cl.subIterator != nil {
		// Let s call the method again, just to let the top part handle this
		return cl.Next()
	}
	return next, false
}

func (cl *charLoop) Reset() {
	cl.index = cl.begin
}

func (cl charLoop) ToArray() (a array) {
	a = array{
		prefix: cl.prefix,
		suffix: cl.suffix,
		index:  int(cl.index- cl.begin),
	}
	for c:= cl.begin;c<= cl.end;c++ {
		a.list = append(a.list, string(c))
	}
	return a
}

func (cl charLoop) String() (s string) {
	return fmt.Sprintf("%s{%s..%s}%s", cl.prefix, string(cl.begin), string(cl.end), cl.suffix)
}

func (cl charLoop) Unique() ALC {
	return uniqueAlc(&cl)
}

func (cl charLoop) ToList() []string {
	return alcToList(&cl)
}