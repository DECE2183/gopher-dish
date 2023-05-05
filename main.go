package main

import (
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

func main() {
	rand.Seed(time.Now().UnixNano())

	baseWorld := world.New(256, 128, WorldTickInterval)
	cell.New(baseWorld, nil)

	go func() {
		for {
			baseWorld.Handle()
		}
	}()

	gui.Run(UITickInterval, baseWorld)
}
