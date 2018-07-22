package redux

import (
	"fmt"
	"reflect"
)

type ReducerMap map[string]interface{}

func CombineReducers(fs ...interface{}) (reducer func(state, action interface{}) (out interface{})) {
	fs2 := []interface{}{}
	for _, o := range fs {
		t := reflect.TypeOf(o)

		switch t.Kind() {
		case reflect.Map:
			assertValidMap(t)
			fs2 = append(fs2, combineMap(castToMap(o)))
		case reflect.Func:
			assertValidReducer(t)
			fs2 = append(fs2, o)
		}
	}

	return combineReducerArr(fs2...)

}

type reducerInfo struct {
	T reflect.Type
	V reflect.Value

	Tstate  reflect.Type
	Taction reflect.Type
}

var g_tinterface = reflect.TypeOf(func(a interface{}) {}).In(0)
var g_tmap = reflect.TypeOf(map[string]interface{}{})

func combineReducerArr(fs ...interface{}) func(state, action interface{}) (out interface{}) {
	vfs := make([]reflect.Value, len(fs))
	for i, fn := range fs {
		vfs[i] = reflect.ValueOf(fn)
	}
	return func(state, action interface{}) (out interface{}) {
		var res []reflect.Value
		vaction := reflect.ValueOf(action)
		for i, vfn := range vfs {
			if 0 == i {
				res = vfn.Call([]reflect.Value{reflect.ValueOf(state), vaction})
			} else {
				res = vfn.Call([]reflect.Value{res[0], vaction})
			}
		}

		if 0 < len(res) {
			out = res[0].Interface()
		} else {
			out = state
		}
		return
	}
}

func combineMap(m map[string]interface{}) (reducer func(state, action interface{}) (out interface{})) {
	byType := map[reflect.Type]map[string]*reducerInfo{}
	for k, r := range m {
		info := &reducerInfo{T: reflect.TypeOf(r)}
		assertValidReducer(info.T)

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
			if byString, has = byType[g_tinterface]; !has {
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
			vfield := getVField(state, fieldname)
			tfield := vfield.Type()
			DBGf("state:%v, fieldname:%v, state type:%v", state, fieldname, tfield)

			if (info.Tstate.Kind() == reflect.Interface) || (tfield == info.Tstate) {
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

func castToMap(m interface{}) map[string]interface{} {
	v := reflect.ValueOf(m)
	return v.Convert(g_tmap).Interface().(map[string]interface{})
}

// because using
//        _, ok := v.(map[string]interface[})
// is not enough
func assertValidMap(t reflect.Type) {
	bValid := (t.Kind() == reflect.Map) &&
		(t.Key().Kind() == reflect.String) &&
		(t.Elem().Kind() == reflect.Interface)
	if !bValid {
		DBG("t.Key()=", t.Key())
		DBG("t.Elem()=", t.Elem())
		msgfmt := "Invalid MapReducer: '%v', should be in the form: 'map[string]interface{}'." +
			"t.Key(): %v, t.Elem(): %v"
		panic(fmt.Sprintf(msgfmt, t, t.Key(), t.Elem()))
	}
}

func assertValidReducer(tfn reflect.Type) {
	bValid := (tfn.Kind() == reflect.Func) &&
		(2 == tfn.NumIn()) &&
		(1 == tfn.NumOut())
	if !bValid {
		msgfmt := "Invalid reducer: '%v', should be in the form: 'func(state, action) (out)'." +
			"  A redux reducer is a func with 2 arguments."
		panic(fmt.Sprintf(msgfmt, tfn))
	}
}
