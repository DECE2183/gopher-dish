package cell

import "math"

// Lively interface implementation

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
