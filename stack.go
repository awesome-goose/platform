package platform

type stack[T any] []any

func NewStack[T any]() *stack[T] {
	return &stack[T]{}
}

func (s *stack[T]) Push(item T) {
	*s = append(*s, item)
}

func (s *stack[T]) Pop() (T, bool) {
	if len(*s) == 0 {
		var zero T
		return zero, false
	}
	lastIndex := len(*s) - 1
	item := (*s)[lastIndex]
	*s = (*s)[:lastIndex]
	return item.(T), true
}

func (s stack[T]) Peek() (T, bool) {
	if len(s) == 0 {
		var zero T
		return zero, false
	}
	return s[len(s)-1].(T), true
}

func (s stack[T]) Len() int {
	return len(s)
}

func (s stack[T]) ExecuteAll(fn func(T) error) error {
	for i := len(s) - 1; i >= 0; i-- {
		if err := fn(s[i].(T)); err != nil {
			return err
		}
	}

	return nil
}
