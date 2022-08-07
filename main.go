package main

import (
	"fmt"
	"gopher-dish/cell"
	"gopher-dish/gui"
	"gopher-dish/world"
	"math/rand"
	"time"
)

const (
	WorldTickInterval = 15 * time.Millisecond
	UITickInterval    = 15 * time.Millisecond
)

var (
	WorldFramerate    int64
	WorldLastTickTime time.Time
	UIFramerate       int64
	UILastTickTime    time.Time
)

func main() {
	rand.Seed(time.Now().UnixNano())

	baseWorld := world.New(256, 128)
	cell.New(baseWorld, nil)
	fmt.Println()
	defer fmt.Println()

	go func() {
		WorldLastTickTime = time.Now()
		for {
			baseWorld.Handle()
			// time.Sleep(WorldTickInterval)

			WorldFramerate = 1000000 / (time.Since(WorldLastTickTime).Microseconds() + 1)
			// if baseWorld.Ticks%world.WorldTicksPerYear == 0 {
			// 	fmt.Printf("\t\t\t\t\t\t\t\r")
			// 	fmt.Printf("FPS: %d | Population: %d | Year: %d\r", WorldFramerate, len(baseWorld.Objects), baseWorld.Year)
			// }

			WorldLastTickTime = time.Now()
		}
	}()

	gui.Run(UITickInterval, baseWorld)
}
