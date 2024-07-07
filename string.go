package gbase

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
