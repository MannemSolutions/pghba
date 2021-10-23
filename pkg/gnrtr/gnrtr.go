package gnrtr

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Gnrtr struct {
	subGnrtrs map[int]subGnrtr
	current    string
	raw string
	allGnrtrs  subGnrtrs
}

func NewGnrtr(s string) (g Gnrtr) {
	g.raw = s
	for _, match := range reIntLoops.FindAllStringSubmatch(g.raw, -1) {
		sg, err := newIntLoop(match[0])
		if err != nil {
			panic(err)
		}
		placeholder := fmt.Sprintf("${%d}", len(g.allGnrtrs))
		g.raw = strings.Replace(g.raw, match[0], placeholder, 1)
		g.allGnrtrs = append(g.allGnrtrs, sg)
	}
	for _, match := range reCharLoops.FindAllStringSubmatch(g.raw, -1) {
		sg, err := newCharLoop(match[0])
		if err != nil {
			panic(err)
		}
		placeholder := fmt.Sprintf("${%d}", len(g.allGnrtrs))
		g.raw = strings.Replace(g.raw, match[0], placeholder, 1)
		g.allGnrtrs = append(g.allGnrtrs, sg)
	}
	for _, match := range reCharLists.FindAllStringSubmatch(g.raw, -1) {
		sg, err := newCharList(match[0])
		if err != nil {
			panic(err)
		}
		placeholder := fmt.Sprintf("${%d}", len(g.allGnrtrs))
		g.raw = strings.Replace(g.raw, match[0], placeholder, 1)
		g.allGnrtrs = append(g.allGnrtrs, sg)
	}
	// We do this in a loop, because array definitions might have array definitions inside them
	// The child arrays would be parsed and replaced on earlier passes, and parent arrays on later passes
	for {
		matches := reArrays.FindAllStringSubmatch(g.raw, -1)
		if matches == nil {
			break
		}
		for _,match := range matches {
			sg, err := newArray(match[0], g.allGnrtrs)
			if err != nil {
				panic(err)
			}
			placeholder := fmt.Sprintf("${%d}", len(g.allGnrtrs))
			g.raw = strings.Replace(g.raw, match[0], placeholder, 1)
			g.allGnrtrs = append(g.allGnrtrs, sg)
		}
	}
	g.buildSubGnrtrs()
	g.setCurrent()

	return g
}

func (g *Gnrtr) setCurrent() string {
	g.current = g.raw
	for i, sg := range g.subGnrtrs {
		g.current = strings.Replace(g.current, fmt.Sprintf("${%d}", i), sg.Current(), 1)
	}
	return g.current
}

func (g Gnrtr) Current() string {
	return g.current
}

func (g Gnrtr) String() (s string) {
	s = g.raw
	for i, sg := range g.subGnrtrs {
		s = strings.Replace(s, fmt.Sprintf("${%d}", i), sg.String(), 1)
	}
	return s
}

func (g *Gnrtr) buildSubGnrtrs() {
	g.subGnrtrs = make(map[int]subGnrtr, 0)
	reSubGenPlaceHolders := regexp.MustCompile(`\${(\d+)}`)
	matches := reSubGenPlaceHolders.FindAllStringSubmatch(g.raw, -1)
	for _, match := range matches {
		gnrtrId, err := strconv.Atoi(match[1])
		if err != nil {
			panic(fmt.Errorf("cannot convert %s to int", match[1]))
		}
		if gnrtrId >= len(g.allGnrtrs) {
			panic(fmt.Errorf("a placeholder references a non existing subGnrtr"))
		}
		g.subGnrtrs[gnrtrId] = g.allGnrtrs[gnrtrId]
	}
}

func (g *Gnrtr) Next() (string, bool) {
	for i := range g.subGnrtrs {
		if _, done := g.subGnrtrs[i].Next(); !done {
			// This one still can move to the next
			return g.setCurrent(), false
		}
		// SubGen at the end, lets start over
		g.subGnrtrs[i].Reset()
	}
	return g.Current(), true
}

func (g Gnrtr) ToArray() (a array) {
	return array{
		list: []string{g.raw},
		index: 0,
		allGnrtrs: g.allGnrtrs,
	}
}

func (g *Gnrtr) Reset() {
	for i := range g.allGnrtrs {
		g.allGnrtrs[i].Reset()
	}
}

func (g Gnrtr) ToList() []string {
	return subGnrtrToList(&g)
}

