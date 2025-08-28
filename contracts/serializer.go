package contracts

type Serializer interface {
	Serialize(data any) ([]byte, error)
}
