package patterns

type SerializableInterface interface {
	ToJson(path Path) string
}
