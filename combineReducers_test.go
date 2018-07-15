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

	reducer1 := redux.CombineReducers([]redux.FieldReducer{
		{Name: "FieldOne", Reducer: func(state, action interface{}) interface{} {
			if _, ok := action.(ActionA); !ok {
				return state
			}
			return redux.Merge(state, action.(ActionA).Payload)
		}},
		{Name: "FieldTwo", Reducer: func(state, action interface{}) interface{} {
			if _, ok := action.(ActionB); !ok {
				return state
			}
			return redux.Merge(state, action.(ActionB).Payload)
		}},
	})

	newState := reducer1(StateA{}, ActionA{Payload: PayloadOne{"New Field One"}})
	v, ok := newState.(StateA)
	assert.True(ok)
	assert.Equal("New Field One", v.FieldOne, "newState=%v", newState)
	assert.Equal("", v.FieldTwo)

	reducer2 := redux.CombineReducers([]redux.FieldReducer{
		{Name: "FieldTwo", Reducer: func(state, action interface{}) interface{} {
			return redux.Merge(state, action.(ActionB).Payload)
		}},
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

	reducer1 := redux.CombineReducers([]redux.FieldReducer{
		{Name: "FieldOne", Reducer: func(state, action interface{}) interface{} {
			return redux.Merge(state, action.(ActionA).Payload)
		}},
		{Name: "SomeNewProperty", Reducer: func(state, action interface{}) interface{} {
			extendedAction := action.(ActionA)
			extendedAction.Payload.SomeNewProperty += " Appended Value"
			return redux.Merge(state, extendedAction.Payload)
		}},
	})

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

	type ActionSub struct {
		Payload PayloadUpdateSub
	}

	subReducer := func(state, action interface{}) interface{} {
		extendedAction, ok := action.(ActionSub)
		if !ok {
			return state
		}
		extendedAction.Payload.Sub.Phone += " Appended Value"
		extendedAction.Payload.Sub.Address += " Appended Value"
		return redux.Merge(state, extendedAction.Payload)
	}

	subReducerPhone := func(state, action interface{}) interface{} {
		extendedAction, ok := action.(ActionUpdatePhone)
		if !ok {
			return state
		}
		extendedAction.Payload.Sub.Phone += " Appended Phone"
		return redux.Merge(state, extendedAction.Payload)
	}

	reducer1 := redux.CombineReducers([]redux.FieldReducer{
		{Name: "FieldOne", Reducer: func(state, action interface{}) interface{} {
			if _, ok := action.(ActionA); !ok {
				return state
			}
			return redux.Merge(state, action.(ActionA).Payload)
		}},
		{Name: "Sub", Reducer: func(state, action interface{}) interface{} {
			if _, ok := action.(ActionUpdatePhone); !ok {
				return state
			}
			return redux.Merge(state, action.(ActionUpdatePhone).Payload)
		}},
		{Name: "Sub", Reducer: subReducer},
		{Name: "Sub", Reducer: redux.CombineReducers([]redux.FieldReducer{
			{Name: "Phone", Reducer: subReducerPhone},
		})},
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
	actionSub := ActionSub{}
	actionSub.Payload.Sub.Phone = "sub phone1234"
	actionSub.Payload.Sub.Address = "some address sub"
	newState = reducer1(StateA{}, actionSub)
	v, ok = newState.(StateA)
	assert.True(ok, "newState=%t: %v", newState, newState)
	assert.Equal("sub phone1234 Appended Value", v.Sub.Phone)
	assert.Equal("some address sub Appended Value", v.Sub.Address)
}
