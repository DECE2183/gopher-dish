package object

type Lively interface {
	Object
	GetAge() uint32
	GetHealth() byte
	GetEnergy() byte
	Reproduce() bool
}
