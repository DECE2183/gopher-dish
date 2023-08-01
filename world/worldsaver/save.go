package worldsaver

import (
	"bytes"
	"encoding/binary"
	"gopher-dish/cell"
	"gopher-dish/world"
	"io"
)

func Save(w *world.World, writer io.Writer) (err error) {
	desc := wDescriptor{
		Width:         w.Width,
		Height:        w.Height,
		Ticks:         w.Ticks,
		Year:          w.Year,
		Epoch:         w.Epoch,
		ObjectCount:   uint64(len(w.Objects)),
		ObjectIdCount: w.ObjectsIdCounter,
	}

	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, desc)

	for id, obj := range w.Objects {
		switch o := obj.(type) {
		case *cell.Cell:
			binary.Write(buf, binary.LittleEndian, uint64(_TYPE_CELL))
			cdesc := wCellDescriptor{
				Id:           id,
				Generation:   o.Generation,
				ParentsChain: o.ParentsChain,
				Age:          o.Age,
				Health:       o.Health,
				Energy:       o.Energy,
				Weight:       o.Weight,
				Died:         o.Died,
				Picked:       o.Picked,
				Genome:       o.Genome,
				Brain:        o.Brain,
				Position:     o.Position,
				Rotation:     o.Rotation,
			}
			for i, bagage := range o.Bagage {
				if bagage == nil {
					continue
				}
				cdesc.Bagage[i] = bagage.GetID()
			}
			binary.Write(buf, binary.LittleEndian, cdesc)
		}
	}

	_, err = buf.WriteTo(writer)
	return
}
