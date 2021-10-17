package gnrtr

import (
	"sort"
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
	if cl, err := newGnrtrCharList(s); err == nil {
		return cl
	}
	if a, err := newGnrtrArray(s); err == nil {
		return a
	}
	if l, err := newGnrtrLoop(s); err == nil {
		return l
	}
	if l, err := newGnrtrCharLoop(s); err == nil {
		return l
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
