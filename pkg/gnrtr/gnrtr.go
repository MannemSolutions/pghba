package gnrtr

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

type Gnrtr interface {
	Current() string
	Next() (string, bool)
	ToArray() (a array)
	Reset()
	String() (s string)
	Unique() Gnrtr
	ToList() []string
}

func NewGnrtr(s string) (g Gnrtr) {
	var allGnrtrs []Gnrtr
	reIntLoops := regexp.MustCompile(`{(\d+)..(\d+)}`)
	for _, match := range reIntLoops.FindAllStringSubmatch(s, -1) {
		g, err := newIntLoop(match[1])
		if err != nil {
			panic(err)
		}
		placeholder := fmt.Sprintf("${%d}", len(allGnrtrs))
		s = strings.Replace(s, match[0], placeholder, 1)
		allGnrtrs = append(allGnrtrs, g)
	}
	reCharLoops := regexp.MustCompile(`{(\S..\S)}`)
	for _, match := range reCharLoops.FindAllStringSubmatch(s, -1) {
		g, err := newCharLoop(match[1])
		if err != nil {
			panic(err)
		}
		placeholder := fmt.Sprintf("${%d}", len(allGnrtrs))
		s = strings.Replace(s, match[0], placeholder, 1)
		allGnrtrs = append(allGnrtrs, g)
	}
	reCharLists := regexp.MustCompile(`\[([^]]+)\]`)
	for _, match := range reCharLists.FindAllStringSubmatch(s, -1) {
		g, err := newCharList(match[1])
		if err != nil {
			panic(err)
		}
		placeholder := fmt.Sprintf("${%d}", len(allGnrtrs))
		s = strings.Replace(s, match[0], placeholder, 1)
		allGnrtrs = append(allGnrtrs, g)
	}
	if a, err := newGnrtrArray(s); err == nil {
		return a
	}

	return nil
}

// uniqueGnrtr creates an array with sorted unique elements
func uniqueGnrtr(g Gnrtr) Gnrtr {
	// Clone so we can reset
	a := g.ToArray()
	a.Reset()
	//Make unique
	unique := make(map[string]bool)
	for {
		next, done := a.Next()
		if done {
			break
		}
		unique[next] = true
	}
	// Build sorted
	var sorted *array
	sorted = &array{}
	for next := range unique {
		sorted.list = append(sorted.list, next)
	}
	sort.Strings(sorted.list)

	return sorted
}

// SortedArray creates an array with sorted unique elements
func gnrtrToList(g Gnrtr) (l []string) {
	// Clone so we can reset
	a := g.ToArray()
	a.Reset()
	//Make unique
	for {
		next, done := a.Next()
		if done {
			break
		}
		l = append(l, next)
	}
	return l
}

//func StrToGnrtr(s string) Gnrtr {
//	return &array{
//		list: []string{s},
//	}
//}
