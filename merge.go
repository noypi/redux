package redux

import (
	"reflect"
)

// merges a and b
// will ignore nil fields and empty strings
func merge(a, b interface{}) (c ReducerResult) {
	DBG("merge in a=", a, "; b=", b)
	at := reflect.TypeOf(a)
	av := reflect.ValueOf(a)

	bt := reflect.TypeOf(b)
	bv := reflect.ValueOf(b)

	c = ReducerResult{}
	if (at.Kind() == reflect.Struct) && (bt.Kind() == reflect.Struct) {
		for i := 0; i < at.NumField(); i++ {
			fat := at.Field(i)
			fav := av.Field(i)
			if !isValidField(fat.Type, fav) {
				continue
			}
			var v interface{} = fav.Interface()

			if _, has := bt.FieldByName(fat.Name); has {
				bav := bv.FieldByName(fat.Name)
				if fat.Type.Kind() == reflect.Struct {
					res := Merge(v, bav.Interface())
					c[fat.Name] = res
					if res2, ok := res.(ReducerResult); ok {
						if res2.CanFlattenTo(fat.Type) {
							c[fat.Name] = res2.ToType(fat.Type)
						}
					}
					continue
				}
				v = bav.Interface()
			}
			c[fat.Name] = v
		}
	}

	if bt.Kind() == reflect.Struct {
		for i := 0; i < bt.NumField(); i++ {
			bat := bt.Field(i)
			bav := bv.Field(i)
			if !isValidField(bat.Type, bav) {
				continue
			}
			if c.Has(bat.Name) {
				continue
			}
			c[bat.Name] = bav.Interface()
		}
	}

	return
}

func Merge(a, b interface{}) (out interface{}) {
	c := merge(a, b)

	if c.CanFlattenTo(a) {
		out = c.ToType(a)
	} else {
		out = c
	}
	DBG("Merge out=", out)

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
