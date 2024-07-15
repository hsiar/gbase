package gbase

import (
	"golang.org/x/exp/constraints"
)

//范型约束定义

type TInt interface {
	constraints.Integer
}

type TNum interface {
	constraints.Integer | constraints.Float
}

// 定义一个约束，要求T是[]byte或string类型
type TByteOrStr interface {
	~[]byte | ~string
}
