package arg_list_comp

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestALCLoopSub(t *testing.T) {
	var next string
	var results []string
	var done bool
	myAlc := "test(ing|{1..3})"
	myFuncName := fmt.Sprintf("NewALC(\"%s\")", myAlc)
	myLoop := NewALC(myAlc)
	assert.NotNil(t, myLoop, "%s should return an iterator", myFuncName)
	assert.Equal(t, myAlc, myLoop.String(), "%s.String() should be \"%s\"", myFuncName, myAlc)

	for {
		next, done = myLoop.Next()
		if done {
			break
		}
		results = append(results, next)
	}
	assert.Len(t, results, 4, "%s should return 4 elements", myFuncName)
	assert.Contains(t, results, "testing", "%s should return \"testing\"", myFuncName)
	assert.Contains(t, results, "test2", "%s should return \"test2\"", myFuncName)
}
