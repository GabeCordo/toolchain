package patterns

type SafeInterface interface {
	Safe() SafeInterface
}
