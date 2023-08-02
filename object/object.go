package object

import (
	"io"
)

const (
	TYPE_CELL = iota
)

type Object interface {
	GetID() uint64
	Prepare()
	Handle(yearChanged, epochChanged bool)
	Save(writer io.Writer) error
}
