package cell

import (
	"gopher-dish/object"
	"math"
)

// Movable interface implementation

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
	c.SpendEnergy(c.Weight)
	pos.X = (pos.X + int32(c.World.Width)) % int32(c.World.Width)
	if !c.World.MoveObject(c, pos) {
		return false
	}
	c.Position = pos
	return true
}

func (c *Cell) MoveInDirection(rot object.Rotation) bool {
	return c.MoveToPosition(c.getRelPos(rot))
}

func (c *Cell) Rotate(rot object.Rotation) bool {
	c.Rotation.Degree = int32(math.Round(float64(c.Rotation.Rotate(rot.Degree).Degree)/45.0)) * 45
	c.SpendEnergy(c.Weight / 4)
	return true
}
