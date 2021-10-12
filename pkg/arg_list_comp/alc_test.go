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
	myAlcDef := "test(ing|{1..3})"
	myFuncName := fmt.Sprintf("NewALC(\"%s\")", myAlcDef)
	myAlc := NewALC(myAlcDef)
	assert.NotNil(t, myAlc, "%s should return an iterator", myFuncName)
	assert.Equal(t, myAlcDef, myAlc.String(), "%s.String() should be \"%s\"", myFuncName, myAlcDef)
	mySortedArray := SortedArray(myAlc)

	for {
		next, done = mySortedArray.Next()
		if done {
			break
		}
		results = append(results, next)
	}
	assert.Len(t, results, 4, "%s should return 4 elements", myFuncName)
	assert.Contains(t, results, "testing", "%s should return \"testing\"", myFuncName)
	assert.Contains(t, results, "test2", "%s should return \"test2\"", myFuncName)
}
