package hba

import (
	"bufio"
	"fmt"
	"os"
)

type File struct {
	path     string
	lines    Lines
	dirty    bool
	numRules int
}

func NewFile(path string) File {
	return File{
		path: path,
	}
}

// Renumber all rules in the .lines sub value.
// Not all lines are rules, and we want to number all rules consistently.
// And know how many rules are there in the lines.
func (f File) renumberRules() {
	for i, line := range f.lines {
		rule, isRule := line.(Rule)
		if !isRule {
			continue
		}
		f.numRules = i + 1
		rule.SetRowNum(f.numRules)
	}
}

func (f File) Read() error {
	var commentBlock Comments
	file, err := os.Open(f.path)
	if err != nil {
		if os.IsNotExist(err) {
			// If the file does not exist, we will create with write. No file is no rules, so we are done here...
			return nil
		}
		return err
	}
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
	err = file.Close()
	if err != nil {
		return err
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	f.lines = lines
	f.renumberRules()
	return nil
}

func (f *File) DeleteRule(r Rule) (found bool) {
	for i, line := range f.lines {
		rule, isRule := line.(Rule)
		if !isRule {
			continue
		}
		if rule.Compare(r) == 0 {
			f.lines = append(f.lines[:i], f.lines[i+1:]...)
			found = true
			f.dirty = true
			f.numRules -= 1
		}
	}
	if found {
		f.renumberRules()
	}
	return found
}

func (f *File) DeleteRules(rs *Rules) (found bool, err error) {
	for {
		next, done, err := rs.Next()
		if done {
			return found, nil
		}
		if err != nil {
			return found, err
		}
		found = found || f.DeleteRule(next)
	}
}

func (f File) InsertRule(r Rule, index int) {
	lines := append(f.lines[:index], r)
	f.lines = append(lines, f.lines[index:]...)
	f.renumberRules()
}

func (f *File) AddRule(r Rule, auto bool) (found bool) {
	if !auto {
		log.Debugf("Adding rule %s (manual)", r.String())
		for i, line := range f.lines {
			if line.RowNum() > r.RowNum() {
				f.InsertRule(r, i)
				f.dirty = true
				return true
			}
		}
	} else {
		log.Debugf("Adding rule %s (auto)", r.String())
		for i, line := range f.lines {
			rule, isRule := line.(Rule)
			if !isRule {
				continue
			}
			if rule.Less(r) {
				f.InsertRule(r, i)
				f.dirty = true
				return true
			}
		}
	}
	log.Debugf("Appending rule %s (not found).", r.String())
	f.lines = append(f.lines, r)
	f.renumberRules()
	f.dirty = true
	return false
}

func (f *File) AddRules(rs *Rules, auto bool) (found bool, err error) {
	for {
		next, done, err := rs.Next()
		if done {
			return found, nil
		}
		if err != nil {
			return found, err
		}
		found = found || f.AddRule(next, auto)
	}
}

func (f *File) Save(force bool) error {
	if !(f.dirty || force) {
		return nil
	}
	file, err := os.Create(f.path)
	if err != nil {
		return err
	}
	for _, line := range f.lines {
		_, err := file.WriteString(line.String() + "\n")
		if err != nil {
			_ = file.Close()
			return err
		}
	}
	err = file.Sync()
	if err != nil {
		return err
	}
	err = file.Close()
	if err != nil {
		return err
	}
	return nil
}
