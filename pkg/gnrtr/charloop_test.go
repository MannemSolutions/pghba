package gnrtr

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGnrtrCharLoop(t *testing.T) {
	myGnrtr := "test_{a..c}"
	//myArrayGnrtr := "test_(a|b|c)"
	myFuncName := fmt.Sprintf("NewGnrtr(\"%s\")", myGnrtr)
	myLoop := NewGnrtr(myGnrtr)
	if !assert.NotNil(t, myLoop, "%s should return a loop iterator", myFuncName) {
		return
	}
	assert.Equal(t, myGnrtr, myLoop.String(), "%s.String() should be \"%s\"", myFuncName, myGnrtr)

	//assert.Equal(t, myArrayGnrtr, myLoop.ToArray().String(),
	//	"%s.ToArray().String() should be \"%s\"", myFuncName, myArrayGnrtr)
	results := myLoop.ToList()
	assert.Len(t, results, 3, "%s should return 3 elements", myFuncName)
	assert.Contains(t, results, "test_a", "%s should return \"test_a\"", myFuncName)
	assert.Contains(t, results, "test_c", "%s should return \"test_c\"", myFuncName)
}
