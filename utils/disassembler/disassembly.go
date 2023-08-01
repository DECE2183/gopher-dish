package disassembler

import (
	"fmt"
	"gopher-dish/cell"
	"gopher-dish/utils"
)

type argType byte

const (
	_ARG_CONST = iota
	_ARG_REG
	_ARG_COND
)

var commandNames = map[cell.Command]string{
	cell.CMD_NOP: "nop",
	// Branch commands
	cell.CMD_CMP:  "cmp",
	cell.CMD_JMP:  "jmp",
	cell.CMD_DIVE: "dive",
	cell.CMD_LIFT: "lift",
	// Memory commands
	cell.CMD_PUT:  "put",
	cell.CMD_SAVE: "save",
	cell.CMD_LOAD: "load",
	// Math commands
	cell.CMD_ADD: "add",
	cell.CMD_SUB: "sub",
	cell.CMD_MUL: "mul",
	cell.CMD_DIV: "div",
	// Application commands
	cell.CMD_MOVE:      "move",
	cell.CMD_ROTATE:    "rot",
	cell.CMD_CHECKPOS:  "cpos",
	cell.CMD_BITE:      "bite",
	cell.CMD_RECYCLE:   "recl",
	cell.CMD_REPRODUCE: "repr",
	// Bagage commands
	cell.CMD_PICKUP:    "pick",
	cell.CMD_DROP:      "drop",
	cell.CMD_BAGSIZE:   "bsize",
	cell.CMD_BAGACTIVE: "bset",
	cell.CMD_BAGENERGY: "bnrg",
	cell.CMD_BAGCHECK:  "bchk",
	// Stats commands
	cell.CMD_GETAGE:     "age",
	cell.CMD_GETHEALTH:  "heal",
	cell.CMD_GETENERGY:  "nrg",
	cell.CMD_GETCOUNTER: "cntr",
}

var commandArgs = map[cell.Command][]argType{
	cell.CMD_NOP: {},
	// Branch commands
	cell.CMD_CMP:  {_ARG_REG, _ARG_REG},
	cell.CMD_JMP:  {_ARG_COND, _ARG_CONST},
	cell.CMD_DIVE: {_ARG_COND, _ARG_CONST},
	cell.CMD_LIFT: {_ARG_COND},
	// Memory commands
	cell.CMD_PUT:  {_ARG_REG, _ARG_CONST},
	cell.CMD_SAVE: {_ARG_REG, _ARG_CONST},
	cell.CMD_LOAD: {_ARG_REG, _ARG_CONST},
	// Math commands
	cell.CMD_ADD: {_ARG_REG, _ARG_REG, _ARG_REG},
	cell.CMD_SUB: {_ARG_REG, _ARG_REG, _ARG_REG},
	cell.CMD_MUL: {_ARG_REG, _ARG_REG, _ARG_REG},
	cell.CMD_DIV: {_ARG_REG, _ARG_REG, _ARG_REG},
	// Application commands
	cell.CMD_MOVE:      {_ARG_REG},
	cell.CMD_ROTATE:    {_ARG_REG},
	cell.CMD_CHECKPOS:  {_ARG_REG, _ARG_REG},
	cell.CMD_BITE:      {},
	cell.CMD_RECYCLE:   {_ARG_CONST},
	cell.CMD_REPRODUCE: {},
	// Bagage commands
	cell.CMD_PICKUP:    {},
	cell.CMD_DROP:      {},
	cell.CMD_BAGSIZE:   {},
	cell.CMD_BAGACTIVE: {},
	cell.CMD_BAGENERGY: {},
	cell.CMD_BAGCHECK:  {},
	// Stats commands
	cell.CMD_GETAGE:     {_ARG_REG},
	cell.CMD_GETHEALTH:  {_ARG_REG},
	cell.CMD_GETENERGY:  {_ARG_REG},
	cell.CMD_GETCOUNTER: {_ARG_REG},
}

var registerNames = map[cell.Command]string{
	cell.R0: "r0",
	cell.R1: "r1",
	cell.R2: "r2",
	cell.R3: "r3",
}

var conditionNames = map[cell.Command]string{
	cell.CND_NONE:   "_none",
	cell.Command(1): "_unk",

	cell.CND_EQ:    "_eq",
	cell.CND_NEQ:   "_neq",
	cell.CND_LESS:  "_less",
	cell.CND_GREAT: "_great",

	cell.CND_SUCCESS: "_succ",
	cell.CND_FAIL:    "_fail",

	cell.Command(128): "_unk",
}

func Disassembly(genome cell.Genome) (code string) {
	var cmditr utils.Iterator

	for cmditr < cell.GenomeLength {
		cmd := genome.Code[cmditr.Inc()]
		cmdName, ok := commandNames[cmd]
		if !ok {
			code += fmt.Sprintf("%-5s;\n", commandNames[cell.CMD_NOP])
			continue
		}
		code += fmt.Sprintf("%-5s", cmdName)
		for i, argt := range commandArgs[cmd] {
			switch argt {
			case _ARG_CONST:
				code += fmt.Sprint(genome.Code[cmditr.Inc()])
			case _ARG_REG:
				code += registerNames[genome.Code[cmditr.Inc()]%cell.RegistersCount]
			case _ARG_COND:
				cond := genome.Code[cmditr.Inc()]
				condCnt := 0
				for cind := 1; cind < 256; cind <<= 1 {
					if (int(cond) & cind) > 0 {
						code += conditionNames[cell.Command(cind)]
						code += " | "
					}
					condCnt++
				}
				if condCnt > 0 {
					code = code[:len(code)-3]
				}
			}

			if i < len(commandArgs[cmd])-1 {
				code += ", "
			}
		}
		code += ";\n"
	}

	code += "\n"
	return
}
