package arg_list_comp

import (
	"fmt"
	"strconv"
	"strings"
)

type loop struct {
	prefix string
	index int
	end int
	suffix string
}

func NewAlcLoop(prefix string, comprehension string, suffix string) (l loop, err error) {
	if ! (strings.HasPrefix(comprehension, "(") && strings.HasSuffix(comprehension, ")")) {
		return loop{}, fmt.Errorf("missing round braces around array comprehension")
	}
	l = loop {
		prefix: prefix,
		suffix: suffix,
		index: 0,
	}
	csplit := strings.Split(comprehension[1:len(comprehension)-1], "..")
	if len (csplit) != 2 {
		return loop{}, fmt.Errorf("invalid format for array comprehension (should be (iStart..iEnd), is %s)", comprehension)
	}
	l.index, err = strconv.Atoi(csplit[0])
	if err != nil {
		return loop{}, fmt.Errorf("invalid format for array comprehension (first part should be integer, but is %s). %e", csplit[0], err)
	}
	l.end, err = strconv.Atoi(csplit[1])
	if err != nil {
		return loop{}, fmt.Errorf("invalid format for array comprehension (second part should be integer, but is %s). %e", csplit[1], err)
	}
	return l, nil
}

func (l loop) Next() (next string, done bool) {
	if l.index > l.end {
		return "", true
	}
	next = fmt.Sprintf("%s%d%s", l.prefix, l.index, l.suffix)
	l.index += 1
	return next, false
}
