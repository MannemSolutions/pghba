package hba

import (
	"fmt"
	"regexp"
)

type Line interface {
	String() string
	Comments() Comments
	Bare() string
	Less(Line) bool
	Compare(Line) int
	RowNum() int
}

type EmptyLine string

// TODO I don't get the regex: \S matches any non-space character. Isn't the logic inverted here?
func NewEmptyLine(line string) (EmptyLine, error) {
	re := regexp.MustCompile(`^\S*$`)
	if re.MatchString(line) {
		return EmptyLine(line), nil
	}
	return EmptyLine(""), fmt.Errorf("line is not empty:\n%s", line)
}

func (e EmptyLine) String() string {
	return string(e)
}

func (e EmptyLine) Comments() Comments {
	return Comments{}
}

func (e EmptyLine) Bare() string {
	return ""
}

func (e EmptyLine) Less(Line) bool {
	return false
}

func (e EmptyLine) Compare(Line) int {
	return 0
}

func (e EmptyLine) RowNum() int {
	return 0
}

type Lines []Line
