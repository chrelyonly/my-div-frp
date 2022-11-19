package service

import (
	"fmt"
	"github.com/fatedier/frp/cmd/frpc/myConfig"
	"github.com/fatedier/frp/cmd/frpc/sub"
	"github.com/fatedier/frp/cmd/frpc/util"
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
		var input = "MJrxpq0gMrV8OAvB"
		var flag = true
		for flag {
			fmt.Println("请重新输入token")
			_, err := fmt.Scanln(&input)
			if err != nil {
				log.Info("输入错误,请重新输入,输入000退出")
			} else {
				if input == "000" {
					flag = false
				}
				//发送请求获取配置文件
				var configRes = util.MyHttpUtil("http://127.0.0.1:8082/frpApi/getFrpConfig?token="+input, "get")
				if configRes.Code == 200 {
					log.Info("********************获取配置文件成功********************")
					log.Info("********************开始启动frp********************")
					//启动frp
					flag = false
					myConfig.ServerAddr = configRes.Data["server_addr"]
					myConfig.ServerPort, _ = strconv.Atoi(configRes.Data["server_port"])
					myConfig.Token = configRes.Data["token"]
					myConfig.Comment = configRes.Data["comment"]
					myConfig.FrpType = configRes.Data["type"]
					myConfig.LocalIp = configRes.Data["local_ip"]
					myConfig.LocalPort, _ = strconv.Atoi(configRes.Data["local_port"])
					myConfig.RemotePort, _ = strconv.Atoi(configRes.Data["remote_port"])
					myConfig.CustomDomains = configRes.Data["custom_domains"]
					sub.Execute()
				} else {
					log.Info("********************获取配置文件失败********************")
					log.Info(configRes.Msg)
				}
			}
		}
		fmt.Println("程序退出")
	} else {
		log.Info("********************获取版本时异常********************")
		log.Info("********************" + res.Msg + "!!!********************")
	}
}
