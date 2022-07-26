package utils

type Iterator int

func (i *Iterator) Inc() (out int) {
	out = int(*i)
	*i++
	return
}

func (i *Iterator) Dec() (out int) {
	out = int(*i)
	*i--
	return
}
