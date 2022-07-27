package object

type Lively interface {
	Object

	GetAge() uint32
	GetHealth() byte
	GetEnergy() byte

	IsDied() bool

	LoseHealth(health byte) bool
	SpendEnergy(energy byte) bool

	HealHealth(health byte) bool
	IncreaseEnergy(energy byte) bool

	Reproduce() bool
	Die() bool
}
