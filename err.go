package gbase

import (
	"fmt"
	"github.com/pkg/errors"
)

//type Err struct {
//	Id int
//	Msg string
//}
//
//
//var list = map[string]*Err{
//	"SERVER-BUSY":{40001,"服务器繁忙,请稍候再试"},
//	"SAVE-ERR":{40002,"保存失败"},
//}

var list = map[string]string{
	"ServerErr": "服务器繁忙,请稍候再试",
	"ParamsErr": "参数有误",
	"DbErr":     "数据错误",
	"SaveErr":   "保存失败",
}

func E(code string, mark ...interface{}) string {
	var (
		e  string
		ok bool
	)
	if e, ok = list[code]; !ok {
		e = ""
	}
	if len(mark) > 0 {
		e += fmt.Sprintf("[%v]", mark[0])
	}
	return e
}

func EServer(mark ...interface{}) string {
	return E("ServerErr", mark...)
}
func EParams(mark ...interface{}) string {
	return E("ParamsErr", mark...)
}
func EDb(mark ...interface{}) string {
	return E("DbErr", mark...)
}
func ESave(mark ...interface{}) string {
	return E("SaveErr", mark...)
}

func NErrf(format string, v ...interface{}) error {
	return errors.New(fmt.Sprintf(format, v...))
}
