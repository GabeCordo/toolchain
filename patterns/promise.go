package patterns

type function func(...any)

type Promise interface {
	OnComplete(function)
	Failed(function)
}
