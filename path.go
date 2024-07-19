package gbase

import (
	"fmt"
	"strconv"
	"strings"
)

type pathItem interface {
	int8 | int | int64 | string
}

// fmt1: ,1,2,3,4,
// fmt2: ,tony,lucy,scot,
type Path[T pathItem] string

func (this Path[T]) Empty() bool {
	return this == "" || this == ","
}

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

func (this Path[T]) ToList() (list []T, err error) {
	str := this.RmHeadEnd()
	if str == "" {
		return
	}
	strs := strings.Split(str.ToString(), ",")

	switch any(list).(type) {
	case []int64:
		for _, s := range strs {
			var v int64
			v, err = strconv.ParseInt(s, 10, 64)
			if err != nil {
				return nil, err
			}
			list = append(list, any(v).(T))
		}
	case []int8:
		for _, s := range strs {
			var v int8
			parsedInt, err := strconv.ParseInt(s, 10, 8)
			if err != nil {
				return nil, err
			}
			v = int8(parsedInt)
			list = append(list, any(v).(T))
		}
	default:
		return nil, fmt.Errorf("unsupported type")
	}

	return list, nil
}

func (this Path[T]) MustToList() (list []T) {
	var err error
	list, err = this.ToList()
	if err != nil {
		panic(fmt.Sprintf("Path.MustToList failed,err:%s", err.Error()))
	}
	return
}

func (this Path[T]) FromList(list []T) Path[T] {
	strList := make([]string, len(list))
	for i, v := range list {
		strList[i] = fmt.Sprintf("%v", v)
	}
	return Path[T](fmt.Sprintf(",%s,", strings.Join(strList, ",")))
}

func (this Path[T]) Has(v any) bool {
	return InArray(v, this.MustToList())
}

//func (this Path[T]) Vd() error {
//
//}
