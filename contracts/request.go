package contracts

type Request interface {
	Read() ([]byte, error)
}

type Param interface {
	Get(key string) string
	GetAll(key string) []string
	Has(key string) bool
	All() map[string][]string
}

type Query interface {
	Get(key string) string
	GetAll(key string) []string
	Has(key string) bool
	All() map[string][]string
}

type Header interface {
	Get(key string) string
	GetAll(key string) []string
	Has(key string) bool
	All() map[string][]string
}

type Body[T any] interface {
	Get() T
}
