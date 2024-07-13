package gbase

import (
	"fmt"
	"github.com/bwmarrin/snowflake"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/go-errors/errors"
	"sync"
	"time"
)

var (
	ccmInst *ChanMap
	once    sync.Once
)

type Chan struct {
	data chan any
}

func (this *Chan) pushData(data any) {
	this.data <- data
}

type ChanMap struct {
	idNode     *snowflake.Node
	list       sync.Map
	expireTime time.Duration
}

func (this *ChanMap) CreateChan(inputKey ...string) (key string, err error) {
	if len(inputKey) > 0 {
		key = inputKey[0]
	} else {
		key = this.idNode.Generate().String()
	}
	ch := &Chan{data: make(chan any)}
	this.list.Store(key, ch)
	go this.autoRemoveChan(key)
	return
}

func (this *ChanMap) autoRemoveChan(key string) {
	time.Sleep(this.expireTime)
	if _, ok := this.list.Load(key); ok {
		this.list.Delete(key)
		hlog.Warnf("Chan with key %s removed due to timeout", key)
	}
}

func (this *ChanMap) PushData(key string, data any) {
	if ch, ok := this.list.Load(key); ok {
		ch.(*Chan).pushData(data)
	}
}

func (this *ChanMap) Exist(key string) bool {
	_, ok := this.list.Load(key)
	return ok
}

func (this *ChanMap) Get(key string) *Chan {
	item, _ := this.list.Load(key)
	return item.(*Chan)
}

func (this *ChanMap) Del(key string) {
	this.list.Delete(key)
}

func (this *ChanMap) Size() (size int) {
	this.list.Range(func(key, value any) bool {
		size++
		return true
	})
	return
}

func (this *ChanMap) SyncGet(key string, outTime ...int64) (resp *Resp, err error) {
	var realOutTime time.Duration

	if !this.Exist(key) {
		return nil, errors.New("not exist this chan")
	}

	if len(outTime) == 0 {
		realOutTime = 10 * time.Second
	} else {
		realOutTime = time.Duration(outTime[0]) * time.Millisecond
	}

	select {
	case data := <-this.Get(key).data:
		resp = NewResp()
		if err = resp.FromX(data); err != nil {
			return nil, fmt.Errorf("SyncGet failed,resp data fmt error:%s", err.Error())
		}
		hlog.Debug("resp", resp.ToString(resp))
		this.Del(key)
		return resp, nil
	case <-time.After(realOutTime):
		this.Del(key)
		return nil, errors.New("Timeout waiting for response")
	}
}

func ChanMapInst() *ChanMap {
	once.Do(func() {
		var err error
		ccmInst = &ChanMap{expireTime: 10 * time.Minute} // 设置清理过期通道的时间，比如10分钟
		if ccmInst.idNode, err = snowflake.NewNode(1); err != nil {
			hlog.Error("ChanMapInst create snowflake node failed,err:%s", err.Error())
		}
	})
	return ccmInst
}
