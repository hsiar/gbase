package gbase

import (
	"fmt"
	"math/rand"
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
