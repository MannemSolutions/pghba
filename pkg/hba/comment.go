package hba

import (
  "fmt"
  "regexp"
)

type Comment struct {
  str  string
  bare string
}

type Comments []Comment;

func NewComment(line string) (c Comment, err error) {
  re := regexp.MustCompile(`\S*#\S*(?P<bare>\S*)`)
	matches := re.FindStringSubmatch(line)
	if matches == nil {
		return c, fmt.Errorf("line %s is not a valid Comment", line)
	}
  c.str = line
  fields := make(map[string]string)
  for id, name := range re.SubexpNames() {
    fields[name] = matches[id]
  }

  c.bare = fields["bare"]
  return c, nil
}

func (c Comment) String() (string) {
  return c.str
}

func (c Comment) Comments() (com []Comment) {
  return []Comment{c}
}

func (c Comment) Bare() (s string) {
  return c.bare
}

func (c Comment) Less(l Line) (less bool) {
  // We do not sort comments
  return false
}
