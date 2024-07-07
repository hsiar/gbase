package gbase

import (
	"fmt"
	"github.com/beego/beego/v2/client/httplib"
	"net/http"
	"net/url"
	"strings"
)

// 域名相关操作
type Url string

func (this Url) ToUrl() *url.URL {
	return Result(url.Parse(this.ToString())).Get()
}

func (this Url) ToString() string {
	return string(this)
}

func (this Url) Unescape() Url {
	n, _ := url.QueryUnescape(this.ToString())
	return Url(n)
}

//func (this Url) FromString(url string) Url {
//	return Url(url)
//}

func (this Url) Domain() (domain string) {
	u := this.ToUrl()
	if u != nil {
		domain = fmt.Sprintf("%s://%s", u.Scheme, u.Host)
	}
	return
}

func (this Url) GetHead() (resp *http.Response, err error) {
	resp, err = http.Head(this.ToString())
	defer resp.Body.Close()
	return
}

func (this Url) Get302Url() (url302 string, err error) {
	var head *http.Response
	if head, err = this.GetHead(); err == nil {
		url302 = head.Request.URL.String()
	}
	return
}

func (this Url) BuildReq() *httplib.BeegoHTTPRequest {
	return httplib.Get(this.ToString())
}

func (this Url) Get() (resp string, err error) {
	return httplib.Get(this.ToString()).String()
}

//func (this Url) PostJson(data CMap) (resp CMap, err error) {
//	var respStr string
//	respStr, err = httplib.Post(this.ToString()).
//		Header("Accept", "application/json, text/javascript, */*; q=0.01").
//		Header("Content-Type", "application/json").
//		Body(data.ToBytes()).String()
//	if err != nil {
//		return
//	}
//	return resp.FromString(respStr), nil
//}

func (this Url) GetPaths() []string {
	return strings.Split(this.ToUrl().Path, "/")
}

func (this Url) GetLastPath() string {
	paths := this.GetPaths()
	if len(paths) > 0 {
		return paths[len(paths)-1]
	}
	return ""
}

func (this Url) GetQueries() url.Values {
	return this.ToUrl().Query()
}

func (this Url) GetQuery(name string) string {
	return this.GetQueries().Get(name)
}

func (this Url) AddQuery(key, value string) Url {
	u := this.ToUrl()
	values := u.Query()
	values.Add(key, value)
	u.RawQuery = values.Encode()
	this = Url(u.String())
	return this
}
