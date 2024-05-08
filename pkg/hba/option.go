package hba

import (
	"fmt"
	"regexp"
	"strings"
)

// The Options type holds options for authentication methods as a key value map
type Options struct {
	str string
	kv  map[string]string
}

// Consumes the next specified option from 'str' and returns it as a key name and its respective value and the remaining string
// Assumes that only options and their values and optionally a comment preceded by '#' are left in 'str'
func nextOption(str string) (key string, value string, rest string) {
	// TODO this regular expression needs some TLC as it's workings are a bit murky right now
	// It contains three named capture groups (key, value, rest) that are referenced later on.
	// What the regex should do:
	// - discard any preceding whitespace
	// - match an option keyword that may not be double quoted (assumption, double check in hba.c)
	// - match an equals sign '=' no space allowed before or after as far as I understand from hba.c
	// - match a one or more values, which may be in double quotes. For certain values, double-quoting may have been
	//   used because double quotes are part of the value itself.
	re := regexp.MustCompile(`^\S*(?P<key>\s*)\S*=\S*(?P<value>[^"]+|"[^"]*")\S+(?P<rest>.*)$`)
	matches := re.FindStringSubmatch(str)
	if matches == nil {
		return "", "", str
	}
	fields := make(map[string]string)
	// walk over the named sub expressions
	for id, name := range re.SubexpNames() {
		fields[name] = matches[id]
	}
	key, exists := fields["key"]
	if !exists {
		return "", "", str
	}
	value, exists = fields["key"]
	if !exists {
		return "", "", str
	}
	return key, value, fields["rest"]
}

func NewOptionsFromString(str string) (o Options, comment Comment, err error) {
	o.str = str
	str = strings.Trim(str, " \t") // Strip surrounding whitespace
	o.kv = make(map[string]string)
	for {
		if str == "" {
			return o, Comment{}, nil
		}
		k, v, str := nextOption(str)
		if k == "" || v == "" {
			return o, Comment{}, fmt.Errorf("could not read option from %s", str)
		}
		o.kv[k] = v
		if len(str) == 0 {
			break
		}
		if strings.HasPrefix(str, "#") {
			comment, err = NewComment(str)
			if err != nil {
				return o, Comment{}, fmt.Errorf("seems like, but could not be parsed as comment:\n%s\n%e", comment, err)
			}
			return o, comment, nil
		}
	}
	return o, Comment{}, nil
}

func (o Options) Len() int {
	return len(o.kv)
}

func (o Options) String() string {
	return o.str
}

func (o Options) Bare() string {
	var opts []string
	for k, v := range o.kv {
		if strings.Contains(v, " ") {
			v = fmt.Sprintf("\"%s\"", v)
		}
		opt := fmt.Sprintf("%s=%s", k, v)
		opts = append(opts, opt)
	}
	return strings.Join(opts, " ")
}
