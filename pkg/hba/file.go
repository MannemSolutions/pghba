package hba

import (
  "bufio"
  "fmt"
  "os"
)

type File struct {
  path string
  lines Lines
}

func NewFile(path string) (f File, err error) {
  var commentBlock Comments
  file, err := os.Open(path)
  if err != nil {
      return f, err
  }
  defer file.Close()

  scanner := bufio.NewScanner(file)
  var lines Lines
  // optionally, resize scanner's capacity for lines over 64K, see next example
  for scanner.Scan() {
    line, err := parseLine(scanner.Text())
    if err != nil {
      return f, err
    }
    switch l := line.(type) {
    case Comment:
      commentBlock = append(commentBlock, l)
    case Rule:
      l.PrependComments(commentBlock)
      commentBlock = Comments{}
    case EmptyLine:
      lines = append(lines, commentBlock, l)
      commentBlock = Comments{}
    }
  }
  if err := scanner.Err(); err != nil {
    return f, err
  }
  f.path = path
  f.lines = lines
  return f, nil
}

func parseLine(fileLine string) (l Line, err error) {
  l, err = NewComment(fileLine)
  if err == nil {
    return l, nil
  }

  return nil, fmt.Errorf("Line %s has an unknown format as a hba line", fileLine)
}

func (f *File) Delete(r Rule) (found bool) {

}