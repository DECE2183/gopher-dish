package cell

import (
	"encoding/binary"
	"gopher-dish/object"
	"io"
	"math"
)

// Object interface implementation

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

	if c.Picked {
		c.SpendEnergy(BaseEnergyDecrement)
		return
	}

	for i := 0; i < GenomeLength && !c.handleCommand(c.currentCommad()); i++ {
		c.SpendEnergy(BaseEnergyDecrement)
	}

	if yearChanged {
		c.Age++
	}
}

func (c *Cell) Save(writer io.Writer) (err error) {
	binary.Write(writer, binary.LittleEndian, uint64(object.TYPE_CELL))

	cdesc := saveCellDescriptor{
		Id:           c.Name,
		Generation:   c.Generation,
		ParentsChain: c.ParentsChain,
		Age:          c.Age,
		Health:       c.Health,
		Energy:       c.Energy,
		Weight:       c.Weight,
		Died:         c.Died,
		Picked:       c.Picked,
		Genome:       c.Genome,
		Brain:        c.Brain,
		Position:     c.Position,
		Rotation:     c.Rotation,
	}

	cdesc.BagageSelected = c.BagageSelected
	cdesc.BagageFullness = c.BagageFullness
	for i, bagage := range c.Bagage {
		if bagage == nil {
			continue
		}
		cdesc.Bagage[i] = bagage.GetID()
	}

	return binary.Write(writer, binary.LittleEndian, cdesc)
}

func (c *Cell) GetEnergy() byte {
	return c.Energy
}

func (c *Cell) SpendEnergy(energy byte) bool {
	energyDec := uint32(math.Round(float64(energy) + float64(c.Age)*AgeInfluenceMultiplier))
	if energyDec < uint32(c.Energy) {
		// Decrement energy
		c.Energy -= byte(energyDec)
	} else {
		// If there is no energy then decrement health
		energyDec -= uint32(c.Energy)
		c.Energy = 0
		healthDec := uint32(math.Round(float64(energyDec) + BaseHealthDecrement + float64(c.Age)*AgeInfluenceMultiplier))
		if healthDec > 255 {
			healthDec = 255
		}
		c.LoseHealth(byte(healthDec))
	}

	return true
}

func (c *Cell) IncreaseEnergy(energy byte) bool {
	if c.Died {
		return false
	}

	if int(c.Energy)+int(energy) > 255 {
		c.Energy = 255
	} else {
		c.Energy += energy
	}

	return true
}
