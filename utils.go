package gbase

import (
	"crypto/md5"
	"fmt"
	"math/rand"
	"reflect"
	"strings"
	"time"
)

func Sn(len int) string {
	t := NewTime()
	return fmt.Sprintf("%s%d%s", t.Cb.Format("YmdHis"), t.Cb.Microsecond(), GetRandomString(len))
}

func GetRandomString(lens int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < lens; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

func InArray(item interface{}, arr interface{}) bool {
	switch reflect.TypeOf(arr).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(arr)
		len := s.Len()
		for i := 0; i < len; i++ {
			v := s.Index(i).Interface()
			//logs.Debug(v)
			if item == v {
				return true
			}
		}
	}

	return false
}

func RandRatio[T TNum](percent, denominator T) bool {
	var randNum T
	// 生成一个0到denominator之间的随机数，对于浮点数，需要生成0到1之间的随机数再乘以denominator
	r := rand.Float64()
	randNum = T(r * float64(denominator))
	return randNum > 0 && percent >= randNum
}

func Password(len int, oriPwd string) (pwd string, salt string) {
	salt = GetRandomString(len)
	defaultPwd := "123456"
	if oriPwd != "" {
		defaultPwd = oriPwd
	}
	pwd = Md5([]byte(defaultPwd + salt))
	return pwd, salt
}

// Md5函数接受一个泛型参数T，T可以是[]byte或string类型
func Md5[T TByteOrStr](input T, upCase ...bool) string {
	var data []byte
	switch v := any(input).(type) {
	case []byte:
		data = v
	case string:
		data = []byte(v)
	}

	if len(upCase) > 0 && upCase[0] {
		return strings.ToUpper(fmt.Sprintf("%X", md5.Sum(data)))
	} else {
		return strings.ToLower(fmt.Sprintf("%x", md5.Sum(data)))
	}
}
