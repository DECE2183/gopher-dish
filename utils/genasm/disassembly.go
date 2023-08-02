package genasm

import (
	"fmt"
	"gopher-dish/cell"
	"gopher-dish/utils"
)

func Disassemble(genome cell.Genome) (code string) {
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
						condCnt++
					}
				}
				if condCnt > 0 {
					code = code[:len(code)-3]
				} else {
					code += conditionNames[0]
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
