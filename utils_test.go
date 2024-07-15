package gbase

import (
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"testing"
)

func TestRandRatio(t *testing.T) {
	for {
		result := RandRatio(0.5, 1.01)
		hlog.Debug(result)
	}
}

func TestMd5(t *testing.T) {
	md5 := Md5([]byte("123456"), true)
	hlog.Debug(md5)
}
