package gbase

import (
	"github.com/cloudwego/hertz/pkg/app/server/binding"
	jsoniter "github.com/json-iterator/go"
	"strings"
)

// 全局基类
type Base struct { /*这里不能写字段名*/
}

func (this *Base) ToString(child interface{}) string {
	str, _ := jsoniter.MarshalToString(child)
	return str
}

func (this *Base) ToBytes(child interface{}) []byte {
	bytes, _ := jsoniter.Marshal(child)
	return bytes
}

func (this *Base) FromForm(child interface{}, formStr string) (err error) {
	list := strings.Split(formStr, "&")
	m := make(map[interface{}]interface{})
	if len(list) > 0 {
		for _, v := range list {
			cur := strings.Split(v, "=")
			if len(cur) == 2 {
				m[cur[0]] = cur[1]
			}
		}
	}
	if len(m) > 0 {
		if bytes, err := jsoniter.Marshal(m); err != nil {
			return err
		} else {
			if err = jsoniter.Unmarshal(bytes, child); err != nil {
				return err
			}
		}
	}
	return
}

func (this *Base) ToCMap(child interface{}) (cm CMap) {
	cm = CMap{}
	_ = cm.FromX(child)
	return
}

// from string,[]byte,map,struct
func (this *Base) FromX(child interface{}, params interface{}) error {
	var (
		jsonBytes []byte
		err       error
		//jsonStr string
	)
	switch params.(type) {
	case string:
		jsonBytes = String(params.(string)).ToBytes()
	case []byte:
		//do nothing
		jsonBytes = params.([]byte)
	default:
		if jsonBytes, err = jsoniter.Marshal(params); err != nil {
			return err
		}
	}
	if err = jsoniter.Unmarshal(jsonBytes, child); err != nil {
		return err
	}
	return nil
}

func (this *Base) Vd(child interface{}) error {
	return binding.Validate(child)
}
