package cell

import (
	"gopher-dish/utils"
	"io"
	"math/rand"
)

const (
	GenomeMutationRate = 2
)

type Genome struct {
	Hash uint64
	Code [GenomeLength]Command
}

func (g Genome) Read(out []byte) (n int, err error) {
	if len(out) < len(g.Code) {
		n = len(out)
	} else {
		n = len(g.Code)
	}

	for i := 0; i < n; i++ {
		out[i] = byte(g.Code[i])
	}

	err = io.EOF
	return
}

func (g Genome) Write(out []byte) (n int, err error) {
	if len(out) < len(g.Code) {
		n = len(out)
	} else {
		n = len(g.Code)
	}

	for i := 0; i < n; i++ {
		g.Code[i] = Command(out[i])
	}

	err = io.EOF
	return
}

func CreateBaseGenome() Genome {
	var newGenome Genome
	var i utils.Iterator

	for i < GenomeLength-32 {
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
		newGenome.Code[i.Inc()] = GenomeLength - 32
	}

	// Check if energy enough to reproduce
	newGenome.Code[i.Inc()] = CMD_GETENERGY
	newGenome.Code[i.Inc()] = R1
	newGenome.Code[i.Inc()] = CMD_PUT
	newGenome.Code[i.Inc()] = R2
	newGenome.Code[i.Inc()] = BaseReproduceEnergyCost * 3
	newGenome.Code[i.Inc()] = CMD_CMP
	newGenome.Code[i.Inc()] = R1
	newGenome.Code[i.Inc()] = R2
	// Return to programm if not enough energy
	newGenome.Code[i.Inc()] = CMD_LIFT
	newGenome.Code[i.Inc()] = CND_LESS | CND_EQ
	// Reproduce elsewise and jump back
	newGenome.Code[i.Inc()] = CMD_PUT
	newGenome.Code[i.Inc()] = R2
	newGenome.Code[i.Inc()] = Command(rand.Uint32() % 256)
	newGenome.Code[i.Inc()] = CMD_REPRODUCE
	newGenome.Code[i.Inc()] = R2
	newGenome.Code[i.Inc()] = CMD_LIFT
	newGenome.Code[i.Inc()] = CND_NONE
	// Jump to start
	newGenome.Code[i.Inc()] = CMD_JMP
	newGenome.Code[i.Inc()] = CND_NONE
	newGenome.Code[i.Inc()] = 0

	newGenome.Hash = genomeHash(newGenome.Code[:])

	return newGenome
}

func (g Genome) Mutate() Genome {
	for i := 0; i < GenomeMutationRate; i++ {
		g.Code[rand.Intn(GenomeLength)] = Command(rand.Intn(256))
	}

	g.Hash = genomeHash(g.Code[:])

	return g
}

func genomeHash(commands []Command) (v uint64) {
	for i, cmd := range commands {
		v += uint64(cmd&0x7F) << (i % 10)
	}
	return
}
