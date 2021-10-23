package gnrtr

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type array struct {
	list       []string
	index      int
	allGnrtrs  subGnrtrs
	subGnrtrs  map[int]subGnrtr
	current    string
	currentRaw string
}

func newArray(s string, ag subGnrtrs) (a *array, err error) {
	match := reArray.FindStringSubmatch(s)
	if match == nil {
		return nil, fmt.Errorf("invalid input to newArray (should have form %s)", reArray.String())
	}
	a = &array{
		list:   strings.Split(match[1], "|"),
		index:  0,
		allGnrtrs: ag,
	}
	a.currentRaw = fmt.Sprintf("%s", a.list[a.index])
	a.setCurrent()
	return a, nil
}

func (a array) Index()  int{
	return a.index
}

func (a array) Current() string {
	return a.current
}

func (a *array) advanceSubGnrtrs() (done bool) {
	for i := range a.subGnrtrs {
		sg := a.subGnrtrs[i]
		if _, done := sg.Next(); !done {
			// This one still can move to the next
			return true
		}
		// At the end, lets start over on this one
		sg.Reset()
	}
	return false
}

func (a *array) rebuildSubGnrtrs() {
	a.subGnrtrs = make(map[int]subGnrtr, 0)
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

func (a *array) setCurrent() string {
	a.current = a.currentRaw
	for id, g := range a.subGnrtrs {
		placeholder := fmt.Sprintf("${%d}", id)
		a.current = strings.Replace(a.current, placeholder, g.Current(), -1)
	}
	return a.current
}

func (a *array) Next() (next string, done bool) {
	if a.advanceSubGnrtrs() {
		// One of my children could advance, I don't need to.
		return a.setCurrent(), false
	}
	a.index += 1
	if a.index >= len(a.list) {
		return "", true
	}
	a.currentRaw = fmt.Sprintf("%s", a.list[a.index])
	a.rebuildSubGnrtrs()

	return a.setCurrent(), false
}

func (a array) String() (s string) {
	s = fmt.Sprintf("(%s)", strings.Join(a.list, "|"))
	reSubGenPlaceHolders := regexp.MustCompile(`\${(\d+)}`)
	matches := reSubGenPlaceHolders.FindAllStringSubmatch(s, -1)
	for _, match := range matches {
		gnrtrId, err := strconv.Atoi(match[1])
		if err != nil {
			panic(fmt.Errorf("cannot convert %s to int", match[1]))
		}
		if gnrtrId >= len(a.allGnrtrs) {
			panic(fmt.Errorf("a placeholder references a non existing subGnrtr"))
		}
		s = strings.Replace(s, fmt.Sprintf("${%d}", gnrtrId), a.allGnrtrs[gnrtrId].String(), 1)
	}
	return s
}

func (a *array) Reset() {
	a.index = 0
	a.setCurrent()
}

func (a array) ToArray() array {
	return array{
		index:  a.index,
		list:   a.list,
	}
}

func (a array) ToList() []string {
	return subGnrtrToList(&a)
}
