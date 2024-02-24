package object

const (
	RelatedDepth = 5
)

type ParentsChain [RelatedDepth]uint64

type Lively interface {
	Object

	GetAge() uint32
	GetGenomeHash() uint64
	GetParentsChain() ParentsChain

	GetHealth() byte
	LoseHealth(health byte) bool
	HealHealth(health byte) bool

	IsDied() bool
	IsKilled() bool
	IsReleated(another Lively) bool

	Reproduce(dir Rotation) bool
	Bite(strength byte) byte
	Die() bool
}
