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
