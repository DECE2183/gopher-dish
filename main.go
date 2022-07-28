package main

import (
	"fmt"
	"gopher-dish/cell"
	"gopher-dish/world"
	"math/rand"
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
	rand.Seed(time.Now().UnixNano())

	myWorld := world.New(256, 64)
	cell.New(myWorld, nil)

	WorldLastTickTime = time.Now()
	for {
		myWorld.Handle()
		// time.Sleep(WorldTickTime)

		WorldFramerate = 1000000 / (time.Since(WorldLastTickTime).Microseconds() + 1)
		if myWorld.Ticks%world.WorldTicksPerYear == 0 {
			fmt.Printf("FPS: %d\nPopulation: %d\nYear: %d\n\n", WorldFramerate, len(myWorld.Objects), myWorld.Year)
		}

		WorldLastTickTime = time.Now()
	}
}
