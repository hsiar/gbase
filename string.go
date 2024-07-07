package gbase

import "unsafe"

type String string

func (this String) ToString() string {
	return string(this)
}

func (this String) StrRmEnd() String {
	if len(this) > 0 {
		return this[0 : len(this)-1]
	} else {
		return this
	}
}

// string与byte[]互转,性能优于直转
func (this String) ToBytes() []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&this))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}
