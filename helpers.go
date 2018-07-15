package redux

import (
	"reflect"
)

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

func isValidField(t reflect.Type, v reflect.Value) (bRet bool) {
	if (t.Kind() == reflect.Ptr) && v.IsNil() {
		return
	}
	if (t.Kind() == reflect.String) && (0 == v.Len()) {
		return
	}

	return true
}

// copies assignable fields from b to a
func copyProps(a, b reflect.Value) {
	at, av := a.Type(), a
	bt, bv := b.Type(), b
	if (at.Kind() != reflect.Struct) || (bt.Kind() != reflect.Struct) {
		return
	}

	for i := 0; i < at.NumField(); i++ {
		fat := at.Field(i)
		bat, has := bt.FieldByName(fat.Name)
		if !has {
			continue
		}

		if bat.Type.AssignableTo(fat.Type) {
			av.Field(i).Set(bv.FieldByName(fat.Name))
		}
	}
}

func canAssignFields(a, b reflect.Type) (bRet bool) {
	if a.NumField() < b.NumField() {
		return
	}
	for i := 0; i < b.NumField(); i++ {
		fb := b.Field(i)
		fa, has := a.FieldByName(fb.Name)
		if !has {
			return
		}
		if !fb.Type.AssignableTo(fa.Type) {
			return
		}
	}

	return true
}
