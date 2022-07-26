package object

type Position struct {
	X, Y int32
}

type Rotation struct {
	Degree int32
}

func (r Rotation) Rotate(degree int32) Rotation {
	r.Degree += degree % 360
	r.Degree = (360 + r.Degree) % 360
	return r
}

type Movable interface {
	Object
	GetPosition() Position
	GetRotation() Rotation
	MoveForward() bool
	MoveBackward() bool
	MoveLeft() bool
	MoveRight() bool
	MoveToPosition(Position) bool
	MoveInDirection(Rotation) bool
	Rotate(Rotation) bool
}
