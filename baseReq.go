package gbase

type BaseReq struct {
	Base
}

func (this *BaseReq) IsAdd(id int64) bool {
	return id == 0
}
