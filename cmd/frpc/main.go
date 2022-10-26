package main

import (
	"github.com/fatedier/frp/cmd/frpc/service"
	"github.com/fatedier/frp/cmd/frpc/sub"
	"github.com/fatedier/frp/pkg/util/log"
	"sync"
)

var waitGroup sync.WaitGroup

func main() {
	//检查版本情况
	service.CheckVersion()
	//多线程
	waitGroup.Add(2)
	//执行完之后启动web页面
	//go controller.InitWeb()
	go sub.Execute()
	waitGroup.Wait()
	log.Info("系统启动完成")
}
