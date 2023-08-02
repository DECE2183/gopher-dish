package worldsaver

import (
	"encoding/binary"
	"fmt"
	"gopher-dish/cell"
	"gopher-dish/object"
	"gopher-dish/world"
	"io"
	"time"
)

func Load(reader io.Reader) (w *world.World, err error) {
	var desc wDescriptor
	err = binary.Read(reader, binary.LittleEndian, &desc)
	if err != nil {
		return
	}

	w = world.New(desc.Width, desc.Height, 16*time.Millisecond)

	w.Ticks = desc.Ticks
	w.Year = desc.Year
	w.Epoch = desc.Epoch
	w.ObjectsIdCounter = desc.ObjectIdCount

	for i := 0; i < int(desc.ObjectCount); i++ {
		var otype uint64
		err = binary.Read(reader, binary.LittleEndian, &otype)
		if err != nil {
			return
		}

		// fmt.Printf("load: %d/%d\r", i, int(desc.ObjectCount))

		switch otype {
		case object.TYPE_CELL:
			var cdesc wCellDescriptor
			err = binary.Read(reader, binary.LittleEndian, &cdesc)
			if err != nil {
				return
			}

			c := &cell.Cell{
				Name:         cdesc.Id,
				Generation:   cdesc.Generation,
				ParentsChain: cdesc.ParentsChain,
				Age:          cdesc.Age,
				Health:       cdesc.Health,
				Energy:       cdesc.Energy,
				Weight:       cdesc.Weight,
				Died:         cdesc.Died,
				Picked:       cdesc.Picked,
				Genome:       cdesc.Genome,
				Brain:        cdesc.Brain,
				Position:     cdesc.Position,
				Rotation:     cdesc.Rotation,
				World:        w,
			}

			w.Objects[cdesc.Id] = c
			w.Places[cdesc.Position.X][cdesc.Position.Y] = c
		default:
			err = fmt.Errorf("unknown object type 0x%X", otype)
			return
		}

	}

	// fmt.Printf("load: %d/%d\r\n", int(desc.ObjectCount), int(desc.ObjectCount))

	/*
		TODO: import cells bagage
		iterate all cells and assign bagage items by ID
	*/

	return
}
