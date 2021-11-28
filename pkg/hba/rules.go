package hba

import (
	"fmt"
	"github.com/mannemsolutions/pghba/pkg/gnrtr"
)

/*
This is a placeholder for all rules combined.
And it can do smart things on a list of rules.
*/

type Rules struct {
	rules     []Rule
}

func NewRules(rowNum int, connTypes string, databases string, users string, addresses string, mask string, method string, options string) (*Rules, error) {
	rules := &Rules{}
	for _, ct := range gnrtr.NewGnrtr(connTypes).ToList() {
		for _, db := range gnrtr.NewGnrtr(databases).ToList() {
			for _, usr := range gnrtr.NewGnrtr(users).ToList() {
				for _, adrs := range gnrtr.NewGnrtr(addresses).ToList() {
					rule, err := NewRule(rowNum, ct, db, usr, adrs, mask, method, options)
					log.Debugf("new rule: %s", rule.String())
					if err != nil {
						return nil, fmt.Errorf("error while generating all rules: %e", err)
					}
					rules.rules = append(rules.rules, rule)
				}
			}
		}
	}
	return rules, nil
}

func (rs Rules) Clone() (clone *Rules) {
	clone = &Rules{}
	for _, rule := range rs.rules {
		clone.rules = append(clone.rules, rule)
	}
	return clone
}

func (rs *Rules) Merge(from Rules) {
	rs.rules = append(rs.rules, from.rules...)
	rs.SmartSort()
}

func (rs *Rules) Sort() {
	// Closure that dictates the order of the Rule structure.
	normalSort := func(r1, r2 *Rule) bool {
		return r1.Less(r2)
	}
	// Sort the rules by the various criteria.
	rulesSortBy(normalSort).Sort(rs.rules)
}

func (rs *Rules) SmartSort() {
	withRowNum := func(r1, r2 *Rule) bool {
		return r1.SortByRowNum(r2)
	}
	rulesSortBy(withRowNum).Sort(rs.rules)
	var i int
	var rule *Rule
	var next *Rule
	for {
		if i >= len(rs.rules) - 1 {
			break
		}
		rule = &rs.rules[i]
		next = &rs.rules[i+1]
		if rule.Compare(next) == 0 {
			// Remove duplicate
			rs.rules = append(rs.rules[:i], rs.rules[i+1:]...)
		} else {
			rule.rowNum = i
			// No duplicate, next...
			i+=1
		}
	}
}

func (rs *Rules) Renumber() {
	for i := range rs.rules {
		rs.rules[i].rowNum = i
	}
}

func (rs *Rules) Add(r Rule) (found bool) {
	for i := range rs.rules {
		if r.Less(rs.rules[i]) {
			// extend slice (rs.rules[i] now is same as rs.rules[i+1])
			rs.rules = append(rs.rules[:i+1], rs.rules[i:]...)
			// set rs.rules[i] to what it actually should be
			rs.rules[i] = r
			return true
		}
	}
	// No hit, then we should append to end
	rs.rules = append(rs.rules, r)
	return false
}

func (rs *Rules) Remove(r Rule) (found bool) {
	var i int
	for {
		if i >= len(rs.rules) {
			return found
		}
		if r.Compare(rs.rules[i]) == 0 {
			rs.rules = append(rs.rules[:i], rs.rules[i+1:]...)
			found = true
		}
	}
}

