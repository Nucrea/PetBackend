package containers

func NewStack[T interface{}]() *Stack[T] {
	return &Stack[T]{[]T{}}
}

type Stack[T interface{}] struct {
	arr []T
}

func (s *Stack[T]) Empty() bool {
	return len(s.arr) <= 0
}

func (s *Stack[T]) Push(val T) {
	s.arr = append(s.arr, val)
}

func (s *Stack[T]) Peek() (T, bool) {
	if len(s.arr) <= 0 {
		var t T
		return t, false
	}
	return s.arr[len(s.arr)-1], true
}

func (s *Stack[T]) Pop() (T, bool) {
	if len(s.arr) <= 0 {
		var t T
		return t, false
	}

	maxIndex := len(s.arr) - 1
	element := s.arr[maxIndex]
	s.arr = s.arr[:maxIndex]
	return element, true
}
