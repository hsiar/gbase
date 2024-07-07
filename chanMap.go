package gbase

import (
	"github.com/bwmarrin/snowflake"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/emirpasic/gods/maps/treemap"
	"github.com/emirpasic/gods/utils"
	"github.com/go-errors/errors"
	"sync"
	"time"
)

var (
	ccmInst *ChanMap
)

type Chan struct {
	mu   sync.Mutex
	data chan any
}

func (this *Chan) pushData(data any) {
	this.mu.Lock()
	defer this.mu.Unlock()
	this.data <- data
}

// 通道字典,用于同步获取服务间响应数据，
// 比如machine模块发送websocket消息后，同步等待websocket响应消息数据
// 服务间响应数据通过rpc修改通道值返回
// 同服务内直接修改通道值返回
type ChanMap struct {
	idNode *snowflake.Node
	list   *treemap.Map
}

func (this *ChanMap) CreateChan() (key int64, err error) {
	key = this.idNode.Generate().Int64()
	this.list.Put(key, &Chan{data: make(chan any)})
	return
}

func (this *ChanMap) PushData(key int64, data any) {
	if this.Exist(key) {
		this.Get(key).pushData(data)
		//ori
		//this.Get(key).data <- data
	}

}

func (this *ChanMap) Exist(key int64) bool {
	_, ok := this.list.Get(key)
	return ok
}

func (this *ChanMap) Get(key int64) *Chan {
	item, _ := this.list.Get(key)
	return item.(*Chan)
}

func (this *ChanMap) DelChan(key int64) {
	this.list.Remove(key)
	//logs.Debug(this)

	//if this.Exist(key) {
	//	//cur := this.Get(key)
	//	//cur.mu.Lock()
	//	//close(cur.data)
	//	//cur.mu.Unlock()
	//	this.list.Delete(key)
	//	logs.Debug(this)
	//	//delete(this.list, key)
	//}
}

// v2 return CResp
// outTime,unit:ms
func (this *ChanMap) SyncGetV2(key int64, outTime ...time.Duration) (resp *CResp, err error) {
	var (
		realOutTime time.Duration
	)
	resp = &CResp{}
	if !this.Exist(key) {
		return nil, errors.New("not exist this chan")
	}

	if len(outTime) == 0 {
		realOutTime = 10 * 1000 //默认10秒
	} else {
		realOutTime = outTime[0]
	}

	select {
	case data := <-this.Get(key).data:
		if err = resp.FromX(data); err != nil {
			return nil, errors.New("SyncGetV2 failed,resp data is not *CResp type")
		}
		hlog.Debug("resp", resp.ToString(resp))
		this.DelChan(key)
		return resp, nil
	case <-time.After(realOutTime * time.Millisecond):
		// 超时处理
		this.DelChan(key)
		return nil, errors.New("Timeout waiting for response")
	}
}

func ChanMapInst() *ChanMap {
	if ccmInst == nil {
		var err error
		ccmInst = &ChanMap{}
		if ccmInst.idNode, err = snowflake.NewNode(1); err != nil {
			hlog.Error("ChanMapInst create snowflake node failed,err:%s", err.Error())
		}
		ccmInst.list = treemap.NewWith(utils.Int64Comparator)
	}
	return ccmInst
}
