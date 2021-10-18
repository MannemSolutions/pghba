package gnrtr

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type array struct {
	prefix     string
	list       []string
	suffix     string
	index      int
	allGnrtrs  []Gnrtr
	subGnrtrs  map[int]Gnrtr
	current    string
	currentRaw string
}

func newGnrtrArray(s string) (a *array, err error) {
	prefix, comprehension, suffix, err := groupChar("(").Parts(s)
	if err != nil {
		return nil, err
	}
	return &array{
		prefix: prefix,
		list:   strings.Split(comprehension, "|"),
		suffix: suffix,
		index:  0,
	}, nil
}

func (a array) Current() string {
	return a.current
}

func (a array) advanceSubGnrtrs() (done bool) {
	for i := range a.subGnrtrs {
		if _, done := a.subGnrtrs[i].Next(); !done {
			// This one still can move to the next
			return true
		}
		// At the end, lets start over on this one
		a.subGnrtrs[i].Reset()
	}
	return false
}

func (a *array) rebuildSubGnrtrs() {
	a.subGnrtrs = make(map[int]Gnrtr, 0)
	reSubGenPlaceHolders := regexp.MustCompile(`\${(\d+)}`)
	matches := reSubGenPlaceHolders.FindAllStringSubmatch(a.currentRaw, -1)
	for _, match := range matches {
		gnrtrId, err := strconv.Atoi(match[1])
		if err != nil {
			panic(fmt.Errorf("cannot convert %s to int", match[1]))
		}
		if gnrtrId >= len(a.allGnrtrs) {
			panic(fmt.Errorf("a placeholder references a non existing subGnrtr"))
		}
		a.subGnrtrs[gnrtrId] = a.allGnrtrs[gnrtrId]
	}
}

func (a *array) rebuildCurrent() {
	a.current = a.currentRaw
	for id, g := range a.subGnrtrs {
		placeholder := fmt.Sprintf("${%d}", id)
		a.current = strings.Replace(a.current, placeholder, g.Current(), -1)
	}
}

func (a *array) Next() (next string, done bool) {
	if a.advanceSubGnrtrs() {
		// One of my children could advance, I don't need to.
		return a.Current(), false
	}
	a.index += 1
	if a.index > len(a.list) + 1 {
		return "", done
	}
	a.currentRaw = fmt.Sprintf("%s%s%s", a.prefix, a.list[a.index], a.suffix)
	a.rebuildSubGnrtrs()

	return a.Current(), false
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
		index:  a.index,
		list:   a.list,
	}
}

func (a array) Unique() Gnrtr {
	return uniqueGnrtr(&a)
}

func (a array) ToList() []string {
	return gnrtrToList(&a)
}
