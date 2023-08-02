package cell

import (
	"encoding/binary"
	"gopher-dish/object"
	"io"
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

	for i, bagage := range c.Bagage {
		if bagage == nil {
			continue
		}
		cdesc.Bagage[i] = bagage.GetID()
	}

	return binary.Write(writer, binary.LittleEndian, cdesc)
}
