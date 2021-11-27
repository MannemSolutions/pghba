package gnrtr

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGnrtrCharList(t *testing.T) {
	myGnrtr := "test_[-ac-e-]"
	myGnrtrStr := "test_[-acde-]"
	myArrayGnrtr := "(test_-|test_a|test_c|test_d|test_e|test_-)"
	myFuncName := fmt.Sprintf("NewGnrtr(\"%s\")", myGnrtr)
	myLoop := NewGnrtr(myGnrtr)
	if !assert.NotNil(t, myLoop, "%s should return a loop iterator", myFuncName) {
		return
	}
	assert.Equal(t, myGnrtrStr, myLoop.String(), "%s.String() should be \"%s\"", myFuncName, myGnrtrStr)
	assert.Equal(t, myArrayGnrtr, myLoop.toArray().String(),
		"%s.ToArray().String() should be \"%s\"", myFuncName, myArrayGnrtr)
	results := myLoop.ToList()
	assert.Len(t, results, 6, "%s should return 4 elements", myFuncName)
	assert.Contains(t, results, "test_-", "%s should return \"test_-\"", myFuncName)
	assert.Contains(t, results, "test_a", "%s should return \"test_a\"", myFuncName)
	assert.Contains(t, results, "test_c", "%s should return \"test_c\"", myFuncName)
}
