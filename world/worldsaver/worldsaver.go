package worldsaver

import (
	"gopher-dish/cell"
	"gopher-dish/object"
)

type wDescriptor struct {
	Width, Height uint32

	Ticks uint64
	Year  uint64
	Epoch uint64

	ObjectCount   uint64
	ObjectIdCount uint64
}

type wCellDescriptor struct {
	Id           uint64
	Generation   uint64
	ParentsChain object.ParentsChain

	Age    uint32
	Health byte
	Energy byte
	Weight byte

	Died   bool
	Picked bool

	Genome cell.Genome
	Brain  cell.Brain

	Bagage [cell.BagageSize]uint64

	Position object.Position
	Rotation object.Rotation
}
