package reack

type DispatchFunc func(action interface{}) interface{}
type ConnectFunc func(component interface{}) interface{}
type StateToPropsFunc func(state, ownProps interface{}) interface{}
type DispatchToPropsFunc func(disp DispatchFunc, ownProps interface{}) interface{}

type CanCreate interface {
	Create() interface{}
}

func Connect(mapStateToProps StateToPropsFunc, mapDispatchToProps DispatchToPropsFunc) ConnectFunc {
	return func(component interface{}) (obj interface{}) {

		//props := mapStateToProps(state)
		//return FlattenToType(props, component)
		return nil
	}
}
