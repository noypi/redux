package redux

import (
	"bytes"
	"fmt"
	"reflect"

	"github.com/noypi/util/reflect"
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
	return util.FlattenToType(this, refType)
}

func (this ReducerResult) CanFlattenTo(refType interface{}) (bRet bool) {
	return util.CanFlattenTo(this, refType)
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
