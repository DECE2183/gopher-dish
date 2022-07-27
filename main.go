package main

import (
	"fmt"
	"gopher-dish/cell"
	"gopher-dish/world"
	"time"
)

const (
	WorldTickTime = 15 * time.Millisecond
)

var (
	WorldFramerate    int64
	WorldLastTickTime time.Time
)

func main() {
	myWorld := world.New(256, 64)
	cell.New(myWorld, nil)

	WorldLastTickTime = time.Now()
	for {
		myWorld.Handle()
		time.Sleep(WorldTickTime)

		WorldFramerate = 1000000 / time.Since(WorldLastTickTime).Microseconds()
		if myWorld.Ticks%world.WorldTicksPerYear == 0 {
			fmt.Printf("FPS: %d\nPopulation: %d\n\n", WorldFramerate, len(myWorld.Objects))
		}

		WorldLastTickTime = time.Now()
	}
}
