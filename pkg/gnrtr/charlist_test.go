package gnrtr_test

import (
	"fmt"
	"github.com/mannemsolutions/pghba/pkg/gnrtr"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAlcCharList(t *testing.T) {
	var next string
	var results []string
	var done bool
	myAlc := "test_[abc]"
	myArrayAlc := "test_(a|b|c)"
	myFuncName := fmt.Sprintf("NewGnrtr(\"%s\")", myAlc)
	myLoop := gnrtr.NewGnrtr(myAlc)
	if !assert.NotNil(t, myLoop, "%s should return a loop iterator", myFuncName) {
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
