package cell

// Pickable interface implementation

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
