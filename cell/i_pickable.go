package cell

import "gopher-dish/object"

// Pickable interface implementation

func (c *Cell) GetWeight() byte {
	return c.Weight
}

func (c *Cell) PickUp() bool {
	if c.Health == 0 && !c.Picked {
		c.World.RemoveObject(c.Name)
		c.Picked = true
		return true
	}

	return false
}

func (c *Cell) Drop(pos object.Position) bool {
	c.Picked = false
	pos.X = (pos.X + int32(c.World.Width)) % int32(c.World.Width)
	if !c.World.PlaceObject(c, pos) {
		return false
	}
	c.Position = pos
	return true
}
