package object

type Object interface {
	GetID() uint64
	Prepare()
	Handle(yearChanged, epochChanged bool)
}
