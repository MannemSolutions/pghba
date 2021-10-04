package arg_list_comp_test

import (
	"github.com/mannemsolutions/pghba/pkg/arg_list_comp"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestALCArray(t *testing.T) {
	var next string
	var results []string
	var done bool
	myAlc := "test(ing||user|er)"
	myFuncName := "NewALC(\"test(ing||user|er)\")"
	myArray := arg_list_comp.NewALC(myAlc)
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
	assert.Len(t, results, 2, "%s should return 4 elements", myFuncName)
	assert.Contains(t, results, "test", "%s should return \"test\"", myFuncName)
	assert.Contains(t, results, "testing", "%s should return \"testing\"", myFuncName)
}