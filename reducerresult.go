package redux

import (
	"bytes"
	"fmt"
	"reflect"
)

type ReducerResult map[string]interface{}

func (this ReducerResult) init(a interface{}) {
	v := reflect.ValueOf(a)
	t := reflect.TypeOf(a)

	if t.Kind() == reflect.Struct {
		for i := 0; i < v.NumField(); i++ {
			fv := v.Field(i)
			ft := t.Field(i)
			if !isValidField(ft.Type, fv) {
				continue
			}
			this[ft.Name] = fv.Interface()
		}
	}
}

func (this ReducerResult) Merge(b ReducerResult) (out ReducerResult) {
	DBG("reducer merge in this=", this, "; b=", b)
	out = ReducerResult{}
	for k, v := range this {
		if v2, has := b[k]; has {
			v = v2
		}
		out[k] = v
	}
	for k, v := range b {
		if _, has := out[k]; has {
			continue
		}
		out[k] = v
	}

	DBG("reducer merge out=", out)
	return
}

func (this ReducerResult) AddMergeField(b ReducerResult, fieldname string) {
	if !this.Has(fieldname) {
		this[fieldname] = b.Get(fieldname)
	} else {
		if (nil != this.Get(fieldname)) && (nil != b.Get(fieldname)) {
			v3 := Merge(this.Get(fieldname), b.Get(fieldname))
			if res3, ok := v3.(ReducerResult); ok {
				this = this.Merge(res3)
			} else {
				this[fieldname] = v3
			}
		}
	}

	return
}

func (this ReducerResult) ToType(refType interface{}) interface{} {
	t, ok := refType.(reflect.Type)
	if !ok {
		t = reflect.TypeOf(refType)
	}

	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	out := reflect.New(t).Elem()

	for k, v := range this {
		if nil == v {
			continue
		}
		fvout := out.FieldByName(k)
		fv := reflect.ValueOf(v)

		if fv.Type().AssignableTo(fvout.Type()) {
			fvout.Set(fv)
		} else {
			copyProps(fvout, fv)
		}
	}

	return out.Interface()
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

func (this ReducerResult) CanFlattenTo(refType interface{}) (bRet bool) {
	t0, ok := refType.(reflect.Type)
	if !ok {
		t0 = reflect.TypeOf(refType)
	}

	t := t0
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return
	}

	if t.NumField() < len(this) {
		DBG("reftype numfield is lesser than ReducerResult, reftype=", t.Name(), "; reducerresult=", this)
		return
	}

	for k, v := range this {
		if nil == v {
			continue
		}
		ft, has := t0.FieldByName(k)
		if !has {
			return
		}
		vt := reflect.TypeOf(v)
		if !ft.Type.AssignableTo(vt) && !canAssignFields(ft.Type, vt) {
			DBG("canflatten was not assignable, ft.name=", ft.Name, "; ft typename=", ft.Type.Name(), "; v=", v)
			DBG("ft.pkg=", ft.PkgPath)
			return
		}
	}

	return true
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

func (this ReducerResult) Has(k string) bool {
	_, ok := this[k]
	return ok
}

func (this ReducerResult) Get(k string) interface{} {
	v, _ := this[k]
	return v
}

func (this ReducerResult) String() string {
	var buf bytes.Buffer
	buf.WriteString("{")
	for k, v := range this {
		buf.WriteString(fmt.Sprintf("[%v: %v] ", k, v))
	}
	buf.WriteString("}")
	return buf.String()
}
