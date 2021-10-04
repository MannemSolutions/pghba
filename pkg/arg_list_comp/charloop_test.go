package arg_list_comp

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestALCCharLoop(t *testing.T) {
	var next string
	var results []string
	var done bool
	myAlc := "test_{a..c}"
	myArrayAlc := "test(a|b|c)"
	myFuncName := fmt.Sprintf("NewALC(\"%s\")", myAlc)
	myLoop := NewALC(myAlc)
	if ! assert.NotNil(t, myLoop, "%s should return a loop iterator", myFuncName) {
		return
	}
	assert.Equal(t, myAlc, myLoop.String(), "%s.String() should be \"%s\"", myFuncName, myAlc)

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
	assert.Contains(t, results, "test_a", "%s should return \"test_a\"", myFuncName)
	assert.Contains(t, results, "test_c", "%s should return \"test_c\"", myFuncName)
}