package main

import (
	"encoding/json"
	"github.com/fatedier/frp/cmd/frpc/sub"
	"github.com/fatedier/frp/pkg/util/log"
	"io/ioutil"
	"net/http"
	"sync"
)

var waitGroup sync.WaitGroup

func main() {
	log.Info("********************初始化配置中********************")
	log.Info("开始连接服务器")
	log.Info("检测版本状态")
	log.Info("等待服务器响应")
	get, err := http.Get("https://chrelyonly.cn/blog/frpApi/checkVersion")
	if err != nil {
		log.Error("链接服务器超时")
	}
	status := get.StatusCode
	if status != 200 {
		log.Info("链接服务器超时")
	}
	body, err := ioutil.ReadAll(get.Body)
	if err != nil {
		log.Info("读取错误")
		return
	}
	//判断版本
	var res R
	_ = json.Unmarshal(body, &res)
	if res.Code == 200 {
		version := "1.0"
		log.Info("本地frp管理版本: " + version)
		log.Info("远端frp管理版本: " + res.Msg)
		if version == res.Msg {
			log.Info("********************当前是最新版本,无需更新********************")
		} else {
			log.Info("********************发现新版本请更新********************")
		}
	} else {
		log.Info("********************获取版本时异常********************")
		log.Info("********************" + res.Msg + "!!!********************")
	}
	waitGroup.Add(2)
	//执行完之后启动web页面
	//go controller.InitWeb()
	go sub.Execute()
	waitGroup.Wait()
	log.Info("系统系统完成")
}

// R 转换结构
type R struct {
	Code    int32             `json:"code"`
	Data    map[string]string `json:"data"`
	Success bool              `json:"success"`
	Msg     string            `json:"msg"`
}
