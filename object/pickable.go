package object

type Pickable interface {
	Object
	GetWeight() byte
	PickUp() bool
	Drop(Position) bool
}
