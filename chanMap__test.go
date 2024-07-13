package gbase

import (
	"fmt"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"testing"
	"time"
)

func TestCChanMap_CreateChan(t *testing.T) {
	var (
		key1, key2 string
		err        error
		resp       *Resp
	)
	ccm := ChanMapInst()
	key1, err = ccm.CreateChan()
	hlog.Debug("first key", key1)
	if key2, err = ccm.CreateChan(); err == nil {
		hlog.Debug("second key:", key2)
		go func() {
			if resp, err = ccm.SyncGet(key2, 5000); err != nil {
				hlog.Debug(err)
			} else {
				hlog.Debugf("get resp:%s", resp.ToString(resp))
			}
		}()

	}
	go func() {
		var r = &Resp{}
		r.Code = 200
		r.Msg = "ok"
		r.Data = 123
		time.Sleep(time.Second * 3)
		ccm.PushData(key2, r)
		//ccm[key2] <- a1.ToString(a1) //Map{"fuck": "you"}
	}()

	go func() {

		for {
			hlog.Debugf("chanmap.list:%d", ccm.Size())
			time.Sleep(time.Second * 1)
		}
	}()

	signalChan := make(chan os.Signal, 1)
	cleanupDone := make(chan bool)
	signal.Notify(signalChan,
		os.Interrupt,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	go func() {
		for _ = range signalChan {
			fmt.Println("\nReceived an interrupt, stopping services...")
			cleanupDone <- true
		}
	}()
	<-cleanupDone

	//globalservice.ProcessHangup()
}

func TestX1(t *testing.T) {
	//initProEnv()
	m := sync.Map{}
	m.Store(1, "1")
	m.Store(2, "2")
	hlog.Debug("%+v", m)
	m.Delete(2)
	hlog.Debug("%+v", m)

	type Ccm struct {
		m sync.Map
	}
	ccm := &Ccm{}
	ccm.m.Store(1, 1)
	ccm.m.Delete(1)
	//ccm.m.
	hlog.Debug(ccm)
}

//func TestXX(t *testing.T) {
//
//	//ccm := ChanMapInst()
//	//key, _ := ccm.CreateChanTest()
//	//ChanMapInst().DelChan(key)
//	//logs.Debug("%+v", ccm.list)
//
//	ccm := &ChanMap{}
//	ccm.list =
//	ccm.idNode, _ = snowflake.NewNode(1)
//	key, _ := ccm.CreateChan()
//	ccm.DelChan(key)
//
//	//key := ccm.idNode.Generate().Int64()
//	//ccm.list.Store(key, &Chan{data: make(chan any)})
//	//ccm.list.Delete(key)
//	logs.Debug(ccm)
//
//}
