package hba

import (
	"fmt"
	"net"
	"regexp"
	"strings"
)

type Rule struct {
	connType ConnType
	database Database
	user User
	address net.IPNet
	method Method
	options Options
	comments string
}

func NewRule(line string) (r Rule, err error) {
	reMethods := strings.Join(methods(), "|")
	re := regexp.MustCompile(`(?P<type>\S+)\s+(?P<database>\S+)\s+(?P<user>\S+)\s+(?P<address>.*?)\s*(?P<method>(`+reMethods+`))(?P<options>.*)`)
	matches := re.FindStringSubmatch(line)
	if matches == nil {
		return Rule{}, fmt.Errorf("line %s is not a valid Rule", line)
	}
	fields := make(map[string]string)
	for id, name := range re.SubexpNames() {
		fields[name] = matches[id]
	}
	r.connType = NewConnType(fields["type"])
	if r.connType == ConnTypeUnknown {
		return r, fmt.Errorf("line %s has an invalid Connection Type %s", line, fields["type"])
	}
	r.database = fields["database"]
	r.user = fields["user"]
	if r.connType != ConnTypeLocal {
		hostOrIp := fields["address"]
	}

	return r, nil
}

func (r Rule) Compare(other Rule) (comparison int) {
	if comparison = r.connType.Compare(other.connType); comparison != 0 {
		return comparison
	}
	if comparison = r.database.Compare(other.database); comparison != 0 {
		return comparison
	}
	if comparison = r.user.Compare(other.user); comparison != 0 {
		return comparison
	}
	return 0
}