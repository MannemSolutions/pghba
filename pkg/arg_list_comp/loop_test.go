package arg_list_comp

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestALCIntLoop(t *testing.T) {
	var next string
	var results []string
	var done bool
	myAlc := "test{1..3}"
	myFuncName := fmt.Sprintf("NewALC(\"%s\")", myAlc)
	myLoop := NewALC(myAlc)
	if ! assert.NotNil(t, myLoop, "%s should return a loop iterator", myFuncName) {
		return
	}
	assert.Equal(t, myLoop.String(), myAlc, "%s.String() should be \"%s\"", myFuncName, myAlc)

	assert.Equal(t, myLoop.ToArray().String(), "test(1|2|3)",
		"%s.ToArray().String() should be \"test(1|2|3)\"", myFuncName)
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

func TestALCCharLoop(t *testing.T) {
	var next string
	var results []string
	var done bool
	myAlc := "test{a..c}"
	myFuncName := fmt.Sprintf("NewALC(\"%s\")", myAlc)
	myLoop := NewALC(myAlc)
	if assert.NotNil(t, myLoop, "%s should return a loop iterator", myFuncName) {
		return
	}
	assert.Equal(t, myLoop.String(), myAlc, "%s.String() should be \"test{1..3}\"", myFuncName)

	assert.Equal(t, myLoop.ToArray().String(), "test(a|b|c)",
		"%s.ToArray().String() should be \"test(1|2|3)\"", myFuncName)
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


func TestALCLoopSub(t *testing.T) {
	var next string
	var results []string
	var done bool
	myAlc := "test(ing|{1..3})"
	myFuncName := fmt.Sprintf("NewALC(\"%s\")", myAlc)
	myLoop := NewALC(myAlc)
	assert.IsType(t, &loop{}, myLoop, "%s should return an iterator", myFuncName)
	assert.Equal(t, myLoop.String(), myAlc, "%s.String() should be \"%s\"", myFuncName, myAlc)

	assert.Equal(t, myLoop.ToArray().String(), "test(1|2|3)",
		"%s.ToArray().String() should be \"test(1|2|3)\"", myFuncName)
	for {
		next, done = myLoop.Next()
		if done {
			break
		}
		results = append(results, next)
	}
	assert.Len(t, results, 3, "%s should return 4 elements", myFuncName)
	assert.Contains(t, results, "testing", "%s should return \"testing\"", myFuncName)
	assert.Contains(t, results, "test2", "%s should return \"test2\"", myFuncName)
}
