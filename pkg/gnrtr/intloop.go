package gnrtr

import (
	"fmt"
	"strconv"
	"strings"
)

type intLoop struct {
	begin  int
	index  int
	end    int
}

func newIntLoop(s string) (l *intLoop, err error) {
	if ! strings.Contains(s, "..") {
		return nil, fmt.Errorf("invalid input to newIntLoop (should contain '..')")
	}
	parts  := strings.Split(s, "..")

	start, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, fmt.Errorf("cannot convert start (%s..) to int in newIntLoop", parts[0])
	}
	end, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, fmt.Errorf("cannot convert end (..%s) to int in newIntLoop", parts[1])
	}
	return &intLoop{
		begin: start,
		end: end,
		index:  0,
	}, nil
}

func (l intLoop) Current() string {
	if l.index > l.end+1 {
		return ""
	}
	return fmt.Sprintf("%d", l.index-1)
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

func (l intLoop) Unique() Gnrtr {
	return uniqueGnrtr(&l)
}

func (l intLoop) ToList() []string {
	return gnrtrToList(&l)
}

func (l intLoop) String() (s string) {
	return fmt.Sprintf("{%d..%d}", l.begin, l.end)
}
