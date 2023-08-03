package cell

import (
	"gopher-dish/object"
	"gopher-dish/world"
)

const (
	GenomeLength   = 256
	MemorySize     = 64
	StackDepth     = 32
	RegistersCount = 4
	SensorsCount   = 4
	BagageSize     = 4
)

const (
	BaseHealth              = 50
	BaseEnergy              = 20
	BaseWeight              = 5
	BaseEnergyDecrement     = 2
	BaseHealthDecrement     = 4
	BaseBiteStrength        = 40
	AgeInfluenceMultiplier  = 0.2
	BaseReproduceEnergyCost = 32
)

// Registers list
const (
	R0 = iota
	R1
	R2
	R3
)

type TriggerSource byte

// Triggers list
const (
	TRG_NONE = iota
)

type Sensor struct {
	JumpPosition  uint64
	TriggerSource TriggerSource
	Triggered     bool
}

type StackState struct {
	JumpPosition    uint64
	JumpRegisters   [RegistersCount]byte
	JumpCompareFlag byte
}

type Brain struct {
	CompareFlag byte
	Registers   [RegistersCount]byte
	Memory      [MemorySize]byte
	Stack       [StackDepth]StackState

	Sensors [SensorsCount]Sensor

	StackCounter   uint64
	CommandCounter uint64
}

type Cell struct {
	Name         uint64
	Generation   uint64
	ParentsChain object.ParentsChain

	Age    uint32
	Health byte
	Energy byte
	Weight byte

	Died   bool
	Killed bool
	Picked bool

	Genome     Genome
	Brain      Brain
	TriggerMap map[TriggerSource]*Sensor

	Bagage [BagageSize]object.Pickable

	World    *world.World
	Position object.Position
	Rotation object.Rotation
}

type saveCellDescriptor struct {
	Id           uint64
	Generation   uint64
	ParentsChain object.ParentsChain

	Age    uint32
	Health byte
	Energy byte
	Weight byte

	Died   bool
	Picked bool

	Genome Genome
	Brain  Brain

	Bagage [BagageSize]uint64

	Position object.Position
	Rotation object.Rotation
}

func New(w *world.World, parent *Cell, pos object.Position) *Cell {
	c := &Cell{Health: BaseHealth, Energy: BaseEnergy, Weight: BaseWeight, World: w}

	if parent != nil {
		for i := 0; i < object.RelatedDepth-1; i++ {
			c.ParentsChain[i+1] = parent.ParentsChain[i]
		}
		c.ParentsChain[0] = parent.Name
		c.Generation = parent.Generation + 1
		c.Genome = parent.Genome.Mutate()
		if c.Energy > parent.Energy {
			c.Energy = parent.Energy
		}
	} else {
		c.Genome = CreateBaseGenome()
	}

	c.Position = pos
	c.Name = w.ReserveID()

	if c.Name > 0 && w.PlaceObject(c, c.Position) {
		return c
	} else {
		return nil
	}
}

func (c *Cell) incCounter() uint64 {
	c.Brain.CommandCounter++
	if c.Brain.CommandCounter >= GenomeLength {
		c.Brain.CommandCounter = 0
	}
	return c.Brain.CommandCounter
}

func (c *Cell) recycle(rType uint64) {
	switch rType {
	case RCL_SUNENERGY:
		c.IncreaseEnergy(c.World.GetSunlightAtPosition(c.Position))
	case RCL_BAGAGE:
	}
}

func (c *Cell) getRelPos(rot object.Rotation) object.Position {
	newPos := c.Position
	switch c.Rotation.Rotate(rot.Degree).Degree {
	case 0:
		newPos.Y--
	case 45:
		newPos.X++
		newPos.Y--
	case 90:
		newPos.X++
	case 135:
		newPos.X++
		newPos.Y++
	case 180:
		newPos.Y++
	case 225:
		newPos.X--
		newPos.Y++
	case 270:
		newPos.X--
	case 315:
		newPos.X--
		newPos.Y--
	}
	newPos.X = (newPos.X + int32(c.World.Width)) % int32(c.World.Width)
	return newPos
}
