package redux

import (
	"sort"
)

type GetStateFunc func() interface{}
type DispatchFunc func(interface{}) interface{}

// func( getstateFunc, dispatchFunc ) (nextfunc)
type MiddlewareFunc func(GetStateFunc, DispatchFunc) DispatchFunc

func ApplyMiddleware(getstate GetStateFunc, disp DispatchFunc, fns ...MiddlewareFunc) (createStore func() Store) {
	//reverse
	sort.Slice(fns, func(i, j int) bool {
		return true
	})
	wrappedStore := &storeWrap{
		getState:    getstate,
		dispatch:    disp,
		middlewares: fns,
	}

	return func() Store {
		return wrappedStore
	}
}

type storeWrap struct {
	getState    GetStateFunc
	dispatch    DispatchFunc
	middlewares []MiddlewareFunc
}

func (this storeWrap) Dispatch(action interface{}) interface{} {
	var next func(interface{}) interface{}
	for i, fn := range this.middlewares {
		if 0 == i {
			next = fn(this.getState, this.dispatch)
		} else {
			next = fn(this.getState, next)
		}
	}
	return next(action)
}

func (this storeWrap) GetState() interface{} {
	return this.getState()
}
