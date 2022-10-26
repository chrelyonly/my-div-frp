package service

import (
	"github.com/fatedier/frp/cmd/frpc/util"
	"github.com/fatedier/frp/pkg/util/log"
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
			log.Info("********************发现新版本,请更新********************")
		}
	} else {
		log.Info("********************获取版本时异常********************")
		log.Info("********************" + res.Msg + "!!!********************")
	}
}
