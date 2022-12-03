package main

//var waitGroup sync.WaitGroup

import (
	"fmt"
	"github.com/fatedier/frp/cmd/frpc/myConfig"
	"github.com/fatedier/frp/cmd/frpc/service"
	"github.com/fatedier/frp/cmd/frpc/sub"
	"github.com/fatedier/frp/pkg/util/log"
)

func main() {

	var flag = true
	for flag {
		//监听键盘输入
		fmt.Println("输入1原生启动")
		fmt.Println("输入2通过网络获取配置文件")
		var input int
		_, err := fmt.Scanln(&input)
		if err != nil {
			log.Info("输入错误")
			return
		}
		switch input {
		case 1:
			//设置为离线模式
			myConfig.IsNet = false
			sub.Execute()
		case 2:
			//设置为网络模式
			myConfig.IsNet = true
			////检查版本情况
			service.CheckVersion()
		case 3:
			//退出
			flag = false
			fmt.Println("程序退出")
		default:
			fmt.Println("输入错误,请重新输入")
		}
	}

	////多线程
	//waitGroup.Add(2)
	////执行完之后启动web页面
	////go controller.InitWeb()
	//go sub.Execute()
	//waitGroup.Wait()
	//log.Info("系统启动完成")
}
