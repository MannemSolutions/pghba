package hba

import (
	"fmt"
	"strings"
)

// Rule holds a parsed or newly created pg_hba.conf rule.
type Rule struct {
	rowNum   int      // Unique row number across all rules
	comments Comments // Line comments
	str      string   // The current string being parsed
	connType ConnType // Connection type for the rule
	database Database // The databases affected by the rule
	user     User     // The users affected by the rule
	address  Address  // The addresses (dis)allowed access by the rule
	method   Method   // The authentication method for the rule
	options  Options  // Authentication options
}

// Create a new rule programmatically by providing the individual values for the relevant fields.
func NewRule(rowNum int,
	connType string,
	database string,
	user string,
	address string,
	mask string,
	method string,
	options string) (r Rule, err error) {
	var addr Address
	ct := NewConnType(connType)
	mtd := NewMethod(method)
	db := Database(database) // TODO this is a method being used as a function?
	usr := User(user)
	if ct != ConnTypeLocal && ct != ConnTypeUnknown {
		addr, err = NewAddress(address)
		if err != nil {
			return Rule{}, err
		}
		err = addr.SetMask(mask)
		if err != nil {
			return Rule{}, err
		}
	}
	opts, _, err := NewOptionsFromString(options)
	if err != nil {
		return Rule{}, err
	}

	return Rule{
		rowNum:   rowNum,
		connType: ct,
		method:   mtd,
		database: db,
		user:     usr,
		address:  addr,
		options:  opts,
	}, nil
}

// Create a new rule structure from an existing line in a pg_hba.conf file.
func NewRuleFromLine(line string) (Rule, error) {
	var address, db string
	var r Rule
	parts := newSplitter(line)
	connType, err := parts.Next()
	if err != nil {
		return Rule{}, err
	}
	r.connType = NewConnType(connType)
	if r.connType == ConnTypeUnknown {
		return Rule{}, fmt.Errorf("could not convert line to a Rule: %s", line)
	}
	db, err = parts.Next()
	if err != nil {
		return Rule{}, err
	}
	r.database = Database(db)
	user, err := parts.Next()
	if err != nil {
		return Rule{}, err
	}
	r.user = User(user)
	r.method = MethodUnknown
	if r.connType != ConnTypeLocal {
		address, err = parts.Next()
		if err != nil {
			return Rule{}, err
		}
		if r.address, err = NewAddress(address); err != nil {
			return Rule{}, fmt.Errorf("line %s has an invalid address: %e", line, err)
		}
		next, err := parts.Next()
		if err != nil {
			return Rule{}, err
		}
		r.method = NewMethod(next)
		if r.method == MethodUnknown {
			if err = r.address.SetMask(next); err != nil {
				return Rule{}, fmt.Errorf("%s is not recognizable as method, or as mask: %e", next, err)
			}
		}
	}
	if r.method == MethodUnknown {
		method, err := parts.Next()
		if err != nil {
			return Rule{}, err
		}
		r.method = NewMethod(method)
		if r.method == MethodUnknown {
			return Rule{}, fmt.Errorf("%s is not recognizable as method", method)
		}
	}

	options, comment, err := NewOptionsFromString(parts.Rest())
	if err != nil {
		return Rule{}, fmt.Errorf("could not parse options from %s", parts.Rest())
	}
	r.options = options
	r.comments = append(r.comments, comment)
	return r, nil
}

// Below the implementation of various methods against 'Rule'.

func (r *Rule) PrependComments(comments Comments) {
	// Counterintuitive, but basically we take the 'comments' argument as a start, and append all comments
	// that where already in r.comments. We store the result in r.comments...
	r.comments = append(comments, r.comments...)
}

// Check to see if one rule is contained in another rule.
func (r Rule) Contains(other Rule) bool {
	if r.database != "" && r.database.Compare(other.database) != 0 {
		return false
	}
	if r.user != "" && r.user.Compare(other.user) != 0 {
		return false
	}
	if r.connType != ConnTypeUnknown {
		if r.connType.Compare(other.connType) != 0 {
			return false
		}
		if r.connType != ConnTypeLocal && other.connType != ConnTypeLocal {
			return r.address.Contains(other.address)
		}
	}
	return true
}

// Check if one rule
func (r Rule) Compare(other Line) (comparison int) {
	o, ok := other.(Rule) // TODO I don't understand this construct
	if !ok {
		// We cannot compare rules with other line types
		return 0
	}
	if comparison = r.connType.Compare(o.connType); comparison != 0 {
		return comparison
	}
	if comparison = r.database.Compare(o.database); comparison != 0 {
		return comparison
	}
	if comparison = r.user.Compare(o.user); comparison != 0 {
		return comparison
	}
	if o.connType != ConnTypeLocal {
		if comparison = r.address.Compare(o.address); comparison != 0 {
			return comparison
		}
	}
	return 0
}

func (r Rule) Bare() string {
	var lines []string
	comments := r.comments.String()
	if len(comments) > 0 {
		lines = append(lines, comments)
	}
	var parts []string
	if r.connType == ConnTypeLocal {
		parts = []string{r.connType.String(), string(r.database), string(r.user), r.method.String()}
	} else {
		parts = []string{r.connType.String(), string(r.database), string(r.user), r.address.String(), r.method.String()}
	}
	if r.options.Len() > 0 {
		parts = append(parts, r.options.Bare())
	}
	lines = append(lines, strings.Join(parts, "\t"))
	return strings.Join(lines, "\n")
}

func (r Rule) String() string {
	if r.str == "" {
		return r.Bare()
	}
	var lines []string
	comments := r.comments.String()
	if len(comments) > 0 {
		lines = append(lines, comments)
	}
	lines = append(lines, r.str)
	return strings.Join(lines, "\n")
}

func (r Rule) Comments() (com Comments) {
	return r.comments
}

func (r Rule) Less(l Line) (less bool) {
	return r.Compare(l) < 0
}

func (r Rule) SortByRowNum(l Line) (less bool) {
	if r.rowNum != l.RowNum() {
		return r.rowNum-l.RowNum() < 0
	}
	return r.Compare(l) < 0
}

func (r *Rule) SetRowNum(rowNum int) {
	r.rowNum = rowNum
}

func (r Rule) RowNum() int {
	return r.rowNum
}

func (r Rule) Clone() Rule {
	return Rule{
		rowNum:   r.rowNum,
		comments: r.comments,
		str:      r.str,
		connType: r.connType,
		database: r.database,
		user:     r.user,
		address:  r.address.Clone(),
		method:   r.method,
		options:  r.options,
	}
}
