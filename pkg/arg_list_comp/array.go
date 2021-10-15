package arg_list_comp

import (
	"fmt"
	"strings"
)

type array struct {
	prefix string
	list []string
	suffix string
	index int
	subIterator ALC
	current string
}

func newAlcArray(s string) (a *array, err error) {
	prefix, comprehension, suffix, err := groupChar("(").Parts(s)
	if err != nil {
		return nil, err
	}
	return &array{
		prefix: prefix,
		list: strings.Split(comprehension, "|"),
		suffix: suffix,
		index: 0,
	}, nil
}

func (a array) Current() string {
	return a.current
}

func (a *array) Next() (next string, done bool) {
	if a.subIterator != nil {
		next, done := a.subIterator.Next()
		if done {
			a.subIterator = nil
		} else {
			a.current = next
			return a.current, false
		}
	}
	if a.index >= len(a.list) {
		a.current = ""
		return a.current, true
	}
	next = fmt.Sprintf("%s%s%s", a.prefix, a.list[a.index], a.suffix)
	a.index += 1
	a.subIterator = NewALC(next)
	if a.subIterator != nil {
		// Let s call the method again, just to let the top part handle this
		return a.Next()
	}
	return next, false
}

func (a array) String() (s string) {
	return fmt.Sprintf("%s(%s)%s", a.prefix, strings.Join(a.list, "|"), a.suffix)
}

func (a *array) Reset() {
	a.current = ""
	a.index = 0
}

func (a array) ToArray() array {
	return array{
		prefix: a.prefix,
		suffix: a.suffix,
		index: a.index,
		list: a.list,
	}
}

func (a array) Unique() ALC {
	return uniqueAlc(&a)
}

func (a array) ToList() []string {
	return alcToList(&a)
}