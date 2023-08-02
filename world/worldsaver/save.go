package worldsaver

import (
	"bytes"
	"encoding/binary"
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

	for _, obj := range w.Objects {
		obj.Save(buf)
	}

	_, err = buf.WriteTo(writer)
	return
}
