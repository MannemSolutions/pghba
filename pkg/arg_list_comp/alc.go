package arg_list_comp

import (
	"sort"
)

type ALC interface {
	Current() string
	Next()      (string, bool)
	ToArray()   (a array)
	Reset()
	String()    (s string)
	Unique() ALC
	ToList() []string
}

func NewALC (s string) (alc ALC){
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
func uniqueAlc(alc ALC) (ALC) {
	// Clone so we can reset
	a := alc.ToArray()
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
func alcToList(alc ALC) (l []string) {
	// Clone so we can reset
	a := alc.ToArray()
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
