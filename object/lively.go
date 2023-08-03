package object

const (
	RelatedDepth = 5
)

type ParentsChain [RelatedDepth]uint64

type Lively interface {
	Object

	GetAge() uint32
	GetHealth() byte
	GetEnergy() byte
	GetGenomeHash() uint64
	GetParentsChain() ParentsChain

	IsDied() bool
	IsKilled() bool
	IsReleated(another Lively) bool

	LoseHealth(health byte) bool
	SpendEnergy(energy byte) bool

	HealHealth(health byte) bool
	IncreaseEnergy(energy byte) bool

	Reproduce(dir Rotation) bool
	Bite(strength byte) byte
	Die() bool
}
