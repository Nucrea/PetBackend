package board

type EventType int

const (
	EventTypeAddPoint EventType = iota
	EventTypeRemovePoint
	EventTypeMovePoint
)

type Event interface {
	GetPrimitiveId() string
	GetType() EventType
}

type event struct {
	PrimitiveId string
	Type        EventType
}

func (e event) GetPrimitiveId() string {
	return e.PrimitiveId
}

func (e event) GetType() EventType {
	return e.Type
}

type PointEventAdd struct {
	event
}

type PointEventRemove struct {
	event
}

type PointEventMove struct {
	event
	X, Y int
}
