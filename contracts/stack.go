package contracts

type Stack[T any] interface {
	Push(T)
	Pop() (T, bool)
	Peek() (T, bool)
	Len() int
	ExecuteAll(fn func(T) error) error
}
