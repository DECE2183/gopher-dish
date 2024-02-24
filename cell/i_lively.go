package cell

import (
	"gopher-dish/object"
	"math"
)

// Lively interface implementation

func (c *Cell) GetAge() uint32 {
	return c.Age
}

func (c *Cell) GetGenomeHash() uint64 {
	return c.Genome.Hash
}

func (c *Cell) GetParentsChain() object.ParentsChain {
	return c.ParentsChain
}

func (c *Cell) GetHealth() byte {
	return c.Health
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

func (c *Cell) HealHealth(health byte) bool {
	if c.Died {
		return false
	}

	c.Health += health
	return true
}

func (c *Cell) IsDied() bool {
	return c.Died
}

func (c *Cell) IsKilled() bool {
	return c.Killed
}

func (c *Cell) IsReleated(another object.Lively) bool {
	ochain := another.GetParentsChain()
	oid := another.GetID()
	for i := range ochain {
		if ochain[i] == c.Name || ochain[i] == c.ParentsChain[i] || c.ParentsChain[i] == oid {
			return true
		}
	}
	return false
}

func (c *Cell) Reproduce(rot object.Rotation) bool {
	if c.Energy <= BaseReproduceEnergyCost/2 {
		c.SpendEnergy(BaseReproduceEnergyCost / 2)
		return false
	}
	pos := c.getRelPos(rot)
	newCell := New(c.World, c, pos)
	if newCell == nil {
		c.SpendEnergy(BaseReproduceEnergyCost / 2)
		return false
	}
	c.SpendEnergy(BaseReproduceEnergyCost)
	return true
}

func (c *Cell) Bite(strength byte) byte {
	if c.Died {
		c.World.RemoveObject(c.Name)
		return c.Energy
	}

	biteStrength := int(math.Round(float64(strength) + float64(c.Age)*AgeInfluenceMultiplier - float64(c.Weight)))
	if biteStrength > 255 {
		biteStrength = 255
	} else if biteStrength <= 0 {
		return 0
	}

	c.SpendEnergy(byte(biteStrength))
	c.Killed = true

	var energy int
	if c.Died {
		energy = int(c.Energy)
		c.World.RemoveObject(c.Name)
	} else {
		energy = biteStrength - int(math.Round(float64(c.Age)*AgeInfluenceMultiplier)) + int(c.Weight)
	}

	if energy > 255 {
		energy = 255
	} else if energy <= 0 {
		return 0
	}

	return byte(energy)
}

func (c *Cell) Die() bool {
	c.Energy += BaseEnergyDecrement
	c.Health = 0
	c.Died = true
	return true
}
