package board

import "github.com/google/uuid"

type Board struct {
	Id         string
	Primitives map[string]Primitive
}

func (b *Board) ProcessEvent(event Event) {
	switch event.GetType() {
	case EventTypeAddPoint:
		p := Point{
			primitive: primitive{
				Id:   uuid.New().String(),
				Type: PrimitiveTypePoint,
			},
			X: 0,
			Y: 0,
		}
		b.Primitives[p.Id] = p
	case EventTypeRemovePoint:
		delete(b.Primitives, event.GetPrimitiveId())
	case EventTypeMovePoint:
		moveEvent := (event).(PointEventMove)
		primitiveId := moveEvent.GetPrimitiveId()

		p := b.Primitives[primitiveId].(Point)
		p.X = moveEvent.X
		p.Y = moveEvent.Y
		b.Primitives[primitiveId] = p
	}
}
