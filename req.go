package gbase

import (
	"bytes"
	"encoding/json"
	"github.com/beego/beego/v2/client/httplib"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"strings"
)

type Req struct {
	client  *httplib.BeegoHTTPRequest
	headers CMap
}

func (this *Req) WithHeaders(headers CMap) *Req {
	for k, _ := range headers {
		this.headers[k] = headers.GetString(k)
	}
	return this
}

func (this *Req) WithJsonHeader() *Req {
	this.headers["Accept"] = "application/json, text/javascript, */*; q=0.01"
	this.headers["Content-Type"] = "application/json"
	return this
}

func (this *Req) Send(method string, apiUrl string, params CMap) (resp *CResp) {
	var (
		err     error
		respStr string
	)
	//this.mu.Lock()
	//defer this.mu.Unlock()
	resp = &CResp{}

	if strings.ToLower(method) == "get" {
		this.client = httplib.Get(apiUrl)
		params.ToUrlParamsStr()
		for k, _ := range params {
			this.client.Param(k, params.GetString(k))
		}
		hlog.Debugf("Send GET API url:%s", this.client.GetRequest().URL)
	} else {
		this.client = httplib.Post(apiUrl)
		if params != nil {
			bf := bytes.NewBuffer([]byte{})
			jsonEncoder := json.NewEncoder(bf)
			jsonEncoder.SetEscapeHTML(false)
			_ = jsonEncoder.Encode(params)
			this.client.Body(bf.Bytes())
			hlog.Debugf("Send POST API url:%s,%s", apiUrl, params.ToString())
		}
	}
	for k, _ := range this.headers {
		this.client.Header(k, this.headers.GetString(k))
	}

	//Header("Accept", "application/json, text/javascript, */*; q=0.01").
	//Header("Content-Type", "application/json").
	//Header("")

	respStr, err = this.client.String()
	if err != nil {
		resp.Code = 100
		resp.Msg = "请求失败"
		return
	}
	hlog.Debugf("%s 响应:%s", apiUrl, respStr)
	resp.Code = 200
	resp.Msg = "ok"
	resp.Data = respStr
	return
}

func NewReq() *Req {
	req := &Req{}
	req.headers = CMap{}
	return req
}
