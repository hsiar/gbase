package gbase

import (
	"fmt"
	"github.com/emirpasic/gods/maps/treemap"
	"github.com/henrylee2cn/ameda"
	jsoniter "github.com/json-iterator/go"
)

type Map map[string]interface{}

func (this Map) FromString(str string) Map {
	_ = jsoniter.UnmarshalFromString(str, &this)
	return this
}

// from struct or map
func (this Map) FromX(v interface{}) error {
	if jsonBytes, err := jsoniter.Marshal(v); err != nil {
		return err
	} else if err = jsoniter.Unmarshal(jsonBytes, &this); err != nil {
		return err
	}
	return nil
}

func (this Map) ToString() string {
	str, _ := jsoniter.MarshalToString(this)
	return str
}

// {b:xx,a:xx,c:xx} => {a:xx,b:xx,c:xx}
func (this Map) ToSortString() string {
	tm := treemap.NewWithStringComparator()
	for k, v := range this {
		tm.Put(k, v)
	}
	if jb, err := tm.ToJSON(); err != nil {
		return ""
	} else {
		return string(jb)
	}
}

func (this Map) ToUrlParamsStr() string {
	var (
		str string
	)
	tm := treemap.NewWithStringComparator()
	for k, v := range this {
		tm.Put(k, v)
	}
	it := tm.Iterator()
	for it.Next() {
		str += fmt.Sprintf("%v=%v&", it.Key(), it.Value())
	}
	return String(str).StrRmEnd().ToString()
}

func (this Map) ToBytes() []byte {
	//jsoniter.
	bytes, _ := jsoniter.Marshal(this)
	return bytes
}

func (this Map) ToTreeMap() (tm *treemap.Map) {
	tm = treemap.NewWithStringComparator()
	for k, v := range this {
		tm.Put(k, v)
	}
	return
}

func (this Map) GetInt(key string) int {
	v, _ := ameda.StringToInt(fmt.Sprintf("%v", this[key]))
	return v
}

func (this Map) GetInt64(key string) int64 {
	//logs.Debug("this.key",key,this[key])
	var (
		value string
	)
	switch this[key].(type) {
	case float64, float32:
		value = fmt.Sprintf("%.f", this[key])
	case string:
		value = this[key].(string)
	case bool:
		value = fmt.Sprintf("%v", this[key])
	default:
		value = fmt.Sprintf("%d", this[key])
	}
	v, _ := ameda.StringToInt64(value)
	return v
}

func (this Map) GetString(key string) string {
	switch this[key].(type) {
	case float64, float32:
		return fmt.Sprintf("%.f", this[key])
	default:
		return fmt.Sprintf("%v", this[key])
	}
}

func (this Map) GetBool(key string) bool {
	switch this[key].(type) {
	case bool:
		return this[key].(bool)
	default:
		return false
	}
}

func (this Map) Exist(keys ...string) bool {
	for _, v := range keys {
		if _, exist := this[v]; !exist {
			return false
		}
	}
	return true
}

func (this Map) RemoveKeys(keys ...string) Map {
	for _, v := range keys {
		delete(this, v)
	}
	return this
}

func (this Map) Combine(cm Map) Map {
	for k, v := range cm {
		this[k] = v
	}
	return this
}
