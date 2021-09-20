package hba

import (
  "bufio"
  "fmt"
  "os"
)

type File struct {
  path string
  lines Lines
  dirty bool
}

func NewFile(path string) File {
  return File{
    path: path,
  }
}

func (f File) Read() error {
  var commentBlock Comments
  file, err := os.Open(f.path)
  if err != nil {
      return err
  }
  defer file.Close()

  scanner := bufio.NewScanner(file)
  var lines Lines
  // optionally, resize scanner's capacity for lines over 64K, see next example
  for scanner.Scan() {
    line := scanner.Text()
    e, err := NewEmptyLine(line)
    if err == nil {
      lines = append(lines, commentBlock, e)
      commentBlock = Comments{}
      continue
    }
    c, err := NewComment(line)
    if err == nil {
      commentBlock = append(commentBlock, c)
      continue
    }
    r, err := NewRuleFromLine(line)
    if err == nil {
      r.PrependComments(commentBlock)
      commentBlock = Comments{}
      continue
    }
    return fmt.Errorf("could not parse this hba line %s", line)
  }
  if err := scanner.Err(); err != nil {
    return err
  }
  f.lines = lines
  return nil
}

func (f *File) DeleteRule(r Rule) (found bool) {
  for i, line:= range f.lines {
    rule, isRule := line.(Rule)
    if ! isRule {
      continue
    }
    if rule.Compare(r) == 0 {
      f.lines = append(f.lines[:i], f.lines[i+1:]...)
      found = true
      f.dirty = true
    }
  }
  return found
}

func (f *File) Save() error {
  file, err := os.Create(f.path)
  if err != nil {
    return err
  }
  defer file.Close()
  for _, line := range f.lines {
    _, err := file.WriteString(line.String() + "\n")
    if err != nil {
      return err
    }
  }
  err = file.Sync()
  if err != nil {
    return err
  }
  return nil
}