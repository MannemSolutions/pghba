package hba

import (
  "fmt"
  "regexp"
  "strings"
)

type Options struct {
  str string
  kv  map[string]string
}

func nextOption(str string) (key string, value string, rest string) {
  re := regexp.MustCompile(`^\S*(?P<key>\s*)\S*=\S*(?P<value>[^"]+|"[^"]*")\S+(?P<rest>.*)$`)
  matches := re.FindStringSubmatch(str)
  if matches == nil {
    return "", "", str
  }
  fields := make(map[string]string)
  for id, name := range re.SubexpNames() {
    fields[name] = matches[id]
  }
  key, exists := fields["key"]
  if ! exists {
    return "", "", str
  }
  value, exists = fields["key"]
  if ! exists {
    return "", "", str
  }
  return key, value, fields["rest"]
}

func NewOptionsFromString(str string) (o Options, comment Comment, err error){
  o.str = str
  str = strings.Trim(str, " \t")
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