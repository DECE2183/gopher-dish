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
	RelatedDepth   = 6
	BagageSize     = 4
)

const (
	BaseHealth              = 60
	BaseEnergy              = 28
	BaseWeight              = 5
	BaseEnergyDecrement     = 3
	BaseHealthDecrement     = 5
	AgeInfluenceMultiplier  = 2.8
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
	ParentsChain [RelatedDepth]uint64

	Age    uint32
	Health byte
	Energy byte
	Weight byte

	Died   bool
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
	ParentsChain [RelatedDepth]uint64

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

func New(w *world.World, parent *Cell) *Cell {
	c := &Cell{Health: BaseHealth, Energy: BaseEnergy, Weight: BaseWeight, World: w}

	if parent != nil {
		for i := 0; i < RelatedDepth-1; i++ {
			c.ParentsChain[i+1] = parent.ParentsChain[i]
		}
		c.ParentsChain[0] = parent.Name
		c.Generation = parent.Generation + 1
		c.Position = parent.Position
		c.Genome = parent.Genome.Mutate()
		if c.Energy > parent.Energy {
			c.Energy = parent.Energy
		}
	} else {
		c.Position = w.GetCenter()
		c.Position.Y /= 6
		c.Genome = CreateBaseGenome()
	}

	c.Name, c.Position = w.AddObject(c)

	if c.Name > 0 {
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
	c.LoseHealth(BaseEnergyDecrement / 2)
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
	switch c.Rotation.Degree {
	case 0:
		newPos.X++
	case 45:
		newPos.X++
		newPos.Y--
	case 90:
		newPos.Y--
	case 135:
		newPos.X--
		newPos.Y--
	case 180:
		newPos.X--
	case 225:
		newPos.X--
		newPos.Y++
	case 270:
		newPos.Y++
	case 315:
		newPos.X++
		newPos.Y++
	}
	newPos.X = newPos.X % int32(c.World.Width)
	return newPos
}
