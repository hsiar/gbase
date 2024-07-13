package gbase

import (
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/hertz-contrib/sse"
)

type SSE struct {
	stream *sse.Stream
}

//params:msg,code,data
func (this *SSE) Pub(params ...interface{}) (err error) {
	var (
		msg  string
		code int = consts.StatusContinue
		data any
	)
	for k, v := range params {
		switch k {
		case 0:
			msg = v.(string)
		case 1:
			code = v.(int)
		case 2:
			data = v
		}
	}
	return this.stream.Publish(&sse.Event{
		Event: "fail",
		Data: Map{
			"code": code,
			"msg":  msg,
			"data": data,
		}.ToBytes(),
	})
}

func (this *SSE) Fail(params ...interface{}) (err error) {
	var (
		msg  string
		code int = consts.StatusContinue
	)
	for index, param := range params {
		switch index {
		case 0:
			msg = param.(string)
		case 1:
			code = param.(int)
		}
	}
	return this.stream.Publish(&sse.Event{
		Event: "fail",
		Data: Map{
			"code": code,
			"msg":  msg,
		}.ToBytes(),
	})
}

func (this *SSE) Success(params ...interface{}) (err error) {
	var (
		msg  string = "ok"
		data any
	)
	for index, param := range params {
		switch index {
		case 0:
			data = param
		case 1:
			msg = param.(string)
		}
	}
	return this.stream.Publish(&sse.Event{
		Event: "success",
		Data: Map{
			"code": consts.StatusOK,
			"msg":  msg,
			"data": data,
		}.ToBytes(),
	})
}
func NewSSE(c *app.RequestContext) (s *SSE) {
	s = &SSE{}
	s.stream = sse.NewStream(c)
	c.SetStatusCode(200)
	c.SetContentType("text/event-stream")
	return
}
