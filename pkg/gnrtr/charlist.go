package gnrtr

import (
	"fmt"
	"strings"
)

// charList is a slice of byte and has a stored index of type int.
// charList stores the list of characters to use for expansion in the format
// from regular expressions. This includes ranges. It is described by the
// regular expression `\[([^]]+)]`
// TODO Question: does the *regex.Regexp type automagically have a .String()
// method? Answer: it does! (should have looked at the link to the definition
// popup. Still, why don't you have to import the regex package here?
type charList struct {
	list  []byte
	index int
}

// Returns a new charList 'cl' containing all the characters described by
// string 's' as bytes in an array
func newCharList(s string) (cl *charList, err error) {
	match := reCharList.FindStringSubmatch(s)
	if match == nil { // verify input first
		return nil, fmt.Errorf("invalid input to newCharList (should have form %s)", reCharList.String())
	}
	if strings.HasPrefix(match[1], "^") {
		return nil, fmt.Errorf("cannot make an iterator of a negative character list (starting with ^)")
	}
	cl = &charList{ // Initialize the charList
		index: 0,
	}
	chars := match[1] // This selects the first substring match
	for i := 0; i < len(chars); i++ {
		if i < len(chars)-1 {
			if chars[i+1] == '-' { // detect character ranges
				start := chars[i]
				end := chars[i+2]
				if end < start {
					return nil, fmt.Errorf("could not parse %s, %s should be before %s", s, string(start),
						string(end))
				} // add each individual character in the range to the character list
				for char := start; char <= end; char++ {
					cl.list = append(cl.list, char)
				}
				i += 2
				continue
			}
		}
		cl.list = append(cl.list, chars[i])
	}
	return cl, nil
}

// All implementations of methods for charList follow below and are self-describing afaics

func (cl charList) clone() subGnrtr {
	return &charList{
		list:  cl.list,
		index: cl.index,
	}
}

func (cl charList) Index() int {
	return cl.index
}

func (cl charList) Current() string {
	if cl.index >= len(cl.list) {
		return ""
	}
	return string(cl.list[cl.index])
}

func (cl *charList) Next() (next string, done bool) {
	cl.index += 1
	next = cl.Current()
	if next == "" {
		done = true
	}
	return next, done
}

func (cl *charList) Reset() {
	cl.index = 0
}

func (cl *charList) toArray() (a array) {
	a = array{
		index: cl.index,
	}
	for i := 0; i < len(cl.list); i++ {
		a.list = append(a.list, string([]byte{cl.list[i]}))
	}
	return a
}

func (cl *charList) String() (s string) {
	return fmt.Sprintf("[%s]", string(cl.list))
}

func (cl charList) ToList() []string {
	return subGnrtrToList(&cl)
}
