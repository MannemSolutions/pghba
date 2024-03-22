package gnrtr

// subGnrtr interface
type subGnrtr interface {
	Current() string      // Return the character in the current position within the generator pattern as a string
	Next() (string, bool) // The next character within the pattern and a boolean to indicate if the last element is reached
	toArray() (a array)
	Reset() // Reset the index, useful for walking a slice without disturbing the original
	ToList() []string
	String() string // Return the whole slice of values as string
	Index() int     // Return the current index location in the slice
	clone() subGnrtr
}

type subGnrtrs []subGnrtr // basically a slice of any type T that implements the subGnrtr interface.

// subGnrtrToList returns a slice of strings by iterating over the values
func subGnrtrToList(g subGnrtr) (l []string) {
	// Clone so we can reset
	clone := g.clone()
	clone.Reset()
	//Make unique
	for {
		next, done := clone.Next()
		if done {
			break
		}
		l = append(l, next)
	}
	return l
}
