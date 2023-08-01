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

// Object interface
func (c *Cell) GetID() uint64 {
	return c.Name
}
func (c *Cell) Prepare() {
	if c.Died {
		return
	} else if c.Health <= 0 {
		c.Die()
		return
	}

	c.executeCommand(c.currentCommad())
}
func (c *Cell) Handle(yearChanged, epochChanged bool) {
	if c.Died {
		return
	} else if c.Health <= 0 {
		c.Die()
		return
	}

	for i := 0; i < GenomeLength && !c.handleCommand(c.currentCommad()); i++ {
		c.SpendEnergy(BaseEnergyDecrement)
	}

	if yearChanged {
		c.Age++
	}
}

// Pickable interface
func (c *Cell) GetWeight() byte {
	return c.Weight
}
func (c *Cell) PickUp() bool {
	if c.Health == 0 && !c.Picked {
		c.Picked = true
		return true
	}

	return false
}
func (c *Cell) Drop() {
	c.Picked = false
}

// Movable interface
func (c *Cell) GetPosition() object.Position {
	return c.Position
}
func (c *Cell) GetRotation() object.Rotation {
	return c.Rotation
}
func (c *Cell) MoveForward() bool {
	return c.MoveInDirection(c.Rotation)
}
func (c *Cell) MoveBackward() bool {
	return c.MoveInDirection(c.Rotation.Rotate(180))
}
func (c *Cell) MoveLeft() bool {
	return c.MoveInDirection(c.Rotation.Rotate(90))
}
func (c *Cell) MoveRight() bool {
	return c.MoveInDirection(c.Rotation.Rotate(270))
}
func (c *Cell) MoveToPosition(pos object.Position) bool {
	c.Position = pos
	return true
}
func (c *Cell) MoveInDirection(rot object.Rotation) bool {
	c.SpendEnergy(c.Weight)
	return c.World.MoveObject(c, c.getRelPos(rot))
}
func (c *Cell) Rotate(rot object.Rotation) bool {
	c.Rotation.Degree = int32(math.Round(float64(c.Rotation.Rotate(rot.Degree).Degree)/45.0)) * 45
	return true
}

// Lively interface
func (c *Cell) GetAge() uint32 {
	return c.Age
}
func (c *Cell) GetHealth() byte {
	return c.Health
}
func (c *Cell) GetEnergy() byte {
	return c.Energy
}
func (c *Cell) IsDied() bool {
	return c.Died
}
func (c *Cell) LoseHealth(health byte) bool {
	if c.Died {
		return false
	}
	if health < c.Health {
		c.Health -= health
	} else {
		c.Die()
	}

	return true
}
func (c *Cell) SpendEnergy(energy byte) bool {
	if c.Energy > 0 {
		// Decrement energy
		energyDec := byte(math.Round(float64(energy) + float64(c.Age)*AgeInfluenceMultiplier))
		if energyDec < c.Energy {
			c.Energy -= energyDec
		} else {
			c.Energy = 0
		}
	} else {
		// If there is no energy then decrement health
		healthDec := byte(math.Round(BaseHealthDecrement * float64(c.Age) * AgeInfluenceMultiplier))
		c.LoseHealth(healthDec)
	}

	return true
}
func (c *Cell) HealHealth(health byte) bool {
	if c.Died {
		return false
	}

	c.Health += health
	return true
}
func (c *Cell) IncreaseEnergy(energy byte) bool {
	if c.Died {
		return false
	}

	c.Energy += energy
	return true
}
func (c *Cell) Reproduce() bool {
	if c.Energy <= BaseReproduceEnergyCost/2 {
		c.SpendEnergy(BaseReproduceEnergyCost)
		return false
	}
	newCell := New(c.World, c)
	if newCell == nil {
		return false
	}
	c.SpendEnergy(BaseReproduceEnergyCost)
	return true
}
func (c *Cell) Die() bool {
	c.Health = 0
	c.Died = true
	c.World.RemoveObject(c.Name)
	return true
}

// Private methods
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
