package gnrtr_test

import (
	"github.com/mannemsolutions/pghba/pkg/gnrtr"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGnrtrArray(t *testing.T) {
	myGnrtr := "test(ing||user|er)"
	myFuncName := "NewGnrtr(\"test(ing||user|er)\")"
	myArray := gnrtr.NewGnrtr(myGnrtr)
	if !assert.NotNil(t, myArray, "%s should return an array iterator", myFuncName) {
		return
	}
	results := myArray.ToList()
	assert.Len(t, results, 4, "%s should return 4 elements", myFuncName)
	assert.Contains(t, results, "test", "%s should return \"test\"", myFuncName)
	assert.Contains(t, results, "testing", "%s should return \"testing\"", myFuncName)
}
