package service

import (
	"encoding/json"
	"fmt"
	"github.com/fatedier/frp/cmd/frpc/entity"
	"github.com/fatedier/frp/cmd/frpc/myConfig"
	"github.com/fatedier/frp/cmd/frpc/sub"
	"github.com/fatedier/frp/cmd/frpc/util"
	"github.com/fatedier/frp/pkg/config"
	"github.com/fatedier/frp/pkg/util/log"
	"strconv"
)

func CheckVersion() {
	log.Info("********************初始化配置中********************")
	log.Info("开始连接服务器")
	log.Info("检测版本状态")
	log.Info("等待服务器响应")
	var res = util.MyHttpUtil("https://chrelyonly.cn/blog/frpApi/checkVersion", "get")
	if res.Code == 200 {
		version := "1.0"
		log.Info("本地frp管理版本: " + version)
		log.Info("远端frp管理版本: " + res.Msg)
		if version == res.Msg {
			log.Info("********************当前是最新版本,无需更新********************")
		} else {
			log.Info("********************发现新版本********************")
		}
		var input = ""
		var flag = true
		for flag {
			fmt.Println("请重新输入token")
			_, err := fmt.Scanln(&input)
			if err != nil {
				log.Info("输入错误,请重新输入")
				return
			}
			//发送请求获取配置文件
			var configRes = util.MyHttpUtil("http://chrelyonly/blog/frpApi/getFrpConfig?token="+input, "get")
			if configRes.Code == 200 {
				//启动frp
				flag = false
				data, _ := json.Marshal(&configRes.Data)
				var frpConfig = entity.FrpConfig{}
				err := json.Unmarshal(data, &frpConfig)
				if err != nil {
					log.Info("********************解析配置文件失败********************")
					return
				}
				log.Info("********************获取配置文件成功********************")
				log.Info("********************准备启动frp********************")
				myConfig.ServerAddr = frpConfig.ServerAddr
				myConfig.ServerPort, _ = strconv.Atoi(frpConfig.ServerPort)
				myConfig.Token = frpConfig.Token
				//循环处理配置文件
				exeJsonList(frpConfig.ListMap)
				sub.Execute()
			} else {
				log.Info("********************获取配置文件失败********************")
				log.Info(configRes.Msg)
			}
		}
		fmt.Println("程序退出")
	} else {
		log.Info("********************获取版本时异常********************")
		log.Info("********************" + res.Msg + "!!!********************")
	}
}

// 处理代理配置文件
func exeJsonList(list []entity.ListMap) {
	var tempFrpConFigList = make(map[string]config.ProxyConf)
	myConfig.FrpcConfig = list
	for i := 0; i < len(list); i++ {
		if list[i].Type == "tcp" {
			//可能有多个tcp代理
			info := config.TCPProxyConf{}
			info.LocalIP = list[i].LocalIp
			info.LocalPort, _ = strconv.Atoi(list[i].LocalPort)
			info.ProxyName = list[i].Comment
			info.RemotePort, _ = strconv.Atoi(list[i].RemotePort)
			tempFrpConFigList[list[i].Comment] = &info
		}
		if list[i].Type == "http" {
			//可能有多个tcp代理
			info := config.HTTPProxyConf{}
			info.LocalIP = list[i].LocalIp
			info.LocalPort, _ = strconv.Atoi(list[i].LocalPort)
			info.ProxyName = list[i].Comment
			info.CustomDomains = []string{list[i].CustomDomains}
			tempFrpConFigList[list[i].Comment] = &info
		}
	}
	config.ListMap = tempFrpConFigList

}
