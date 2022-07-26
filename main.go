package main

import (
	"gopher-dish/cell"
	"gopher-dish/world"
)

func main() {
	myWorld := world.New(256, 64)
	myCell := cell.New(myWorld, nil)

	for {
		myCell.Handle()
	}
}
