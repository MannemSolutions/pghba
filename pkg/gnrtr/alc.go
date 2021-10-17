package gnrtr

import (
	"sort"
)

type Gnrtr interface {
	Current() string
	Next()      (string, bool)
	ToArray()   (a array)
	Reset()
	String()    (s string)
	Unique() Gnrtr
	ToList() []string
}

func NewGnrtr (s string) (alc Gnrtr){
	if cl, err := newAlcCharList(s); err == nil {
		return cl
	}
	if a, err := newAlcArray(s); err == nil {
		return a
	}
	if l, err := newAlcLoop(s); err == nil {
		return l
	}
	if l, err := newAlcCharLoop(s); err == nil {
		return l
	}
	return nil
}

// uniqueAlc creates an array with sorted unique elements
func uniqueAlc(g Gnrtr) Gnrtr {
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
func alcToList(g Gnrtr) (l []string) {
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

//func StrToALC(s string) ALC {
//	return &array{
//		list: []string{s},
//	}
//}

