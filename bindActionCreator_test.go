package redux_test

import (
	"testing"

	. "github.com/noypi/redux"
	assertpkg "github.com/stretchr/testify/assert"
)

func TestBindCreator_x01(t *testing.T) {
	assert := assertpkg.New(t)

	fnActionCreator := func(param string) (newAction string) {
		return param + " appended newAction"
	}

	var result1 string
	fnDispatch := func(action interface{}) {
		result1 = action.(string)
	}

	fn := BindActionCreator(fnDispatch, fnActionCreator)

	assert.Equal("param1 appended newAction", fn("param1"))
	assert.Equal("param1 appended newAction", result1)
}
