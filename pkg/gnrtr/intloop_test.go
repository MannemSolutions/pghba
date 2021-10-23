package gnrtr

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGnrtrIntLoop(t *testing.T) {
	myGnrtr := "test{1..3}"
	//myArrayGnrtr := "test(1|2|3)"
	myFuncName := fmt.Sprintf("NewGnrtr(\"%s\")", myGnrtr)
	myLoop := NewGnrtr(myGnrtr)
	if !assert.NotNil(t, myLoop, "%s should return a loop iterator", myFuncName) {
		return
	}
	assert.Equal(t, myLoop.String(), myGnrtr, "%s.String() should be \"%s\"", myFuncName, myGnrtr)

	//assert.Equal(t, myArrayGnrtr, myLoop.ToArray().String(),
	//	"%s.ToArray().String() should be \"%s\"", myFuncName, myArrayGnrtr)
	results := myLoop.ToList()
	assert.Len(t, results, 3, "%s should return 3 elements", myFuncName)
	assert.Contains(t, results, "test1", "%s should return \"test1\"", myFuncName)
	assert.Contains(t, results, "test3", "%s should return \"test3\"", myFuncName)
}
