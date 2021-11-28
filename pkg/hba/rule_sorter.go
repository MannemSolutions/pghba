package hba

import (
	"sort"
)

// By is the type of a "less" function that defines the ordering of its Rule arguments.
type rulesSortBy func(p1, p2 *Rule) bool

// Sort is a method on the function type, By, that sorts the argument slice according to the function.
func (rulesSortBy rulesSortBy) Sort(rules []Rule) {
	ps := &ruleSorter{
		rules: rules,
		by:      rulesSortBy, // The Sort method's receiver is the function (closure) that defines the sort order.
	}
	sort.Sort(ps)
}

// ruleSorter joins a By function and a slice of []Rule to be sorted.
type ruleSorter struct {
	rules []Rule
	by      func(p1, p2 *Rule) bool // Closure used in the Less method.
}

// Len is part of sort.Interface.
func (s *ruleSorter) Len() int {
	return len(s.rules)
}

// Swap is part of sort.Interface.
func (s *ruleSorter) Swap(i, j int) {
	s.rules[i], s.rules[j] = s.rules[j], s.rules[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (s *ruleSorter) Less(i, j int) bool {
	return s.by(&s.rules[i], &s.rules[j])
}


