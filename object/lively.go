package object

type Lively interface {
	Object

	GetAge() uint32
	GetHealth() byte
	GetEnergy() byte
	GetGenomeHash() uint64

	IsDied() bool

	LoseHealth(health byte) bool
	SpendEnergy(energy byte) bool

	HealHealth(health byte) bool
	IncreaseEnergy(energy byte) bool

	Reproduce(dir Rotation) bool
	Die() bool
}
