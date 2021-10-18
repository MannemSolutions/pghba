package gnrtr

import (
	"fmt"
	"strings"
)

type charLoop struct {
	begin  byte
	index  byte
	end    byte
}

func newCharLoop(s string) (l *charLoop, err error) {
	if ! strings.Contains(s, "..") {
		return nil, fmt.Errorf("invalid input to newIntLoop (should contain '..')")
	}
	parts  := strings.Split(s, "..")

	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format for character loop comprehension (should be {cStart..cEnd}, is %s)", s)
	}
	if len(parts[0]) != 1 {
		return nil, fmt.Errorf("invalid format for character loop comprehension (cStart should be 1 character), is %s)", parts[0])
	}
	if len(parts[1]) != 1 {
		return nil, fmt.Errorf("invalid format for character loop comprehension (cEnd should be 1 character), is %s)", parts[1])
	}
	l = &charLoop{
		begin: []byte(parts[0])[0],
		end: []byte(parts[1])[0],
	}
	if l.begin > l.end {
		return nil, fmt.Errorf("invalid format for character loop comprehension (cStart should be before cEnd), is %s)", s)
	}
	l.index = l.begin
	return l, nil
}

func (cl charLoop) Current() string {
	if cl.index > cl.end+1 {
		return ""
	}
	return string(cl.index-1)
}

func (cl *charLoop) Next() (next string, done bool) {
	cl.index += 1
	next = cl.Current()
	if next == "" {
		done = true
	}
	return next, done
}

func (cl *charLoop) Reset() {
	cl.index = cl.begin
}

func (cl charLoop) ToArray() (a array) {
	a = array{
		index:  int(cl.index - cl.begin),
	}
	for c := cl.begin; c <= cl.end; c++ {
		a.list = append(a.list, string(c))
	}
	return a
}

func (cl charLoop) String() (s string) {
	return fmt.Sprintf("{%s..%s}", string(cl.begin), string(cl.end))
}

func (cl charLoop) Unique() Gnrtr {
	return uniqueGnrtr(&cl)
}

func (cl charLoop) ToList() []string {
	return gnrtrToList(&cl)
}
