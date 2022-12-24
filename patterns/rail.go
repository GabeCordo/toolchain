package patterns

type RailInterface interface {
	Value(func(any))
	Unknown(func(error))
}
