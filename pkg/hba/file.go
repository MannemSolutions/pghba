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
func (f *File) renumberRules() {
	for i, line := range f.lines {
		rule, isRule := line.(Rule)
		if !isRule {
			continue
		}
		f.numRules = i + 1
		rule.SetRowNum(f.numRules)
		f.lines[i] = rule
	}
}

func (f *File) Read() error {
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
			lines = append(lines, r)
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
	f.dirty = true
	f.renumberRules()
	return nil
}

func (f *File) DeleteRule(r Rule) (found bool) {
	var i = 0
	for {
		if i >= len(f.lines) {
			break
		}
		if rule, isRule := f.lines[i].(Rule); !isRule {
			log.Debugf("Not a rule %s", f.lines[i].String())
			i+=1
			continue
		} else if r.Contains(rule) {
			log.Debugf("Removing rule %s", rule.String())
			f.lines = append(f.lines[:i], f.lines[i+1:]...)
			found = true
			f.dirty = true
			f.numRules -= 1
		} else {
			i+=1
			log.Debugf("Leaving rule %s", rule.String())
		}
	}
	if found {
		f.renumberRules()
	}
	return found
}

func (f *File) DeleteRules(rs *Rules) (found bool) {
	for _, r := range rs.rules {
		log.Debugf("Cleaning rule %s", r.String())
		found = f.DeleteRule(r) || found
	}
	return found
}

func (f *File) InsertRule(r Rule, index int) {
	f.lines = append(f.lines[:index+1], f.lines[index:]...)
	f.lines[index+1] = r
	f.renumberRules()
	log.Debugf("Insert rule %s on line %d. Lines: %d", r.String(), index, len(f.lines))
}

func (f *File) AddRule(r Rule, auto bool) (changed bool) {
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
			cmp := r.Compare(rule)
			if cmp == 0 {
				return false
			} else if cmp < 0 {
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
	return true
}

func (f *File) AddRules(rs Rules, auto bool) (changed bool) {
	log.Debugf("Count rules: %d", len(rs.rules))
	for _, r := range rs.rules {
		// First run AddRule and then add to found, or AddRule will not be run when found is true
		changed = f.AddRule(r, auto) || changed
	}
	return changed
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
