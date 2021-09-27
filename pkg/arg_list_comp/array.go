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
}

func NewAlcArray(prefix string, comprehension string, suffix string) (a array, err error) {
	if ! (strings.HasPrefix(comprehension, "{") && strings.HasSuffix(comprehension, "}")) {
		return a, fmt.Errorf("missing curly braces around array comprehension '%s'", comprehension)
	}
	comprehension = comprehension[1:len(comprehension)-1]
	return array{
		prefix: prefix,
		list: strings.Split(comprehension, ","),
		suffix: suffix,
		index: 0,
	}, nil
}

func (a array) Next() (next string, done bool) {
	if a.index >= len(a.list) {
		return "", true
	}
	next = fmt.Sprintf("%s%s%s", a.prefix, a.list[a.index], a.suffix)
	a.index += 1
	return next, false
}