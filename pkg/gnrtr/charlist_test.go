package gnrtr_test

import (
	"fmt"
	"github.com/mannemsolutions/pghba/pkg/gnrtr"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGnrtrCharList(t *testing.T) {
	var next string
	var results []string
	var done bool
	myGnrtr := "test_[ac-e]"
	myArrayGnrtr := "test_(a|c|d|e)"
	myFuncName := fmt.Sprintf("NewGnrtr(\"%s\")", myGnrtr)
	myLoop := gnrtr.NewGnrtr(myGnrtr)
	if !assert.NotNil(t, myLoop, "%s should return a loop iterator", myFuncName) {
		return
	}
	assert.Equal(t, myGnrtr, myLoop.String(), "%s.String() should be \"%s\"", myFuncName, myGnrtr)
	assert.Equal(t, myArrayGnrtr, myLoop.ToArray().String(),
		"%s.ToArray().String() should be \"%s\"", myFuncName, myArrayGnrtr)
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
