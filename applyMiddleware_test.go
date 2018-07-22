package redux_test

import (
	"fmt"

	. "github.com/noypi/redux"
	assertpkg "github.com/stretchr/testify/assert"

	"testing"
)

func TestApplyMiddleware_x01(t *testing.T) {
	assert := assertpkg.New(t)

	var results []interface{}

	fnNewMid := func(i int) MiddlewareFunc {
		return func(fn1 GetStateFunc, fn2 DispatchFunc) DispatchFunc {
			return func(action interface{}) interface{} {
				DBG("fnMid i=", i)
				results = append(results, fmt.Sprintf("fnMid%d called", i))
				return fn2(action)
			}
		}
	}

	fns := []MiddlewareFunc{}
	for i := 0; i < 3; i++ {
		fns = append(fns, fnNewMid(i))
	}

	fnGetState := func() interface{} {
		return nil
	}

	fnDispatch := func(action interface{}) interface{} {
		return action
	}

	createStore := ApplyMiddleware(fnGetState, fnDispatch, fns...)
	store := createStore()
	store.Dispatch("test dispatch")

	assert.Equal(3, len(results))
	for i, res := range results {
		assert.Equal(fmt.Sprintf("fnMid%d called", i), res)
	}
}
