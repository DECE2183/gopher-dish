package main

import (
	"fmt"
	"gopher-dish/cell"
	"gopher-dish/gui"
	"gopher-dish/utils"
	"gopher-dish/utils/disassembler"
	"gopher-dish/world"
	"gopher-dish/world/worldsaver"
	"math/rand"
	"os"
	"time"
)

const (
	WorldTickInterval = 16 * time.Millisecond
	UITickInterval    = 33 * time.Millisecond
)

func main() {
	rand.Seed(time.Now().UnixNano())

	var i utils.Iterator
	for int(i) < len(os.Args) {
		arg := os.Args[i.Inc()]
		switch arg {
		case "-d", "--disassembly":
			path := os.Args[i.Inc()]

			f, err := os.Open(path)
			if err != nil {
				panic(err)
			}
			defer f.Close()

			w, err := worldsaver.Load(f)
			if err != nil {
				panic(err)
			}

			var (
				cellObj *cell.Cell
				ok      bool
			)

			for !ok {
				cellid := rand.Uint64() % uint64(w.ObjectsIdCounter)
				cellObj, ok = w.GetObject(cellid).(*cell.Cell)
			}

			code := disassembler.Disassembly(cellObj.Genome)
			fmt.Println(code)
			os.Exit(0)
		}
	}

	baseWorld := world.New(348, 128, WorldTickInterval)
	cell.New(baseWorld, nil)

	go func() {
		for {
			baseWorld.Handle()
		}
	}()

	gui.Run(UITickInterval, baseWorld)
}
