package redux

type Store interface {
	GetState() interface{}
	Dispatch(action interface{}) interface{}
}
