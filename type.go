package redux

type HasType interface {
	GetType() interface{}
}

type Type int

func (this Type) GetType() interface{} {
	return this
}
