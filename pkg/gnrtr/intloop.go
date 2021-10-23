package gnrtr

import (
	"fmt"
	"strconv"
)

type intLoop struct {
	begin  int
	index  int
	end    int
}

func newIntLoop(s string) (l *intLoop, err error) {
	match := reIntLoop.FindStringSubmatch(s)
	if match == nil {
		return nil, fmt.Errorf("invalid input to newIntLoop (should have form %s)", reIntLoop.String())
	}
	l = &intLoop{}
	l.begin, err = strconv.Atoi(match[1])
	if err != nil {
		return nil, fmt.Errorf("cannot convert start (%s..) to int in newIntLoop", match[0])
	}
	l.end, err = strconv.Atoi(match[2])
	if err != nil {
		return nil, fmt.Errorf("cannot convert end (..%s) to int in newIntLoop", match[1])
	}
	l.Reset()
	return l, nil
}

func (l intLoop) Clone()  subGnrtr {
	return &intLoop{
		begin: l.begin,
		index: l.index,
		end: l.end,
	}
}

func (l intLoop) Index() int {
	return l.index
}

func (l intLoop) Current() string {
	if l.index > l.end {
		return ""
	}
	return fmt.Sprintf("%d", l.index)
}

func (l *intLoop) Next() (next string, done bool) {
	l.index += 1
	next = l.Current()
	if next == "" {
		done = true
	}
	return l.Current(), done
}

func (l *intLoop) Reset() {
	l.index = l.begin
}

func (l intLoop) ToArray() (a array) {
	a = array{
		index:  l.index - l.begin,
	}
	for i := l.begin; i <= l.end; i++ {
		a.list = append(a.list, fmt.Sprintf("%d", i))
	}
	return a
}

func (l intLoop) ToList() []string {
	return subGnrtrToList(&l)
}

func (l intLoop) String() (s string) {
	return fmt.Sprintf("{%d..%d}", l.begin, l.end)
}
