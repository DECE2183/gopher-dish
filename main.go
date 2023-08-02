package main

import (
	"fmt"
	"gopher-dish/cell"
	"gopher-dish/gui"
	"gopher-dish/utils"
	"gopher-dish/utils/genasm"
	"gopher-dish/world"
	"gopher-dish/world/worldsaver"
	"math/rand"
	"os"
	"strconv"
	"time"
)

const (
	WorldTickInterval = 16 * time.Millisecond
	UITickInterval    = 33 * time.Millisecond
)

func main() {
	var baseWorld *world.World
	var seed int64

	var i utils.Iterator
	for int(i) < len(os.Args) {
		arg := os.Args[i.Inc()]

		switch arg {
		case "-s", "--seed":
			num, err := strconv.ParseInt(os.Args[i.Inc()], 10, 64)
			if err != nil {
				panic(err)
			}

			seed = num

		case "-w", "--world":
			if len(os.Args) < int(i)+1 {
				fmt.Println("Missing path to the world file")
				os.Exit(22)
			}
			path := os.Args[i.Inc()]

			f, err := os.Open(path)
			if err != nil {
				panic(err)
			}
			defer f.Close()

			baseWorld, err = worldsaver.Load(f)
			if err != nil {
				panic(err)
			}

		case "-i", "--info":
			if baseWorld == nil {
				fmt.Println("You need to load the world first by '-w' or '--world' command")
				os.Exit(22)
			}

			fmt.Println("World info:")
			fmt.Printf("    Dimensions: [%d, %d]\n", baseWorld.Width, baseWorld.Height)
			fmt.Printf("    Tick:       %d\n", baseWorld.Ticks)
			fmt.Printf("    Year:       %d\n", baseWorld.Year)
			fmt.Printf("    Epoch:      %d\n", baseWorld.Epoch)
			fmt.Printf("    ID counter: %d\n", baseWorld.ObjectsIdCounter)
			fmt.Printf("    Population: %d\n", len(baseWorld.Objects))

		case "-d", "--disassembly":
			if baseWorld == nil {
				fmt.Println("You need to load the world first by '-w' or '--world' command")
				os.Exit(22)
			}

			num, err := strconv.ParseUint(os.Args[i.Inc()], 10, 64)
			if err != nil {
				panic(err)
			}

			var (
				cellObj *cell.Cell
				ok      bool
			)

			cellid := num % uint64(baseWorld.ObjectsIdCounter)
			cellObj, ok = baseWorld.GetObject(cellid).(*cell.Cell)

			if !ok {
				fmt.Printf("Cell with id %d not found", cellid)
				os.Exit(22)
			}

			code := genasm.Disassemble(cellObj.Genome)
			fmt.Println(code)

		case "-q", "--exit":
			os.Exit(0)
		}
	}

	if seed == 0 {
		seed = time.Now().UnixNano()
	}
	rand.Seed(seed)

	if baseWorld == nil {
		baseWorld = world.New(348, 128, WorldTickInterval)
		cell.New(baseWorld, nil)
	}

	go func() {
		for {
			baseWorld.Handle()
		}
	}()

	gui.Run(UITickInterval, baseWorld)
}
