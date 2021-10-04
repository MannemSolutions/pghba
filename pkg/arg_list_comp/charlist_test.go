package arg_list_comp_test

import (
	"fmt"
	"github.com/mannemsolutions/pghba/pkg/arg_list_comp"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAlcCharList(t *testing.T) {
	var next string
	var results []string
	var done bool
	myAlc := "test[abc]"
	myArrayAlc := "test(a|b|c)"
	myFuncName := fmt.Sprintf("NewALC(\"%s\")", myAlc)
	myLoop := arg_list_comp.NewALC(myAlc)
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
	assert.Contains(t, results, "test1", "%s should return \"test1\"", myFuncName)
	assert.Contains(t, results, "test3", "%s should return \"test3\"", myFuncName)
}