package gnrtr

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
)

type Gnrtr struct {
	index      int
	current    string
	currentRaw string
	allGnrtrs  subGnrtrs
}

func NewGnrtr(s string) (g *Gnrtr) {
	g = &Gnrtr{
		currentRaw: s,
	}
	for _, match := range reIntLoops.FindAllStringSubmatch(g.currentRaw, -1) {
		sg, err := newIntLoop(match[0])
		if err != nil {
			panic(err)
		}
		placeholder := fmt.Sprintf("${%d}", len(g.allGnrtrs))
		g.currentRaw = strings.Replace(g.currentRaw, match[0], placeholder, 1)
		g.allGnrtrs = append(g.allGnrtrs, sg)
	}
	for _, match := range reCharLoops.FindAllStringSubmatch(g.currentRaw, -1) {
		sg, err := newCharLoop(match[0])
		if err != nil {
			panic(err)
		}
		placeholder := fmt.Sprintf("${%d}", len(g.allGnrtrs))
		g.currentRaw = strings.Replace(g.currentRaw, match[0], placeholder, 1)
		g.allGnrtrs = append(g.allGnrtrs, sg)
	}
	for _, match := range reCharLists.FindAllStringSubmatch(g.currentRaw, -1) {
		sg, err := newCharList(match[0])
		if err != nil {
			panic(err)
		}
		placeholder := fmt.Sprintf("${%d}", len(g.allGnrtrs))
		g.currentRaw = strings.Replace(g.currentRaw, match[0], placeholder, 1)
		g.allGnrtrs = append(g.allGnrtrs, sg)
	}
	// We do this in a loop, because array definitions might have array definitions inside them
	// The child arrays would be parsed and replaced on earlier passes, and parent arrays on later passes
	for {
		matches := reArrays.FindAllStringSubmatch(g.currentRaw, -1)
		if matches == nil {
			break
		}
		for _, match := range matches {
			sg, err := newArray(match[0], g.allGnrtrs)
			if err != nil {
				panic(err)
			}
			placeholder := fmt.Sprintf("${%d}", len(g.allGnrtrs))
			g.currentRaw = strings.Replace(g.currentRaw, match[0], placeholder, 1)
			g.allGnrtrs = append(g.allGnrtrs, sg)
		}
	}
	g.Reset()

	return g
}

func (g Gnrtr) clone() subGnrtr {
	clone := &Gnrtr{
		index: g.index,
		currentRaw:   g.currentRaw,
	}
	for _, sg := range g.allGnrtrs {
		clone.allGnrtrs = append(clone.allGnrtrs, sg.clone())
	}
	return clone
}

func (g *Gnrtr) setCurrent() string {
	g.current = g.currentRaw

	for {
		matches := rePlaceholders.FindAllStringSubmatch(g.current, -1)
		if len(matches) == 0 {
			break
		}
		newValue := g.current
		for _, match := range matches {
			i, err := strconv.Atoi(match[1])
			if err != nil {
				// This should not be possible!!!
				log.Panicf("cannot parse %s to int (%s) in (g *Gnrtr).setCurrent()", match[1], err.Error())
			}
			sg := g.allGnrtrs[i]
			g.current = strings.Replace(g.current, match[0], sg.Current(), 1)
		}
		if g.current == newValue {
			log.Panicf("We seem to be in a deadloop here. File a bug with the exact call to the executable.")
		}
	}
	return g.current
}

func (g Gnrtr) Index() int {
	return g.index
}

func (g Gnrtr) Current() string {
	return g.current
}

func (g Gnrtr) String() (s string) {
	s = g.currentRaw
	for i, sg := range g.subGnrtrs() {
		s = strings.Replace(s, fmt.Sprintf("${%d}", i), sg.String(), 1)
	}
	return s
}

func (g Gnrtr) subGnrtrs() (sg subGnrtrs) {
	reSubGenPlaceHolders := regexp.MustCompile(`\${(\d+)}`)
	matches := reSubGenPlaceHolders.FindAllStringSubmatch(g.currentRaw, -1)
	for _, match := range matches {
		gnrtrId, err := strconv.Atoi(match[1])
		if err != nil {
			panic(fmt.Errorf("cannot convert %s to int", match[1]))
		}
		if gnrtrId >= len(g.allGnrtrs) {
			panic(fmt.Errorf("a placeholder references a non existing subGnrtr"))
		}
		sg = append(sg, g.allGnrtrs[gnrtrId])
	}
	return sg
}

func (g *Gnrtr) buildSubGnrtrs() (sg subGnrtrs) {
	reSubGenPlaceHolders := regexp.MustCompile(`\${(\d+)}`)
	matches := reSubGenPlaceHolders.FindAllStringSubmatch(g.currentRaw, -1)
	for _, match := range matches {
		gnrtrId, err := strconv.Atoi(match[1])
		if err != nil {
			panic(fmt.Errorf("cannot convert %s to int", match[1]))
		}
		if gnrtrId >= len(g.allGnrtrs) {
			panic(fmt.Errorf("a placeholder references a non existing subGnrtr"))
		}
		sg = append(sg, g.allGnrtrs[gnrtrId])
	}
	return sg
}

func (g *Gnrtr) Next() (string, bool) {
	g.index += 1
	if g.index == 0 {
		return g.Current(), false
	}
	sgs := g.subGnrtrs()
	for i := range sgs {
		if _, done := sgs[i].Next(); !done {
			// This one still can move to the next
			return g.setCurrent(), false
		}
		// SubGen at the end, lets start over
		sgs[i].Reset()
	}
	return g.Current(), true
}

func (g Gnrtr) toArray() (a array) {
	a = array{
		list:      g.ToList(),
		index:     0,
		allGnrtrs: g.allGnrtrs,
	}
	a.setCurrent()
	return a
}

func (g *Gnrtr) Reset() {
	g.index = -1
	for i := range g.allGnrtrs {
		g.allGnrtrs[i].Reset()
	}
	g.buildSubGnrtrs()
	g.setCurrent()
}

func (g *Gnrtr) ToList() (list []string) {
	list = subGnrtrToList(g)
	return list
}
