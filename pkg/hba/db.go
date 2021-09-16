package hba

type Database string;

func (d Database) Compare(other Database) int {
	if d == other {
		return 0
	}
	if d == "all" {
		return -1
	} else if other == "all" {
		return 1
	}
	if d == "samerole" || d == "samegroup" {
		return -1
	} else if other == "samerole" || other == "samegroup" {
		return 1
	}
	if d == "sameuser" {
		return -1
	} else if other == "sameuser" {
		return 1
	}
	if d < other {
		return -1
	}
	return 1
}