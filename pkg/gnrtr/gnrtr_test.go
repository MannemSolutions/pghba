package gnrtr

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGnrtrLoopSub(t *testing.T) {
	myGnrtrDef := "test(ing|{1..3})"
	myFuncName := fmt.Sprintf("NewGnrtr(\"%s\")", myGnrtrDef)
	myGnrtr := NewGnrtr(myGnrtrDef)
	assert.NotNil(t, myGnrtr, "%s should return an iterator", myFuncName)
	assert.Equal(t, myGnrtrDef, myGnrtr.String(), "%s.String() should be \"%s\"", myFuncName, myGnrtrDef)

	results := myGnrtr.ToList()
	assert.Len(t, results, 4, "%s should return 4 elements", myFuncName)
	assert.Contains(t, results, "testing", "%s should return \"testing\"", myFuncName)
	assert.Contains(t, results, "test2", "%s should return \"test2\"", myFuncName)
}
