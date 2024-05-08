package hba

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNextOption(t *testing.T) {
	myOption1 := " ldapuser\t = \"my text\" key2=someother"
	key, value, rest := nextOption(myOption1)
	if !assert.Equal(t, "ldapuser", key, "Expecting next key for %s to be %s but received %s", myOption1, "ldapuser", key) {
		return
	}
	if !assert.Equal(t, "my text", value, "Expecting next value for %s to be %s but received %s", myOption1, "my text", value) {
		return
	}
	if !assert.Equal(t, "key2=someother", rest, "Expecting remainder for %s to be %s but received %s", myOption1, "key2=someother", rest) {
		return
	}
}
