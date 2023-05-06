package main

import (
	"gopher-dish/cell"
	"gopher-dish/gui"
	"gopher-dish/world"
	"math/rand"
	"time"
)

const (
	WorldTickInterval = 16 * time.Millisecond
	UITickInterval    = 33 * time.Millisecond
)

func main() {
	rand.Seed(time.Now().UnixNano())

	baseWorld := world.New(348, 128, WorldTickInterval)
	cell.New(baseWorld, nil)

	go func() {
		for {
			baseWorld.Handle()
		}
	}()

	gui.Run(UITickInterval, baseWorld)
}
