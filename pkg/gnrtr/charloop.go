package gnrtr

import (
	"fmt"
	"strings"
)

type charLoop struct {
	prefix string
	begin  byte
	index  byte
	end    byte
	suffix string
}

func newAlcCharLoop(s string) (l *charLoop, err error) {
	prefix, comprehension, suffix, err := groupChar("{").Parts(s)
	if err != nil {
		return nil, err
	}
	l = &charLoop{
		prefix: prefix,
		suffix: suffix,
	}
	csplit := strings.Split(comprehension, "..")
	if len(csplit) != 2 {
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
	if cl.index > cl.end+1 {
		return ""
	}
	return fmt.Sprintf("%s%s%s", cl.prefix, string(cl.index-1), cl.suffix)
}

func (cl *charLoop) Next() (next string, done bool) {
	cl.index += 1
	next = cl.Current()
	if next == "" {
		done = true
	}
	return next, done
}

func (cl *charLoop) Reset() {
	cl.index = cl.begin
}

func (cl charLoop) ToArray() (a array) {
	a = array{
		prefix: cl.prefix,
		suffix: cl.suffix,
		index:  int(cl.index - cl.begin),
	}
	for c := cl.begin; c <= cl.end; c++ {
		a.list = append(a.list, string(c))
	}
	return a
}

func (cl charLoop) String() (s string) {
	return fmt.Sprintf("%s{%s..%s}%s", cl.prefix, string(cl.begin), string(cl.end), cl.suffix)
}

func (cl charLoop) Unique() Gnrtr {
	return uniqueAlc(&cl)
}

func (cl charLoop) ToList() []string {
	return alcToList(&cl)
}
