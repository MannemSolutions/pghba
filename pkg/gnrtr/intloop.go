package gnrtr

import (
	"fmt"
	"strconv"
)

// intLoop holds the description of a continuous range of integers starting at
// 'begin', ending at 'end' with the current position being 'index'. This type
// has a set of methods defined to adhere to the subGnrtr interface.
type intLoop struct {
	begin int
	index int
	end   int
}

// Return 'l' as intLoop parsed from string 's'.
func newIntLoop(s string) (l *intLoop, err error) {
	match := reIntLoop.FindStringSubmatch(s) // find occurrences of the intLoop pattern.
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

// Below follow the implementations of the required methods to comply with the
// subGnrtr interface.

// clone() returns a value copied clone of intLoop.
func (l intLoop) clone() subGnrtr {
	return &intLoop{
		begin: l.begin,
		index: l.index,
		end:   l.end,
	}
}

// Index() returns the current index into intLoop
func (l intLoop) Index() int {
	return l.index
}

// Return the current value of intLoop as string. Return an empty string when there are
// no more elements in the intLoop
func (l intLoop) Current() string {
	if l.index > l.end {
		return ""
	}
	return fmt.Sprintf("%d", l.index) // Returning string because of subGnrtr interface.
}

// return the next value of intLoop with done set to false, or an empty string
// and done set to true if there are no more elements.
func (l *intLoop) Next() (next string, done bool) {
	l.index += 1
	next = l.Current()
	if next == "" {
		done = true
	}
	return l.Current(), done
}

// Reset() resets the index of intLoop to its first element
func (l *intLoop) Reset() {
	l.index = l.begin
}

func (l intLoop) toArray() (a array) {
	a = array{
		index: l.index - l.begin,
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
