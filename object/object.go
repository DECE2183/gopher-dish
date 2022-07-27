package object

type Object interface {
	GetID() uint64
	GetInstance() Object
	Handle(yearChanged, epochChanged bool)
}
