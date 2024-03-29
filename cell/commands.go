package cell

import (
	"gopher-dish/object"
	"math"
	"math/rand"
)

func truncCmd(cmd Command, max uint64) uint64 {
	return uint64(cmd) % max
}

type Command byte
type CommandHandler func(*Cell)

type CommandDescriptor struct {
	handler     CommandHandler
	synchronous bool
}

// Commands list
const (
	// No operation command
	CMD_NOP = iota
	// Branch commands
	CMD_CMP  // + compare
	CMD_JMP  // + jump
	CMD_DIVE // + dive into function
	CMD_LIFT // + lift from function
	// Memory commands
	CMD_PUT  // + put const to reg
	CMD_RAND // + put random value to reg
	CMD_SAVE // + put reg value to mem
	CMD_LOAD // + load mem value to reg
	// Math commands
	CMD_ADD // +
	CMD_SUB // +
	CMD_MUL // +
	CMD_DIV // +
	// Application commands
	CMD_MOVE        // + move self cell
	CMD_ROTATE      // + rotate self cell
	CMD_CHECKPOS    // + get type of object at near position
	CMD_CHECKREL    // + get releations with cell at near position
	CMD_BITE        // + bite another cell
	CMD_SHAREENERGY // + share energy with near cell
	CMD_RECYCLE     // + recycle something to energy
	CMD_REPRODUCE   // + reproduce
	// Bagage commands
	CMD_PICKUP    // + pickup something
	CMD_DROP      // + drop selected item from bag
	CMD_BAGSIZE   // + count items in bag
	CMD_BAGACTIVE // + set active item in bag
	CMD_BAGENERGY // + get energy of selected item in bag
	CMD_BAGCHECK  // + get type of selected item in bag
	// Stats commands
	CMD_GETAGE     // + get self age
	CMD_GETHEALTH  // + get self health
	CMD_GETENERGY  // + get self energy
	CMD_GETCOUNTER // + get current command counter

	CMD_ENUM_SIZE
)

// Directions list
const (
	DIR_0 = iota
	DIR_45
	DIR_90
	DIR_135
	DIR_180
	DIR_225
	DIR_270
	DIR_315
)

// Compare conditions list
const (
	CND_NONE = iota

	CND_EQ    = 1 << iota
	CND_NEQ   = 1 << iota
	CND_LESS  = 1 << iota
	CND_GREAT = 1 << iota

	CND_SUCCESS = 1 << iota
	CND_FAIL    = 1 << iota

	CND_ENUM_SIZE = 256
)

// Object types list
const (
	OBJ_EMPTY = iota

	OBJ_WALL
	OBJ_BODY
	OBJ_DEAD

	OBJ_RELATED
	OBJ_UNRELATED
)

// Recycle targets list
const (
	RCL_NONE = iota
	RCL_SUNENERGY
	RCL_BAGAGE
	RCL_ENUM_SIZE
)

var commandMap = map[Command]CommandDescriptor{
	CMD_NOP: {func(c *Cell) {
		c.incCounter()
	}, false},

	// Compare values in registers and put result to CompareFlag
	CMD_CMP: {func(c *Cell) {
		op1 := c.Brain.Registers[truncCmd(c.Genome.Code[c.incCounter()], RegistersCount)]
		op2 := c.Brain.Registers[truncCmd(c.Genome.Code[c.incCounter()], RegistersCount)]

		c.Brain.CompareFlag = 0

		if op1 == op2 {
			c.Brain.CompareFlag |= CND_EQ
		} else if op1 != op2 {
			c.Brain.CompareFlag |= CND_NEQ
		}

		if op1 < op2 {
			c.Brain.CompareFlag |= CND_LESS
		} else if op1 > op2 {
			c.Brain.CompareFlag |= CND_GREAT
		}

		c.incCounter()
	}, false},
	// Jump according to condition in CompareFlag
	CMD_JMP: {func(c *Cell) {
		cond := truncCmd(c.Genome.Code[c.incCounter()], CND_ENUM_SIZE)
		pos := truncCmd(c.Genome.Code[c.incCounter()], GenomeLength)

		if cond == CND_NONE || cond&uint64(c.Brain.CompareFlag) > 0 {
			c.Brain.CommandCounter = pos
		} else {
			c.incCounter()
		}
	}, false},
	// Memorize current registers and command counter and dive into subprogramm
	CMD_DIVE: {func(c *Cell) {
		cond := truncCmd(c.Genome.Code[c.incCounter()], CND_ENUM_SIZE)
		pos := truncCmd(c.Genome.Code[c.incCounter()], GenomeLength)

		if c.Brain.StackCounter >= StackDepth {
			c.incCounter()
			return
		}

		cnt := c.Brain.StackCounter
		c.Brain.Stack[cnt].JumpPosition = c.Brain.CommandCounter + 1
		if c.Brain.Stack[cnt].JumpPosition >= GenomeLength {
			c.Brain.Stack[cnt].JumpPosition = 0
		}
		c.Brain.Stack[cnt].JumpRegisters = c.Brain.Registers
		c.Brain.Stack[cnt].JumpCompareFlag = c.Brain.CompareFlag
		c.Brain.StackCounter++

		if cond == CND_NONE || cond&uint64(c.Brain.CompareFlag) > 0 {
			c.Brain.CommandCounter = pos
		} else {
			c.incCounter()
		}
	}, false},
	// Return to main programm
	CMD_LIFT: {func(c *Cell) {
		cond := truncCmd(c.Genome.Code[c.incCounter()], CND_ENUM_SIZE)

		if c.Brain.StackCounter == 0 {
			c.incCounter()
			return
		}
		c.Brain.StackCounter--

		if cond == CND_NONE || cond&uint64(c.Brain.CompareFlag) > 0 {
			cnt := c.Brain.StackCounter
			c.Brain.CompareFlag = c.Brain.Stack[cnt].JumpCompareFlag
			c.Brain.Registers = c.Brain.Stack[cnt].JumpRegisters
			c.Brain.CommandCounter = c.Brain.Stack[cnt].JumpPosition
		} else {
			c.incCounter()
		}
	}, false},

	// Put value to register
	CMD_PUT: {func(c *Cell) {
		reg := truncCmd(c.Genome.Code[c.incCounter()], RegistersCount)
		val := byte(c.Genome.Code[c.incCounter()])

		c.Brain.Registers[reg] = val
		c.incCounter()
	}, false},
	// Put random value to register
	CMD_RAND: {func(c *Cell) {
		reg := truncCmd(c.Genome.Code[c.incCounter()], RegistersCount)
		val := byte(rand.Uint32() % 256)

		c.Brain.Registers[reg] = val
		c.incCounter()
	}, false},
	// Save value from register to memory
	CMD_SAVE: {func(c *Cell) {
		reg := truncCmd(c.Genome.Code[c.incCounter()], RegistersCount)
		mem := truncCmd(c.Genome.Code[c.incCounter()], MemorySize)

		c.Brain.Memory[mem] = c.Brain.Registers[reg]
		c.incCounter()
	}, false},
	// Load value from memory to registor
	CMD_LOAD: {func(c *Cell) {
		reg := truncCmd(c.Genome.Code[c.incCounter()], RegistersCount)
		mem := truncCmd(c.Genome.Code[c.incCounter()], MemorySize)

		c.Brain.Registers[reg] = c.Brain.Memory[mem]
		c.incCounter()
	}, false},

	// Add two regs command
	CMD_ADD: {func(c *Cell) {
		dest := truncCmd(c.Genome.Code[c.incCounter()], RegistersCount)
		src := truncCmd(c.Genome.Code[c.incCounter()], RegistersCount)
		op := truncCmd(c.Genome.Code[c.incCounter()], RegistersCount)

		c.Brain.Registers[dest] = c.Brain.Registers[src] + c.Brain.Registers[op]
		c.incCounter()
	}, false},
	// Subtract two regs command
	CMD_SUB: {func(c *Cell) {
		dest := truncCmd(c.Genome.Code[c.incCounter()], RegistersCount)
		src := truncCmd(c.Genome.Code[c.incCounter()], RegistersCount)
		op := truncCmd(c.Genome.Code[c.incCounter()], RegistersCount)

		c.Brain.Registers[dest] = c.Brain.Registers[src] - c.Brain.Registers[op]
		c.incCounter()
	}, false},
	// Multiply two regs command
	CMD_MUL: {func(c *Cell) {
		dest := truncCmd(c.Genome.Code[c.incCounter()], RegistersCount)
		src := truncCmd(c.Genome.Code[c.incCounter()], RegistersCount)
		op := truncCmd(c.Genome.Code[c.incCounter()], RegistersCount)

		c.Brain.Registers[dest] = c.Brain.Registers[src] * c.Brain.Registers[op]
		c.incCounter()
	}, false},
	// Divide two regs command
	CMD_DIV: {func(c *Cell) {
		dest := truncCmd(c.Genome.Code[c.incCounter()], RegistersCount)
		src := truncCmd(c.Genome.Code[c.incCounter()], RegistersCount)
		op := truncCmd(c.Genome.Code[c.incCounter()], RegistersCount)

		if c.Brain.Registers[op] == 0 {
			c.Brain.Registers[dest] = 0xFF
		} else {
			c.Brain.Registers[dest] = c.Brain.Registers[src] / c.Brain.Registers[op]
		}

		c.incCounter()
	}, false},

	// Relative move command
	CMD_MOVE: {func(c *Cell) {
		dirReg := truncCmd(c.Genome.Code[c.incCounter()], RegistersCount)
		dir := int32(c.Brain.Registers[dirReg]%8) * 45

		moveSuc := c.MoveInDirection(object.Rotation{Degree: dir})

		if moveSuc {
			c.Brain.CompareFlag = CND_SUCCESS
		} else {
			c.Brain.CompareFlag = CND_FAIL
		}
		c.incCounter()
	}, true},
	// Relative rotation  command
	CMD_ROTATE: {func(c *Cell) {
		dirReg := truncCmd(c.Genome.Code[c.incCounter()], RegistersCount)
		dir := int32(c.Brain.Registers[dirReg]%8) * 45

		rotSuc := c.Rotate(object.Rotation{Degree: dir})

		if rotSuc {
			c.Brain.CompareFlag = CND_SUCCESS
		} else {
			c.Brain.CompareFlag = CND_FAIL
		}
		c.incCounter()
	}, false},
	// Check what is located at near cell
	CMD_CHECKPOS: {func(c *Cell) {
		dest := truncCmd(c.Genome.Code[c.incCounter()], RegistersCount)
		dirReg := truncCmd(c.Genome.Code[c.incCounter()], RegistersCount)
		dir := int32(c.Brain.Registers[dirReg]%8) * 45

		pos := c.getRelPos(object.Rotation{Degree: dir})
		if pos.Y < 0 || pos.Y >= int32(c.World.Height) {
			c.Brain.Registers[dest] = OBJ_WALL
			c.incCounter()
			return
		}

		o := c.World.GetObjectAtPosition(pos)
		if o == nil {
			c.Brain.Registers[dest] = OBJ_EMPTY
		} else {
			switch v := o.(type) {
			case object.Lively:
				if v.IsDied() {
					c.Brain.Registers[dest] = OBJ_DEAD
				} else {
					c.Brain.Registers[dest] = OBJ_BODY
				}
			}
		}

		c.incCounter()
	}, true},
	// Check releations of near cell
	CMD_CHECKREL: {func(c *Cell) {
		dest := truncCmd(c.Genome.Code[c.incCounter()], RegistersCount)
		dirReg := truncCmd(c.Genome.Code[c.incCounter()], RegistersCount)
		dir := int32(c.Brain.Registers[dirReg]%8) * 45

		pos := c.getRelPos(object.Rotation{Degree: dir})
		if pos.Y < 0 || pos.Y >= int32(c.World.Height) {
			c.Brain.Registers[dest] = OBJ_WALL
			c.incCounter()
			return
		}

		o := c.World.GetObjectAtPosition(pos)
		if o == nil {
			c.Brain.Registers[dest] = OBJ_EMPTY
		} else {
			switch v := o.(type) {
			case *Cell:
				if v.IsDied() {
					c.Brain.Registers[dest] = OBJ_DEAD
				} else {
					if c.IsReleated(v) {
						c.Brain.Registers[dest] = OBJ_RELATED
					} else {
						c.Brain.Registers[dest] = OBJ_UNRELATED
					}
				}
			default:
				c.Brain.Registers[dest] = OBJ_EMPTY
			}
		}

		c.incCounter()
	}, true},
	// Bite another cell
	CMD_BITE: {func(c *Cell) {
		dirReg := truncCmd(c.Genome.Code[c.incCounter()], RegistersCount)
		dir := int32(c.Brain.Registers[dirReg]%8) * 45

		pos := c.getRelPos(object.Rotation{Degree: dir})
		if pos.Y < 0 || pos.Y >= int32(c.World.Height) {
			c.Brain.CompareFlag = CND_FAIL
			c.incCounter()
			return
		}

		o := c.World.GetObjectAtPosition(pos)
		if o == nil {
			c.Brain.CompareFlag = CND_FAIL
			c.incCounter()
			return
		}

		other, ok := o.(object.Lively)
		if !ok {
			c.Brain.CompareFlag = CND_FAIL
			c.incCounter()
			return
		}

		biteStrength := int(math.Round(BaseBiteStrength + float64(c.Weight)))
		if biteStrength > 255 {
			biteStrength = 255
		}

		c.IncreaseEnergy(other.Bite(byte(biteStrength)))
		c.Brain.CompareFlag = CND_SUCCESS
		c.incCounter()
	}, true},
	// Share energy with near cell
	CMD_SHAREENERGY: {func(c *Cell) {
		dirReg := truncCmd(c.Genome.Code[c.incCounter()], RegistersCount)
		dir := int32(c.Brain.Registers[dirReg]%8) * 45

		pos := c.getRelPos(object.Rotation{Degree: dir})
		if pos.Y < 0 || pos.Y >= int32(c.World.Height) {
			c.Brain.CompareFlag = CND_FAIL
			c.incCounter()
			return
		}

		o := c.World.GetObjectAtPosition(pos)
		if o == nil {
			c.Brain.CompareFlag = CND_FAIL
			c.incCounter()
			return
		}

		other, ok := o.(object.Lively)
		if !ok {
			c.Brain.CompareFlag = CND_FAIL
			c.incCounter()
			return
		}

		shenergy := BaseReproduceEnergyCost / 2
		c.SpendEnergy(byte(shenergy + BaseEnergyDecrement))

		if int(c.Energy) < shenergy {
			shenergy = int(c.Energy)
		}

		other.IncreaseEnergy(byte(shenergy))
		c.Brain.CompareFlag = CND_SUCCESS
		c.incCounter()
	}, true},
	// Recycle stuff to energy
	CMD_RECYCLE: {func(c *Cell) {
		recycleType := truncCmd(c.Genome.Code[c.incCounter()], RCL_ENUM_SIZE)
		ok := c.recycle(recycleType)
		if !ok {
			c.Brain.CompareFlag = CND_FAIL
			c.incCounter()
			return
		}
		c.Brain.CompareFlag = CND_SUCCESS
		c.incCounter()
	}, true},
	// Reproduce
	CMD_REPRODUCE: {func(c *Cell) {
		dirReg := truncCmd(c.Genome.Code[c.incCounter()], RegistersCount)
		dir := int32(c.Brain.Registers[dirReg]%8) * 45
		ok := c.Reproduce(object.Rotation{Degree: dir})
		if !ok {
			c.Brain.CompareFlag = CND_FAIL
			c.incCounter()
			return
		}
		c.Brain.CompareFlag = CND_SUCCESS
		c.incCounter()
	}, true},

	// Pickup something
	CMD_PICKUP: {func(c *Cell) {
		dirReg := truncCmd(c.Genome.Code[c.incCounter()], RegistersCount)
		dir := int32(c.Brain.Registers[dirReg]%8) * 45

		if c.BagageFullness >= BagageSize {
			c.Brain.CompareFlag = CND_FAIL
			c.incCounter()
			return
		}

		pos := c.getRelPos(object.Rotation{Degree: dir})
		if pos.Y < 0 || pos.Y >= int32(c.World.Height) {
			c.Brain.CompareFlag = CND_FAIL
			c.incCounter()
			return
		}

		o := c.World.GetObjectAtPosition(pos)
		if o == nil {
			c.Brain.CompareFlag = CND_FAIL
			c.incCounter()
			return
		}

		other, ok := o.(object.Pickable)
		if !ok {
			c.Brain.CompareFlag = CND_FAIL
			c.incCounter()
			return
		}

		ok = other.PickUp()
		if !ok {
			c.Brain.CompareFlag = CND_FAIL
			c.incCounter()
			return
		}

		c.Brain.CompareFlag = CND_SUCCESS
		c.Bagage[c.BagageFullness] = other
		c.BagageFullness++
		c.incCounter()
	}, true},
	// Drop selected item from bag
	CMD_DROP: {func(c *Cell) {
		dirReg := truncCmd(c.Genome.Code[c.incCounter()], RegistersCount)
		dir := int32(c.Brain.Registers[dirReg]%8) * 45

		if c.BagageFullness == 0 || c.Bagage[c.BagageSelected] == nil {
			c.Brain.CompareFlag = CND_FAIL
			c.incCounter()
			return
		}

		pos := c.getRelPos(object.Rotation{Degree: dir})
		if pos.Y < 0 || pos.Y >= int32(c.World.Height) {
			c.Brain.CompareFlag = CND_FAIL
			c.incCounter()
			return
		}

		o := c.World.GetObjectAtPosition(pos)
		if o != nil {
			c.Brain.CompareFlag = CND_FAIL
			c.incCounter()
			return
		}

		other, ok := o.(object.Pickable)
		if !ok {
			c.Brain.CompareFlag = CND_FAIL
			c.incCounter()
			return
		}

		ok = other.Drop(pos)
		if !ok {
			c.Brain.CompareFlag = CND_FAIL
			c.incCounter()
			return
		}

		c.Brain.CompareFlag = CND_SUCCESS
		c.Bagage[c.BagageSelected] = nil
		c.incCounter()
	}, true},
	// Count items in bag
	CMD_BAGSIZE: {func(c *Cell) {
		reg := truncCmd(c.Genome.Code[c.incCounter()], RegistersCount)
		c.Brain.Registers[reg] = byte(c.BagageFullness)
		c.incCounter()
	}, false},
	// Set active item in bag
	CMD_BAGACTIVE: {func(c *Cell) {
		c.BagageSelected = uint32(truncCmd(c.Genome.Code[c.incCounter()], BagageSize))
		c.incCounter()
	}, false},
	// Get ebergy of selected item in bag
	CMD_BAGENERGY: {func(c *Cell) {
		reg := truncCmd(c.Genome.Code[c.incCounter()], RegistersCount)
		if c.BagageFullness == 0 || c.Bagage[c.BagageSelected] == nil {
			c.Brain.Registers[reg] = 0
			c.incCounter()
			return
		}

		other, ok := c.Bagage[c.BagageSelected].(object.Object)
		if !ok {
			c.Brain.Registers[reg] = 0
			c.incCounter()
			return
		}

		c.Brain.Registers[reg] = other.GetEnergy()
		c.incCounter()
	}, false},
	// Get type of selected item in bag
	CMD_BAGCHECK: {func(c *Cell) {
		reg := truncCmd(c.Genome.Code[c.incCounter()], RegistersCount)
		if c.BagageFullness == 0 || c.Bagage[c.BagageSelected] == nil {
			c.Brain.Registers[reg] = OBJ_EMPTY
			c.incCounter()
			return
		}

		switch other := c.Bagage[c.BagageSelected].(type) {
		case object.Lively:
			if other.IsDied() {
				c.Brain.Registers[reg] = OBJ_DEAD
			} else {
				c.Brain.Registers[reg] = OBJ_BODY
			}
		default:
			c.Brain.Registers[reg] = OBJ_EMPTY
		}

		c.incCounter()
	}, false},

	// Get self age/10 and write to register
	CMD_GETAGE: {func(c *Cell) {
		dest := truncCmd(c.Genome.Code[c.incCounter()], RegistersCount)
		c.Brain.Registers[dest] = byte(c.GetAge() / 10)
		c.incCounter()
	}, false},
	// Get self health and write to register
	CMD_GETHEALTH: {func(c *Cell) {
		dest := truncCmd(c.Genome.Code[c.incCounter()], RegistersCount)
		c.Brain.Registers[dest] = c.GetHealth()
		c.incCounter()
	}, false},
	// Get self energy and write to register
	CMD_GETENERGY: {func(c *Cell) {
		dest := truncCmd(c.Genome.Code[c.incCounter()], RegistersCount)
		c.Brain.Registers[dest] = c.GetEnergy()
		c.incCounter()
	}, false},
	// Get self command counter and write to register
	CMD_GETCOUNTER: {func(c *Cell) {
		dest := truncCmd(c.Genome.Code[c.incCounter()], RegistersCount)
		c.Brain.Registers[dest] = byte(c.Brain.CommandCounter + 1)
		c.incCounter()
	}, false},
}

func (c *Cell) executeCommand(cmd Command) {
	cmdDesc, exists := commandMap[cmd]
	if !exists {
		return
	}
	cmdDesc.handler(c)
}

func (c *Cell) handleCommand(cmd Command) bool {
	cmdDesc, exists := commandMap[cmd]
	if !exists {
		return false
	}

	if cmdDesc.synchronous {
		return true
	}

	cmdDesc.handler(c)
	return false
}

func (c *Cell) currentCommad() Command {
	return c.Genome.Code[c.Brain.CommandCounter]
}
