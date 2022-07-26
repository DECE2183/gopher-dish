package cell

import (
	"fmt"
	"gopher-dish/object"
)

func truncCmd(cmd Command, max uint64) uint64 {
	return uint64(cmd) % max
}

type Command byte
type CommandHandler func(*Cell)

// Commands list
const (
	// No operation command
	CMD_NOP = iota
	// Branch commands
	CMD_CMP  // +
	CMD_JMP  // +
	CMD_DIVE // +
	CMD_LIFT // +
	// Memory commands
	CMD_PUT  // +
	CMD_SAVE // +
	CMD_LOAD // +
	// Math commands
	CMD_ADD // +
	CMD_SUB // +
	CMD_MUL // +
	CMD_DIV // +
	// Application commands
	CMD_MOVE   // +
	CMD_ROTATE // +
	CMD_CHECKPOS
	CMD_BITE
	CMD_RECYCLE   // +
	CMD_REPRODUCE // +
	// Bagage commands
	CMD_PICKUP
	CMD_DROP
	CMD_BAGAGESIZE
	CMD_BAGAGESLOTENERGY
	CMD_SLOTOBJECTTYPE
	// Stats commands
	CMD_GETAGE     // +
	CMD_GETHEALTH  // +
	CMD_GETENERGY  // +
	CMD_GETCOUNTER // +

	CMD_ENUM_SIZE
)

// Compare conditions list
const (
	CND_NONE = iota

	CND_EQ    = 1 << 0
	CND_NEQ   = 1 << 1
	CND_LESS  = 1 << 2
	CND_GREAT = 1 << 3

	CND_SUCCESS = 1 << 6
	CND_FAIL    = 1 << 7

	CND_ENUM_SIZE = 256
)

// Object types list
const (
	CND_WALL = iota
	CND_BODY
	CND_DEAD

	CND_RELATED
	CND_UNRELATED
)

// Recycle targets list
const (
	RCL_NONE = iota
	RCL_SUNENERGY
	RCL_BAGAGE
	RCL_ENUM_SIZE
)

var commandMap = map[Command]CommandHandler{
	CMD_NOP: func(c *Cell) {
		c.incCounter()
	},

	// Compare values in registers and put result to CompareFlag
	CMD_CMP: func(c *Cell) {
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
	},
	// Jump according to condition in CompareFlag
	CMD_JMP: func(c *Cell) {
		cond := truncCmd(c.Genome.Code[c.incCounter()], CND_ENUM_SIZE)
		pos := truncCmd(c.Genome.Code[c.incCounter()], GenomeLength)

		if cond == CND_NONE || cond&uint64(c.Brain.CompareFlag) > 0 {
			c.Brain.CommandCounter = pos
		} else {
			c.incCounter()
		}
	},
	// Memorize current registers and command counter and dive into subprogramm
	CMD_DIVE: func(c *Cell) {
		cond := truncCmd(c.Genome.Code[c.incCounter()], CND_ENUM_SIZE)
		pos := truncCmd(c.Genome.Code[c.incCounter()], GenomeLength)

		if c.Brain.StackCounter >= StackDepth {
			c.incCounter()
			return
		}

		cnt := c.Brain.StackCounter
		c.Brain.Stack[cnt].JumpPosition = c.Brain.CommandCounter + 1
		c.Brain.Stack[cnt].JumpRegisters = c.Brain.Registers
		c.Brain.Stack[cnt].JumpCompareFlag = c.Brain.CompareFlag
		c.Brain.StackCounter++

		if cond == CND_NONE || cond&uint64(c.Brain.CompareFlag) > 0 {
			c.Brain.CommandCounter = pos
		} else {
			c.incCounter()
		}
	},
	// Return to main programm
	CMD_LIFT: func(c *Cell) {
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
	},

	// Put value to register
	CMD_PUT: func(c *Cell) {
		reg := truncCmd(c.Genome.Code[c.incCounter()], RegistersCount)
		val := byte(c.Genome.Code[c.incCounter()])

		c.Brain.Registers[reg] = val
		c.incCounter()
	},
	// Save value from register to memory
	CMD_SAVE: func(c *Cell) {
		reg := truncCmd(c.Genome.Code[c.incCounter()], RegistersCount)
		mem := truncCmd(c.Genome.Code[c.incCounter()], MemorySize)

		c.Brain.Memory[mem] = c.Brain.Registers[reg]
		c.incCounter()
	},
	// Load value from memory to registor
	CMD_LOAD: func(c *Cell) {
		reg := truncCmd(c.Genome.Code[c.incCounter()], RegistersCount)
		mem := truncCmd(c.Genome.Code[c.incCounter()], MemorySize)

		c.Brain.Registers[reg] = c.Brain.Memory[mem]
		c.incCounter()
	},

	// Add two regs command
	CMD_ADD: func(c *Cell) {
		dest := truncCmd(c.Genome.Code[c.incCounter()], RegistersCount)
		src := truncCmd(c.Genome.Code[c.incCounter()], RegistersCount)
		op := truncCmd(c.Genome.Code[c.incCounter()], RegistersCount)

		c.Brain.Registers[dest] = c.Brain.Registers[src] + c.Brain.Registers[op]
		c.incCounter()
	},
	// Subtract two regs command
	CMD_SUB: func(c *Cell) {
		dest := truncCmd(c.Genome.Code[c.incCounter()], RegistersCount)
		src := truncCmd(c.Genome.Code[c.incCounter()], RegistersCount)
		op := truncCmd(c.Genome.Code[c.incCounter()], RegistersCount)

		c.Brain.Registers[dest] = c.Brain.Registers[src] - c.Brain.Registers[op]
		c.incCounter()
	},
	// Multiply two regs command
	CMD_MUL: func(c *Cell) {
		dest := truncCmd(c.Genome.Code[c.incCounter()], RegistersCount)
		src := truncCmd(c.Genome.Code[c.incCounter()], RegistersCount)
		op := truncCmd(c.Genome.Code[c.incCounter()], RegistersCount)

		c.Brain.Registers[dest] = c.Brain.Registers[src] * c.Brain.Registers[op]
		c.incCounter()
	},
	// Divide two regs command
	CMD_DIV: func(c *Cell) {
		dest := truncCmd(c.Genome.Code[c.incCounter()], RegistersCount)
		src := truncCmd(c.Genome.Code[c.incCounter()], RegistersCount)
		op := truncCmd(c.Genome.Code[c.incCounter()], RegistersCount)

		c.Brain.Registers[dest] = c.Brain.Registers[src] / c.Brain.Registers[op]
		c.incCounter()
	},

	// Relative move command
	CMD_MOVE: func(c *Cell) {
		dirReg := truncCmd(c.Genome.Code[c.incCounter()], RegistersCount)
		dir := (dirReg % 8) * 45

		moveSuc := c.MoveInDirection(object.Rotation{Degree: int32(dir)})

		if moveSuc {
			c.Brain.CompareFlag = CND_SUCCESS
		} else {
			c.Brain.CompareFlag = CND_FAIL
		}
		c.incCounter()
	},
	// Relative rotation  command
	CMD_ROTATE: func(c *Cell) {
		dirReg := truncCmd(c.Genome.Code[c.incCounter()], RegistersCount)
		dir := (dirReg % 8) * 45

		rotSuc := c.Rotate(object.Rotation{Degree: int32(dir)})

		if rotSuc {
			c.Brain.CompareFlag = CND_SUCCESS
		} else {
			c.Brain.CompareFlag = CND_FAIL
		}
		c.incCounter()
	},
	// Recycle stuff to energy
	CMD_RECYCLE: func(c *Cell) {
		recycleType := truncCmd(c.Genome.Code[c.incCounter()], RCL_ENUM_SIZE)
		c.recycle(recycleType)
		c.incCounter()
	},
	// Reproduce
	CMD_REPRODUCE: func(c *Cell) {
		c.Reproduce()
		c.incCounter()
	},

	// Get self age/10 and write to register
	CMD_GETAGE: func(c *Cell) {
		dest := truncCmd(c.Genome.Code[c.incCounter()], RegistersCount)
		c.Brain.Registers[dest] = byte(c.GetAge() / 10)
		c.incCounter()
	},
	// Get self health and write to register
	CMD_GETHEALTH: func(c *Cell) {
		dest := truncCmd(c.Genome.Code[c.incCounter()], RegistersCount)
		c.Brain.Registers[dest] = c.GetHealth()
		c.incCounter()
	},
	// Get self energy and write to register
	CMD_GETENERGY: func(c *Cell) {
		dest := truncCmd(c.Genome.Code[c.incCounter()], RegistersCount)
		c.Brain.Registers[dest] = c.GetEnergy()
		c.incCounter()
	},
	// Get self command counter and write to register
	CMD_GETCOUNTER: func(c *Cell) {
		dest := truncCmd(c.Genome.Code[c.incCounter()], RegistersCount)
		c.Brain.Registers[dest] = byte(c.Brain.CommandCounter + 1)
		c.incCounter()
	},
}

func (c *Cell) handleCommand(cmd Command) {
	h, exists := commandMap[cmd]
	if !exists {
		_ = fmt.Errorf("Command doesn't exists: %v\n", cmd)
		return
	}

	h(c)
}

func (c *Cell) currentCommad() Command {
	return c.Genome.Code[c.Brain.CommandCounter]
}
