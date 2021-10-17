package gnrtr

import (
	"fmt"
	"strconv"
	"strings"
)

type intLoop struct {
	prefix string
	begin int
	index int
	end int
	suffix string
}

func newAlcLoop(s string) (l *intLoop, err error) {
	prefix, comprehension, suffix, err := groupChar("{").Parts(s)
	if err != nil {
		return nil, err
	}
	l = &intLoop {
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

func (l intLoop) Current() string {
	if l.index > l.end +1 {
		return ""
	}
	return fmt.Sprintf("%s%d%s", l.prefix, l.index - 1 , l.suffix)
}

func (l *intLoop) Next() (next string, done bool) {
	l.index += 1
	next = l.Current()
	if next == "" {
		done = true
	}
	return l.Current(), done
}

func (l *intLoop) Reset() {
	l.index = l.begin
}

func (l intLoop) ToArray() (a array) {
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

func (l intLoop) Unique() Gnrtr {
	return uniqueAlc(&l)
}

func (l intLoop) ToList() []string {
	return alcToList(&l)
}

func (l intLoop) String() (s string) {
	return fmt.Sprintf("%s{%d..%d}%s", l.prefix, l.begin, l.end, l.suffix)
}
