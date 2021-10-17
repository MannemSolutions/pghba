package gnrtr_test

import (
	"github.com/mannemsolutions/pghba/pkg/gnrtr"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestALCArray(t *testing.T) {
	var next string
	var results []string
	var done bool
	myAlc := "test(ing||user|er)"
	myFuncName := "NewGnrtr(\"test(ing||user|er)\")"
	myArray := gnrtr.NewGnrtr(myAlc)
	if ! assert.NotNil(t, myArray, "%s should return an array iterator", myFuncName) {
		return
	}
	for {
		next, done = myArray.Next()
		if done {
			break
		}
		results = append(results, next)
	}
	assert.Len(t, results, 4, "%s should return 4 elements", myFuncName)
	assert.Contains(t, results, "test", "%s should return \"test\"", myFuncName)
	assert.Contains(t, results, "testing", "%s should return \"testing\"", myFuncName)
}