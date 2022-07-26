package cell

import (
	"gopher-dish/utils"
)

type Genome struct {
	Hash uint64
	Code [GenomeLength]Command
}

func CreateBaseGenome() Genome {
	var newGenome Genome
	var i utils.Iterator

	for i < GenomeLength-16 {
		// Recycle sun 3 times
		newGenome.Code[i.Inc()] = CMD_NOP
		newGenome.Code[i.Inc()] = CMD_RECYCLE
		newGenome.Code[i.Inc()] = RCL_SUNENERGY
		newGenome.Code[i.Inc()] = CMD_RECYCLE
		newGenome.Code[i.Inc()] = RCL_SUNENERGY
		newGenome.Code[i.Inc()] = CMD_RECYCLE
		newGenome.Code[i.Inc()] = RCL_SUNENERGY
		newGenome.Code[i.Inc()] = CMD_RECYCLE
		newGenome.Code[i.Inc()] = RCL_SUNENERGY
		newGenome.Code[i.Inc()] = CMD_RECYCLE
		newGenome.Code[i.Inc()] = RCL_SUNENERGY
		newGenome.Code[i.Inc()] = CMD_RECYCLE
		newGenome.Code[i.Inc()] = RCL_SUNENERGY
		// Jump to energy check block
		newGenome.Code[i.Inc()] = CMD_DIVE
		newGenome.Code[i.Inc()] = CND_NONE
		newGenome.Code[i.Inc()] = GenomeLength - 16
	}

	// Check if energy enough to reproduce
	newGenome.Code[i.Inc()] = CMD_GETENERGY
	newGenome.Code[i.Inc()] = R1
	newGenome.Code[i.Inc()] = CMD_PUT
	newGenome.Code[i.Inc()] = R2
	newGenome.Code[i.Inc()] = BaseEnergy * 3
	newGenome.Code[i.Inc()] = CMD_CMP
	newGenome.Code[i.Inc()] = R1
	newGenome.Code[i.Inc()] = R2
	// Return to programm if not enough energy
	newGenome.Code[i.Inc()] = CMD_LIFT
	newGenome.Code[i.Inc()] = CND_LESS | CND_EQ
	// Reproduce elsewise and jump back
	newGenome.Code[i.Inc()] = CMD_REPRODUCE
	newGenome.Code[i.Inc()] = CMD_LIFT
	newGenome.Code[i.Inc()] = CND_NONE
	// Jump to start
	newGenome.Code[i.Inc()] = CMD_JMP
	newGenome.Code[i.Inc()] = CND_NONE
	newGenome.Code[i.Inc()] = 0

	return newGenome
}

func (g Genome) Mutate() Genome {
	return g
}
