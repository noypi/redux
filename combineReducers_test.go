package redux_test

import (
	"log"
	"testing"

	"github.com/noypi/redux"
	assertpkg "github.com/stretchr/testify/assert"
)

func init() {
	redux.EnableDebugging()
}

func TestCombineReducers_x01(t *testing.T) {
	assert := assertpkg.New(t)

	type StateA struct {
		FieldOne string
		FieldTwo string
	}

	type PayloadOne struct {
		FieldOne string
	}
	type ActionA struct {
		Payload PayloadOne
	}

	type PayloadTwo struct {
		FieldTwo string
	}
	type ActionB struct {
		Payload PayloadTwo
	}

	reducer1 := redux.CombineReducers(redux.ReducerMap{
		"FieldOne": func(state interface{}, action ActionA) interface{} {
			return redux.Merge(state, action.Payload)
		},
		"FieldTwo": func(state interface{}, action ActionB) interface{} {
			return redux.Merge(state, action.Payload)
		},
	})

	newState := reducer1(StateA{}, ActionA{Payload: PayloadOne{"New Field One"}})
	v, ok := newState.(StateA)
	assert.True(ok)
	assert.Equal("New Field One", v.FieldOne, "newState=%v", newState)
	assert.Equal("", v.FieldTwo)

	reducer2 := redux.CombineReducers(redux.ReducerMap{
		"FieldTwo": func(state interface{}, action ActionB) interface{} {
			return redux.Merge(state, action.Payload)
		},
	})

	log.Println("------- trying reducer2")

	newState = reducer2(newState, ActionB{Payload: PayloadTwo{"Second try"}})
	v, ok = newState.(StateA)
	assert.True(ok)
	assert.Equal("New Field One", v.FieldOne, "newState=%v", newState)
	assert.Equal("Second try", v.FieldTwo)
}

func TestCombineReducers_withNewProperty(t *testing.T) {
	assert := assertpkg.New(t)

	type StateA struct {
		FieldOne string
	}

	type PayloadOne struct {
		FieldOne        string
		SomeNewProperty string
	}
	type PayloadB struct {
		NewPropStruct struct {
			NewA string
			NewB string
		}
	}
	type ActionA struct {
		Payload PayloadOne
	}

	reducer1 := redux.CombineReducers(
		func(state interface{}, action ActionA) interface{} {
			return redux.Merge(state, action.Payload)
		},
		func(state interface{}, action ActionA) interface{} {
			action.Payload.SomeNewProperty += " Appended Value"
			redux.DBG("trying to merge from SomeNewProperty")
			return redux.Merge(state, action.Payload)
		},
	)

	newState := reducer1(StateA{}, ActionA{PayloadOne{"New Field One", "unknown"}})
	v, ok := newState.(redux.ReducerResult)
	assert.True(ok)
	assert.Equal(2, len(v))
	assert.Equal("New Field One", v.Get("FieldOne"), "newState=%v", newState)
	assert.Equal("unknown Appended Value", v.Get("SomeNewProperty"))

}

func TestCombineReducers_withSubField(t *testing.T) {
	assert := assertpkg.New(t)

	type SubField struct {
		Phone   string
		Address string
	}
	type StateA struct {
		FieldOne string
		Sub      SubField
	}

	type PayloadOne struct {
		FieldOne        string
		SomeNewProperty string
	}

	type PayloadUpdatePhone struct {
		Sub struct {
			Phone string
		}
	}
	type PayloadUpdateSub struct {
		Sub struct {
			Phone   string
			Address string
		}
	}
	type ActionA struct {
		Payload PayloadOne
	}
	type ActionUpdatePhone struct {
		Payload PayloadUpdatePhone
	}

	type ActionUpdateAddress struct {
		Payload PayloadUpdateSub
	}

	subReducerAddress := func(state string, action ActionUpdateAddress) interface{} {
		state += " Before Address Append"
		action.Payload.Sub.Address += " Appended Address"
		return redux.Merge(state, action.Payload)
	}

	subReducerPhone := func(state string, action ActionUpdatePhone) interface{} {
		state += " Before Phone Append"
		action.Payload.Sub.Phone += " Appended Phone"
		return redux.Merge(state, action.Payload)
	}

	reducer1 := redux.CombineReducers(redux.ReducerMap{
		"FieldOne": func(state interface{}, action ActionA) interface{} {
			return redux.Merge(state, action.Payload)
		},
		"Sub": redux.CombineReducers(redux.ReducerMap{
			"Phone":   subReducerPhone,
			"Address": subReducerAddress,
		}),
	})

	state0 := reducer1(StateA{}, ActionA{PayloadOne{FieldOne: "my field one"}})

	log.Println("-------------- test1")
	updatePhone := PayloadUpdatePhone{}
	updatePhone.Sub.Phone = "phone1234"
	newState := reducer1(state0, ActionUpdatePhone{updatePhone})
	log.Printf("newstate %t: %v", newState, newState)
	v, ok := newState.(StateA)
	log.Println("v=", v)
	assert.True(ok, "newState=%t: %v", newState, newState)
	assert.Equal("my field one", v.FieldOne)
	assert.Equal("", v.Sub.Address)
	assert.Equal("phone1234 Appended Phone", v.Sub.Phone)

	log.Println("-------------- test2")
	actionAddress := ActionUpdateAddress{}
	actionAddress.Payload.Sub.Address = "some address sub"

	actionPhone := ActionUpdatePhone{}
	actionPhone.Payload.Sub.Phone = "sub phone1234"

	newState = reducer1(StateA{}, actionAddress)
	v, ok = newState.(StateA)
	assert.True(ok, "newState=%t: %v", newState, newState)
	assert.Equal("some address sub Appended Address", v.Sub.Address)

	newState = reducer1(StateA{}, actionPhone)
	v, ok = newState.(StateA)
	assert.True(ok, "newState=%t: %v", newState, newState)
	assert.Equal("sub phone1234 Appended Phone", v.Sub.Phone)
}
