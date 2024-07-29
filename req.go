package gbase

import (
	"bytes"
	"encoding/json"
	"github.com/beego/beego/v2/client/httplib"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"strings"
	"time"
)

type Req struct {
	client  *httplib.BeegoHTTPRequest
	headers Map
	url     string
	method  string

	//1.application/json 2.application/x-www-form-urlencoded 3....
	contentType int8

	timeout int
}

func (this *Req) WithHeaders(headers Map) *Req {
	for k, _ := range headers {
		this.headers[k] = headers.GetString(k)
	}
	return this
}

func (this *Req) WithJsonHeader() *Req {
	this.headers["Accept"] = "application/json, text/javascript, */*; q=0.01"
	this.headers["Content-Type"] = "application/json"
	this.contentType = 1
	return this
}

//func (this *Req) WithJsonHeader() *Req {
//	this.headers["Accept"] = "application/json, text/javascript, */*; q=0.01"
//	this.headers["Content-Type"] = "application/json"
//	this.contentType = 2
//	return this
//}

//func (this *Req) With(key string, value string) *Req {
//
//}

func (this *Req) WithUrl(url string) *Req {
	this.url = url
	return this
}

func (this *Req) WithMethod(method string) *Req {
	this.method = method
	return this
}

func (this *Req) WithTimeout(timeout int) *Req {
	this.timeout = timeout
	return this
}

func (this *Req) IsJsonContentType() bool {
	return this.contentType == 1
}

func (this *Req) IsGet() bool {
	return strings.ToUpper(this.method) == "GET"
}

func (this *Req) IsPost() bool {
	return strings.ToUpper(this.method) == "POST"
}

func (this *Req) GetClient() *httplib.BeegoHTTPRequest {
	return this.client
}

// 依赖this.url,this.method
func (this *Req) Build() *Req {
	this.client = httplib.NewBeegoRequest(this.url, strings.ToUpper(this.method))
	if this.timeout > 0 {
		this.client.SetTimeout(60*time.Second, time.Second*time.Duration(this.timeout))
	}
	return this
}

// 依赖this.Build
func (this *Req) PostFile(filePath string, key ...string) *Req {
	reqKey := "file"
	if len(key) > 0 {
		reqKey = key[0]
	}
	this.client.PostFile(reqKey, filePath)
	return this
}

func (this *Req) Do(params ...Map) (resp *Resp) {
	var (
		err        error
		respStr    string
		paramsData Map
	)
	//this.mu.Lock()
	//defer this.mu.Unlock()
	resp = &Resp{}
	if len(params) > 0 {
		paramsData = params[0]
	}

	switch {
	case this.IsGet():
		if this.client == nil {
			this.client = httplib.Get(this.url)
		}
		if paramsData != nil {
			for k, _ := range paramsData {
				this.client.Param(k, paramsData.GetString(k))
			}
		}
		hlog.Debugf("Send GET API url:%s", this.client.GetRequest().URL)
	case this.IsPost():
		if this.client == nil {
			this.client = httplib.Post(this.url)
		}

		if paramsData != nil {
			if this.IsJsonContentType() {
				bf := bytes.NewBuffer([]byte{})
				jsonEncoder := json.NewEncoder(bf)
				jsonEncoder.SetEscapeHTML(false)
				_ = jsonEncoder.Encode(paramsData)
				this.client.Body(bf.Bytes())
			} else {
				for k, _ := range paramsData {
					this.client.Param(k, paramsData.GetString(k))
				}
			}
			hlog.Debugf("Send POST API url:%s,%s", this.url, paramsData.ToString())
		}
	default:
		hlog.Errorf("不支持的请求类型:%s", this.method)
		return NewFailResp(100, "不支持的请求类型")
	}

	for k, _ := range this.headers {
		this.client.Header(k, this.headers.GetString(k))
	}

	respStr, err = this.client.String()
	hlog.Debugf("%s 响应:%s", this.url, respStr)
	if err != nil {
		return NewFailResp(100, "请求失败")
	}
	return NewSuccessResp(200, respStr, "ok")
}

// 简化版：NewReq().WithJsonHeader().Send(method,url,params)
func (this *Req) Send(method string, apiUrl string, params Map) (resp *Resp) {
	this.method = method
	this.url = apiUrl
	this.Build()
	return this.Do(params)
}

func NewReq() *Req {
	req := &Req{}
	req.headers = Map{}
	return req
}
