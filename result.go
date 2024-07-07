package gbase

import "github.com/cloudwego/hertz/pkg/common/hlog"

// 通用泛型结果处理类
type ResultType any

type ResultStruct[T ResultType] struct {
	Result T
	Err    error
}

func Result[T ResultType](result T, err error) *ResultStruct[T] {
	return &ResultStruct[T]{Result: result, Err: err}
}

func (this *ResultStruct[T]) Get() T {
	if this.Err != nil {
		hlog.Errorf("ResultStruct get failed,err:%s", this.Err.Error())
		panic(this.Err)
	}
	return this.Result
}

func (this *ResultStruct[T]) GetOr(defaultValue T) T {
	if this.Err != nil {
		return defaultValue
	}
	return this.Result
}
