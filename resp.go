package gbase

import (
	jsoniter "github.com/json-iterator/go"
)

type Resp struct {
	Base
	Code   int         `json:"code"`
	Status int         `json:"status"`
	Msg    string      `json:"msg"`
	Data   interface{} `json:"data"`
}

func (this *Resp) WithCode(code int) *Resp {
	this.Code = code
	return this
}
func (this *Resp) WithMsg(msg string) *Resp {
	this.Msg = msg
	return this
}
func (this *Resp) WithData(data any) *Resp {
	this.Data = data
	return this
}

//func (this *Resp) FromString(str string) *Resp {
//	_ = jsoniter.UnmarshalFromString(str, this)
//	return this
//}

// from struct or map
func (this *Resp) FromX(v interface{}) error {
	var (
		jsonBytes []byte
		err       error
	)
	switch v.(type) {
	case string:
		jsonBytes = String(v.(string)).ToBytes()
	case []byte:
		//do nothing
		jsonBytes = v.([]byte)
	default:
		if jsonBytes, err = jsoniter.Marshal(v); err != nil {
			return err
		}
	}
	if err = jsoniter.Unmarshal(jsonBytes, this); err != nil {
		return err
	}
	return nil
}

func (this *Resp) IsSuccess() bool {
	return this.Code == 200 || this.Status == 200
}

// deprecated
// pData必须为指针
//func (this *Resp) DataTo(pData any) error {
//	if bData, err := jsoniter.Marshal(this.Data); err != nil {
//		return err
//	} else {
//		return jsoniter.Unmarshal(bData, pData)
//	}
//
//}

// 增加了this.data支持string,[]byte类型，原先this.data支持其它结构对象类型
func (this *Resp) DataTo(pData any) error {
	var (
		jsonBytes []byte
		err       error
	)
	switch this.Data.(type) {
	case string:
		jsonBytes = String(this.Data.(string)).ToBytes() //util.Str2bytes(this.Data.(string))
	case []byte:
		//do nothing
		jsonBytes = this.Data.([]byte)
	default:
		if jsonBytes, err = jsoniter.Marshal(this.Data); err != nil {
			return err
		}
	}
	if err = jsoniter.Unmarshal(jsonBytes, pData); err != nil {
		return err
	}

	return nil

}

func NewResp() (obj *Resp) {
	obj = &Resp{}
	return
}

func NewFailResp(code int, msg string, data ...any) (resp *Resp) {
	resp = NewResp().WithCode(code).WithMsg(msg)
	if len(data) > 0 {
		resp.WithData(data[0])
	}
	return
}

func NewSuccessResp(code int, data any, msg ...string) (resp *Resp) {
	resp = NewResp().WithCode(code).WithData(data)
	if len(msg) > 0 {
		resp.WithMsg(msg[0])
	}
	return
}
