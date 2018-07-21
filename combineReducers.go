package redux

import (
	"fmt"
	"reflect"
)

type Reducer func(state, action interface{}) (newState interface{})

func (fn Reducer) Combine(fn2 Reducer) Reducer {
	return func(state, action interface{}) interface{} {
		return fn2(fn(state, action), action)
	}
}

func DefaultReducer(state, action interface{}) (out interface{}) {
	return state
}

func CombineReducersArr(a, a1 Reducer, as ...Reducer) Reducer {
	r0 := a.Combine(a1)
	for _, r := range as {
		r0 = r0.Combine(r)
	}

	return r0
}

func combineReducersArr(as ...Reducer) Reducer {
	if 1 == len(as) {
		return as[0]
	} else if 2 == len(as) {
		return CombineReducersArr(as[0], as[1])
	} else if 2 < len(as) {
		return CombineReducersArr(as[0], as[1], as[2:]...)
	}

	return DefaultReducer
}

type FieldReducer struct {
	Name    string
	Reducer Reducer
}

type ReducersList []Reducer

type ReducerMap map[string]interface{}

func (this ReducersList) Append(fr FieldReducer) (ls ReducersList) {
	return append(this, fr.Reducer)
}

/*
func (this ReducerMap) Add(frs []FieldReducer) {
	for _, fr := range frs {
		ls, _ := this[fr.Name]
		ls = ls.Append(fr)
		this[fr.Name] = ls
	}
}
*/
type reducerInfo struct {
	T reflect.Type
	V reflect.Value

	Tstate  reflect.Type
	Taction reflect.Type
}

var tinterface = reflect.TypeOf(func(a interface{}) {}).In(0)

func CombineReducers(m ReducerMap) (reducer func(state, action interface{}) (out interface{})) {
	byType := map[reflect.Type]map[string]*reducerInfo{}
	for k, r := range m {
		info := &reducerInfo{T: reflect.TypeOf(r)}
		bValid := (info.T.Kind() == reflect.Func) &&
			(2 == info.T.NumIn()) &&
			(1 == info.T.NumOut())
		if !bValid {
			msgfmt := "Invalid reducer: '%v', should be in the form: 'func(state, action) (out)'." +
				"  A redux reducer is a func with 2 arguments."
			panic(fmt.Sprintf(msgfmt, info))
		}
		info.Tstate = info.T.In(0)
		info.Taction = info.T.In(1)
		info.V = reflect.ValueOf(r)

		byString, has := byType[info.Taction]
		if !has {
			byString = map[string]*reducerInfo{}
			byType[info.Taction] = byString
		}
		byString[k] = info
	}

	return func(state, action interface{}) (out interface{}) {
		DBG(">>> combine state=", state, "; action=", action)
		taction := reflect.TypeOf(action)
		byString, has := byType[taction]
		if !has {
			DBGf("taction is not found, taction: %v", taction)
			if byString, has = byType[tinterface]; !has {
				DBG("still not found in tinterface")
				return state
			}
		}

		res := ReducerResult{}
		if stateres, ok := state.(ReducerResult); ok {
			res = res.Merge(stateres)
		} else {
			res.init(state)
		}

		vaction := reflect.ValueOf(action)

		for fieldname, info := range byString {
			DBGf("state:%v, fieldname:%v", state, fieldname)
			vfield := getVField(state, fieldname)

			if (info.Tstate.Kind() == reflect.Interface) || (vfield.Type() == info.Tstate) {
				DBGf("executing f=%v", info.T)
				vres := info.V.Call([]reflect.Value{vfield, vaction})
				state2 := vres[0].Interface()

				if res2, ok := state2.(ReducerResult); ok {
					res = res.Merge(res2)
				} else {
					res[fieldname] = state2
				}

			} else {
				DBGf("ignoring reducer: %v, because vfield is not same vfield: %v, state:%v, ",
					info.T, vfield.Type(), reflect.TypeOf(state))
			}
		}

		out = res
		if res.CanFlattenTo(state) {
			out = res.ToType(state)
		}

		DBG("<<< combine out=", out)
		return
	}
}

/*
func CombineReducers2(frs []FieldReducer) Reducer {
	m := ReducerMap{}
	m.Add(frs)

	m2 := map[string]Reducer{}
	for fieldname, reducers := range m {
		m2[fieldname] = combineReducersArr(reducers...)
	}

	return func(state, action interface{}) (out interface{}) {

		DBG(">>> combine state=", state, "; action=", action)
		res := ReducerResult{}
		if stateres, ok := state.(ReducerResult); ok {
			res = res.Merge(stateres)
		} else {
			res.init(state)
		}
		DBG("res=", res)

		for fieldname, reducer := range m2 {
			DBG("+fieldname=", fieldname)
			fvalue := getFieldValue(state, fieldname)
			DBG("fvalue=", fvalue)

			state2 := reducer(fvalue, action)
			if res2, ok := state2.(ReducerResult); ok {
				res = res.Merge(res2)
			} else {
				res[fieldname] = state2
			}
			DBG("-fieldname=", fieldname, " - done")
		}

		out = res
		if res.CanFlattenTo(state) {
			out = res.ToType(state)
		}

		DBG("<<< combine out=", out)

		return
	}
}
*/
