package hba

import (
	"fmt"
	"regexp"
	"strings"
)

type Rule struct {
	comments Comments
	str      string
	connType ConnType
	database Database
	user User
	address Address
	method Method
	options Options
}

func NewRule(connType string, database string, user string, address string, method string, options string) (Rule, error) {
	ct := NewConnType(connType)
	mtd := NewMethod(method)
	db := Database(database)
	usr := User(user)
	addr, err := NewAddress(address)
	if err != nil {
		return Rule{}, err
	}
	opts, _, err := NewOptionsFromString(options)
	if err != nil {
		return Rule{}, err
	}

	if ct == ConnTypeUnknown  || mtd == MethodUnknown {
		return Rule{}, fmt.Errorf("New Rule has an invalid connection type (%s) or method (%s)", connType, method)
	}
	return Rule{
		connType: ct,
		method: mtd,
		database: db,
		user: usr,
		address: addr,
		options: opts,
	}, nil
}

type splitter struct {
	rest string
}

func newSplitter(str string) splitter{
	return splitter{rest: str}
}

func (s *splitter) Next() (string, error) {
	re := regexp.MustCompile(`^\S*(?P<part>([^"]+|"[^"]*"))\S+(?P<rest>.*)$`)
	matches := re.FindStringSubmatch(s.rest)
	if matches == nil {
		return "", fmt.Errorf("could not find next part in %s", s.rest)
	}
	fields := make(map[string]string)
	for id, name := range re.SubexpNames() {
		fields[name] = matches[id]
	}
	part, exists := fields["part"]
	if ! exists {
		return "", fmt.Errorf("next part seems empty in %s", s.rest)
	}
	s.rest = fields["rest"]
	return part, nil
}

func (s splitter) Rest() string {
	return s.rest
}

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
	if err !=  nil {
		return Rule{}, fmt.Errorf("could not parse options from %s", parts.Rest())
	}
	r.options = options
	r.comments = append(r.comments, comment)
	return r, nil
}

func (r *Rule) PrependComments(comments Comments) {
	// Counterintuitive, but basically we take the comments argument as a start, and append all comments
	// that where already in r.comments. We store the result in r.comments...
	r.comments = append(comments, r.comments...)
}

func (r Rule) Compare(other Line) (comparison int) {
	o, ok := other.(Rule)
	if ! ok {
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
	return strings.Join(parts, "\t")
}

func (r Rule) String() string {
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