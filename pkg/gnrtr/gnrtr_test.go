package gnrtr

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGnrtrLoopSub(t *testing.T) {
	var next string
	var results []string
	var done bool
	myGnrtrDef := "test(ing|{1..3})"
	myFuncName := fmt.Sprintf("NewGnrtr(\"%s\")", myGnrtrDef)
	myGnrtr := NewGnrtr(myGnrtrDef)
	assert.NotNil(t, myGnrtr, "%s should return an iterator", myFuncName)
	assert.Equal(t, myGnrtrDef, myGnrtr.String(), "%s.String() should be \"%s\"", myFuncName, myGnrtrDef)

	for {
		next, done = myGnrtr.Next()
		if done {
			break
		}
		results = append(results, next)
	}
	assert.Len(t, results, 4, "%s should return 4 elements", myFuncName)
	assert.Contains(t, results, "testing", "%s should return \"testing\"", myFuncName)
	assert.Contains(t, results, "test2", "%s should return \"test2\"", myFuncName)
}
