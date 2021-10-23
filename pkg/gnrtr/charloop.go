package gnrtr

import (
	"fmt"
)

type charLoop struct {
	begin  byte
	index  byte
	end    byte
}

func newCharLoop(s string) (l *charLoop, err error) {
	match := reCharLoop.FindStringSubmatch(s)
	if match == nil {
		return nil, fmt.Errorf("invalid input to newIntLoop (should have form %s)", reCharLoop.String())
	}

	l = &charLoop{
		begin: []byte(match[1])[0],
		end: []byte(match[2])[0],
	}
	if l.begin > l.end {
		return nil, fmt.Errorf("invalid format for character loop comprehension (cStart should be before cEnd), is %s)", s)
	}
	l.index = l.begin
	return l, nil
}

func (cl charLoop) Clone()  subGnrtr {
	return &charLoop{
		begin: cl.begin,
		index: cl.index,
		end: cl.end,
	}
}

func (cl charLoop) Index() int {
	return int(cl.index - cl.begin)
}

func (cl charLoop) Current() string {
	if cl.index > cl.end {
		return ""
	}
	return string(cl.index)
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

func (cl charLoop) ToList() []string {
	return subGnrtrToList(&cl)
}
