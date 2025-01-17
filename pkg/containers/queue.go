package containers

func NewQueue[T interface{}]() Queue[T] {
	return Queue[T]{
		front: NewStack[T](),
		back:  NewStack[T](),
	}
}

type Queue[T interface{}] struct {
	front, back *Stack[T]
}

func (m *Queue[T]) Push(val T) {
	m.front.Push(val)
}

func (m *Queue[T]) swapStacks() {
	if m.back.Empty() {
		for !m.front.Empty() {
			val, _ := m.front.Pop()
			m.back.Push(val)
		}
	}
}

func (m *Queue[T]) Pop() T {
	m.swapStacks()

	val, _ := m.back.Pop()
	return val
}

func (m *Queue[T]) Peek() T {
	m.swapStacks()

	val, _ := m.back.Peek()
	return val
}

func (m *Queue[T]) Empty() bool {
	return m.front.Empty() && m.back.Empty()
}
