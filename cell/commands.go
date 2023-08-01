package cell

import (
	"gopher-dish/object"
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
	CMD_PUT  // + put const to mem
	CMD_SAVE // + put reg value to mem
	CMD_LOAD // + load mem value to reg
	// Math commands
	CMD_ADD // +
	CMD_SUB // +
	CMD_MUL // +
	CMD_DIV // +
	// Application commands
	CMD_MOVE     // + move self cell
	CMD_ROTATE   // + rotate self cell
	CMD_CHECKPOS // + get type of object at near position
	CMD_BITE
	CMD_RECYCLE   // + recycle something to energy
	CMD_REPRODUCE // + reproduce
	// Bagage commands
	CMD_PICKUP    // - pickup something
	CMD_DROP      // - drop selected item from bag
	CMD_BAGSIZE   // - count items in bag
	CMD_BAGACTIVE // - set active item in bag
	CMD_BAGENERGY // - get ebergy of selected item in bag
	CMD_BAGCHECK  // - get type of selected item in bag
	// Stats commands
	CMD_GETAGE     // + get self age
	CMD_GETHEALTH  // + get self health
	CMD_GETENERGY  // + get self energy
	CMD_GETCOUNTER // + get current command counter

	CMD_ENUM_SIZE
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
		dir := (c.Brain.Registers[dirReg] % 8) * 45

		moveSuc := c.MoveInDirection(object.Rotation{Degree: int32(dir)})

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
		dir := (c.Brain.Registers[dirReg] % 8) * 45

		rotSuc := c.Rotate(object.Rotation{Degree: int32(dir)})

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
		dir := (c.Brain.Registers[dirReg] % 8) * 45

		pos := c.getRelPos(object.Rotation{Degree: int32(dir)})
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
	// Recycle stuff to energy
	CMD_RECYCLE: {func(c *Cell) {
		recycleType := truncCmd(c.Genome.Code[c.incCounter()], RCL_ENUM_SIZE)
		c.recycle(recycleType)
		c.incCounter()
	}, true},
	// Reproduce
	CMD_REPRODUCE: {func(c *Cell) {
		c.Reproduce()
		c.incCounter()
	}, true},

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
