package gbase

import (
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"testing"
)

func TestPath_FromString(t *testing.T) {
	pt := Path[int](",1,2,")
	pt.RmHeadEnd()
	hlog.Debug(pt)
}

func TestPath_FromList(t *testing.T) {

	ids := []int64{1, 2, 3}
	path := Path[int64]("").FromList(ids)
	hlog.Debug(path)
}
