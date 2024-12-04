package board

type PrimitiveType int

const (
	PrimitiveTypePoint = iota
	PrimitiveTypeLine
	PrimitiveTypeCircle
)

type Primitive interface {
	GetId() string
	GetType() PrimitiveType
}

type primitive struct {
	Id   string
	Type PrimitiveType
}

func (p primitive) GetId() string {
	return p.Id
}

func (p primitive) GetType() PrimitiveType {
	return p.Type
}

type Point struct {
	primitive
	X, Y int
}

type Line struct {
	primitive
	Points [2]Point
}

type Circle struct {
	primitive
	Center Point
	Radius int
}
