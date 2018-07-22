package redux

import (
	"reflect"
)

func getVField(a interface{}, fieldname string) (b reflect.Value) {
	v := reflect.ValueOf(a)
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() == reflect.Struct {
		fv := v.FieldByName(fieldname)
		DBG("getVField a=", a, ", fieldname=", fieldname)
		if fv.IsValid() {
			b = fv
		} else {
			b = reflect.ValueOf(ReducerResult{})
		}
	} else if v.Kind() == reflect.Map {
		fv := v.MapIndex(reflect.ValueOf(fieldname))
		bIsNil := (fv.Kind() == reflect.Ptr) && fv.IsNil()
		if !bIsNil && fv.IsValid() {
			b = fv
		} else {
			b = reflect.ValueOf(ReducerResult{})
		}
	}
	DBG("getVField out=", b)
	return
}

/*
func getFieldValue(a interface{}, fieldname string) (b interface{}) {
	v := reflect.ValueOf(a)
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() == reflect.Struct {
		fv := v.FieldByName(fieldname)
		DBG("getFieldValue a=", a, ", fieldname=", fieldname)
		if fv.IsValid() {
			b = fv.Interface()
		} else {
			b = ReducerResult{}
		}
	} else if v.Kind() == reflect.Map {
		fv := v.MapIndex(reflect.ValueOf(fieldname))
		bIsNil := (fv.Kind() == reflect.Ptr) && fv.IsNil()
		if !bIsNil && fv.IsValid() {
			b = fv.Interface()
		} else {
			b = ReducerResult{}
		}
	}
	DBG("getFieldValue out=", b)
	return
}
*/
func isValidField(t reflect.Type, v reflect.Value) (bRet bool) {
	if (t.Kind() == reflect.Ptr) && v.IsNil() {
		return
	}
	if (t.Kind() == reflect.String) && (0 == v.Len()) {
		return
	}

	return true
}
