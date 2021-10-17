package gnrtr

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestALCIntLoop(t *testing.T) {
	var next string
	var results []string
	var done bool
	myAlc := "test{1..3}"
	myArrayAlc := "test(1|2|3)"
	myFuncName := fmt.Sprintf("NewGnrtr(\"%s\")", myAlc)
	myLoop := NewGnrtr(myAlc)
	if !assert.NotNil(t, myLoop, "%s should return a loop iterator", myFuncName) {
		return
	}
	assert.Equal(t, myLoop.String(), myAlc, "%s.String() should be \"%s\"", myFuncName, myAlc)

	assert.Equal(t, myArrayAlc, myLoop.ToArray().String(),
		"%s.ToArray().String() should be \"%s\"", myFuncName, myArrayAlc)
	for {
		next, done = myLoop.Next()
		if done {
			break
		}
		results = append(results, next)
	}
	assert.Len(t, results, 3, "%s should return 3 elements", myFuncName)
	assert.Contains(t, results, "test1", "%s should return \"test1\"", myFuncName)
	assert.Contains(t, results, "test3", "%s should return \"test3\"", myFuncName)
}
