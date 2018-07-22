package redux

import (
	"reflect"
)

func BindActionCreator(dispatch func(action interface{}), actionCreatorFunc interface{}) (newActionCreator func(...interface{}) interface{}) {
	tfn := reflect.TypeOf(actionCreatorFunc)
	assertValidActionCreator(tfn)
	vfn := reflect.ValueOf(actionCreatorFunc)

	return func(args ...interface{}) (newAction interface{}) {
		vargs := make([]reflect.Value, len(args))
		for i, arg := range args {
			vargs[i] = reflect.ValueOf(arg)
		}
		vres := vfn.Call(vargs)
		newAction = vres[0].Interface()
		dispatch(newAction)
		return
	}
}

func assertValidActionCreator(tfn reflect.Type) {
	bValid := (tfn.Kind() == reflect.Func) &&
		(1 <= tfn.NumOut())
	if !bValid {
		panic("invalid actionCreatorFunc. must be a func with at least one return value.")
	}
}
