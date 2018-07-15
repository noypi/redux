package redux

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

type ReducerMap map[string]ReducersList

func (this ReducersList) Append(fr FieldReducer) (ls ReducersList) {
	return append(this, fr.Reducer)
}

func (this ReducerMap) Add(frs []FieldReducer) {
	for _, fr := range frs {
		ls, _ := this[fr.Name]
		ls = ls.Append(fr)
		this[fr.Name] = ls
	}
}

func CombineReducers(frs []FieldReducer) Reducer {
	m := ReducerMap{}
	m.Add(frs)

	return func(state, action interface{}) (out interface{}) {

		DBG(">>> combine state=", state, "; action=", action)
		res := ReducerResult{}
		if stateres, ok := state.(ReducerResult); ok {
			res = res.Merge(stateres)
		} else {
			res.init(state)
		}
		DBG("res=", res)

		for fieldname, reducers := range m {
			DBG("+fieldname=", fieldname)
			fvalue := getFieldValue(state, fieldname)
			DBG("fvalue=", fvalue)
			reducer := combineReducersArr(reducers...)

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
