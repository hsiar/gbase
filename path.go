package gbase

type pathItem interface {
	int8 | int | int64 | string
}

// fmt1: ,1,2,3,4,
// fmt2: ,tony,lucy,scot,
type Path[T pathItem] string

func (this Path[T]) FromString(pathStr string) Path[T] {
	this = Path[T](pathStr)
	return this
}
func (this Path[T]) ToString() string {
	return string(this)
}

func (this Path[T]) RmHeadEnd() Path[T] {
	if len(this) > 0 {
		this = this[1:]
	}
	if len(this) > 0 {
		this = this[0 : len(this)-1]
	}
	return this
}

//func (this Path[T]) Vd() error {
//
//}
