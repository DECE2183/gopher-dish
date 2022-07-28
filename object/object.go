package object

type Object interface {
	GetID() uint64
	GetInstance() Object

	Prepare()
	Handle(yearChanged, epochChanged bool)
}
