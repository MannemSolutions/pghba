package hba

import (
  "fmt"
  "regexp"
	"strings"
)

type Comment struct {
  str  string
  bare string
}

func NewComment(line string) (c Comment, err error) {
  re := regexp.MustCompile(`\S*#\S*(?P<bare>.*?)\S*`)
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

func (c Comment) String() string {
  return c.str
}

func (c Comment) Comments() (com Comments) {
  return []Comment{c}
}

func (c Comment) Bare() (s string) {
	return fmt.Sprintf("# %s", c.bare)
}

func (c Comment) Less(_ Line) (less bool) {
  // We do not sort comments
  return false
}

func (c Comment) Compare(Line) int {
	return -1
}

func (c Comment) RowNum() int {
	return 0
}

type Comments []Comment;

func (cb Comments) String() string {
	var strList []string
	for _, c:= range cb {
		strList = append(strList, c.String())
	}
	return strings.Join(strList, "\n")
}

func (cb Comments) Comments() (com Comments) {
	return cb
}

func (cb Comments) Bare() (bare string) {
	var bareList []string
	for _, c:= range cb {
		bareList = append(bareList, c.Bare())
	}
	return strings.Join(bareList, "\n")
}

func (cb Comments) Less(_ Line) (less bool) {
	// We do not sort comments
	return false
}

func (cb Comments) Compare(Line) int {
	return -1
}

func (cb Comments) RowNum() int {
	return 0
}
