package main

import (
	"github.com/fatedier/frp/cmd/frpc/controller"
	"github.com/fatedier/frp/cmd/frpc/sub"
	"github.com/fatedier/frp/pkg/util/log"
	"net/http"
	"sync"
)

var waitGroup sync.WaitGroup

func main() {
	log.Info("开始连接服务器")
	log.Info("检测版本状态")
	log.Info("等待服务器响应")
	get, err := http.Get("https://chrelyonly.cn/blog/blogconfig/blogConfig")
	if err != nil {
		log.Error("链接服务器超时")
	}
	status := get.StatusCode
	if status != 200 {
		log.Error("链接服务器超时")
	}
	err = get.Body.Close()
	if err != nil {
		log.Error("关闭错误")
		return
	}
	waitGroup.Add(2)
	go controller.InitWeb()
	go sub.Execute()
	waitGroup.Wait()
	log.Info("系统系统完成")
	//执行完之后启动web页面

}
