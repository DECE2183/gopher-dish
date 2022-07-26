package cell

import (
	"gopher-dish/object"
	"gopher-dish/world"
	"math"
)

const (
	GenomeLength   = 256
	MemorySize     = 64
	StackDepth     = 32
	RegistersCount = 4
	SensorsCount   = 4
	RelatedDepth   = 6
)

const (
	BaseHealth = 50
	BaseEnergy = 25
	BaseWeight = 5
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
	Picked bool

	Genome     Genome
	Brain      Brain
	TriggerMap map[TriggerSource]*Sensor

	Bagage []object.Pickable

	World    *world.World
	Position object.Position
	Rotation object.Rotation
}

func New(world *world.World, parent *Cell) *Cell {
	c := &Cell{Health: BaseHealth, Energy: BaseEnergy, Weight: BaseWeight, World: world}

	c.Name = world.AddObject(c)
	if parent != nil {
		c.Generation = parent.Generation + 1
		for i := 0; i < RelatedDepth-1; i++ {
			c.ParentsChain[i+1] = parent.ParentsChain[i]
		}
		c.ParentsChain[0] = parent.Name
		c.Genome = parent.Genome.Mutate()
	} else {
		c.Genome = CreateBaseGenome()
	}

	return c
}

// Object interface
func (c Cell) GetID() uint64 {
	return c.Name
}
func (c Cell) GetInstance() object.Object {
	return c.World.GetObject(c.Name)
}
func (c Cell) Handle() {
	self := c.GetInstance().(*Cell)
	self.handleCommand(self.currentCommad())
}

// Pickable interface
func (c Cell) GetWeight() byte {
	return c.Weight
}
func (c Cell) PickUp() bool {
	if c.Health == 0 && !c.Picked {
		c.Picked = true
		return true
	}
	return false
}
func (c Cell) Drop() {
	c.Picked = false
}

// Movable interface
func (c Cell) GetPosition() object.Position {
	return c.Position
}
func (c Cell) GetRotation() object.Rotation {
	return c.Rotation
}
func (c Cell) MoveForward() bool {
	return c.MoveInDirection(c.Rotation)
}
func (c Cell) MoveBackward() bool {
	return c.MoveInDirection(c.Rotation.Rotate(180))
}
func (c Cell) MoveLeft() bool {
	return c.MoveInDirection(c.Rotation.Rotate(90))
}
func (c Cell) MoveRight() bool {
	return c.MoveInDirection(c.Rotation.Rotate(270))
}
func (c Cell) MoveToPosition(pos object.Position) bool {
	self, success := c.GetInstance().(*Cell)
	if !success {
		return false
	}

	self.Position = pos
	return true
}
func (c Cell) MoveInDirection(rot object.Rotation) bool {
	self, success := c.GetInstance().(*Cell)
	if !success {
		return false
	}

	switch self.Rotation.Degree {
	case 0:
		self.Position.X++
	case 45:
		self.Position.X++
		self.Position.Y--
	case 90:
		self.Position.Y--
	case 135:
		self.Position.X--
		self.Position.Y--
	case 180:
		self.Position.X--
	case 225:
		self.Position.X--
		self.Position.Y++
	case 270:
		self.Position.Y++
	case 315:
		self.Position.X++
		self.Position.Y++
	}

	return true
}
func (c Cell) Rotate(rot object.Rotation) bool {
	self, success := c.GetInstance().(*Cell)
	if !success {
		return false
	}

	self.Rotation.Degree = int32(math.Round(float64(self.Rotation.Rotate(rot.Degree).Degree)/45.0)) * 45
	return true
}

// Lively interface
func (c Cell) GetAge() uint32 {
	return c.Age
}
func (c Cell) GetHealth() byte {
	return c.Health
}
func (c Cell) GetEnergy() byte {
	return c.Energy
}
func (c Cell) Reproduce() bool {
	return false
}

// Private methods
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
	case RCL_BAGAGE:
	}
}
