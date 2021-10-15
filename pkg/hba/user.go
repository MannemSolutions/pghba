package hba

import "strings"

type User string

func (u User) Compare(other User) int {
	if u == other {
		return 0
	}
	if u == "all" {
		return -1
	} else if other == "all" {
		return 1
	}
	if strings.HasPrefix(string(u), "+") {
		return -1
	} else if strings.HasPrefix(string(u), "+") {
		return 1
	}
	if u < other {
		return -1
	}
	return 1
}
