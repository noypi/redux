package redux_test

import (
	"testing"

	. "github.com/noypi/redux"
	assertpkg "github.com/stretchr/testify/assert"
)

func TestCombinReducer_withReducerArr(t *testing.T) {
	assert := assertpkg.New(t)

	type User struct {
		Name  string
		Phone string
	}

	type UserUpdate struct {
		Name  *string
		Phone *string
	}

	type StateA struct {
		User User
	}

	type ActionUpdateUser struct {
		Payload UserUpdate
	}

	type ActionUpdateUser2 struct {
		Payload map[string]interface{}
	}

	fnHandleUpdatePhone := func(phone, action string) string {
		return action
	}

	fnHandleUpdateUser := func(user User, action ActionUpdateUser) (out interface{}) {
		return Merge(user, action.Payload)
	}

	reducer := CombineReducers(ReducerMap{
		"User": CombineReducers(
			fnHandleUpdateUser,
			ReducerMap{
				"Phone": fnHandleUpdatePhone,
			}),
	})

	newPhone := "4567"
	state := reducer(StateA{
		User{Name: "some name", Phone: "phone1234"}},
		ActionUpdateUser{Payload: UserUpdate{Phone: &newPhone}})

	vstate, ok := state.(StateA)
	assert.True(ok)
	assert.Equal("some name", vstate.User.Name, "vstate=%v", vstate)
	assert.Equal("4567", vstate.User.Phone)
}
