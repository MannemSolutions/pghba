package arg_list_comp

import (
	"fmt"
	"strconv"
	"strings"
)

type loop struct {
	prefix string
	begin int
	index int
	end int
	suffix string
	subIterator ALC
}

func newAlcLoop(s string) (l *loop, err error) {
	prefix, comprehension, suffix, err := parts(s, "{")
	if err != nil {
		return nil, err
	}
	l = &loop {
		prefix: prefix,
		suffix: suffix,
		index: 0,
	}
	csplit := strings.Split(comprehension, "..")
	if len (csplit) != 2 {
		return nil, fmt.Errorf("invalid format for array comprehension (should be (iStart..iEnd), is %s)", comprehension)
	}
	l.begin, err = strconv.Atoi(csplit[0])
	if err != nil {
		return nil, fmt.Errorf("invalid format for array comprehension, first part (%s) should be an integer: %e", csplit[0], err)
	}

	l.index = l.begin
	l.end, err = strconv.Atoi(csplit[1])
	if err != nil {
		return nil, fmt.Errorf("invalid format for array comprehension, second part (%s) should be an integer: %e", csplit[1], err)
	}
	return l, nil
}

func (l *loop) Next() (next string, done bool) {
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
	next = fmt.Sprintf("%s%d%s", l.prefix, l.index, l.suffix)
	l.index += 1
	l.subIterator = NewALC(next)
	if l.subIterator != nil {
		// Let s call the method again, just to let the top part handle this
		return l.Next()
	}
	return next, false
}

func (l *loop) Reset() {
	l.index = l.begin
}

func (l loop) ToArray() (a array) {
	a = array{
		prefix: l.prefix,
		suffix: l.suffix,
		index: l.index-l.begin,
	}
	for i:=l.begin;i<=l.end;i++ {
		a.list = append(a.list, fmt.Sprintf("%d", i))
	}
	return a
}


func (l loop) String() (s string) {
	return fmt.Sprintf("%s{%d..%d}%s", l.prefix, l.begin, l.end, l.suffix)
}
