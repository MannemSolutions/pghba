package gnrtr

type subGnrtr interface {
	Current() string
	Next() (string, bool)
	ToArray() (a array)
	Reset()
	ToList() []string
	String() string
	Index() int
	Clone() subGnrtr
}

type subGnrtrs []subGnrtr

// SortedArray creates an array with sorted unique elements
func subGnrtrToList(g subGnrtr) (l []string) {
	// Clone so we can reset
	clone := g.Clone()
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
