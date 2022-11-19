// Copyright 2021 The frp Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package config

import (
	"bytes"
	"fmt"
	"github.com/fatedier/frp/cmd/frpc/myConfig"
	"os"
	"path/filepath"
)

func ParseClientConfig(filePath string) (
	cfg ClientCommonConf,
	pxyCfgs map[string]ProxyConf,
	visitorCfgs map[string]VisitorConf,
	err error,
) {
	var content []byte
	content, err = GetRenderedConfFromFile(filePath)
	if err != nil {
		return
	}
	configBuffer := bytes.NewBuffer(nil)
	configBuffer.Write(content)

	// Parse common section.
	cfg, err = UnmarshalClientConfFromIni(content)
	if err != nil {
		return
	}
	cfg.Complete()
	if err = cfg.Validate(); err != nil {
		err = fmt.Errorf("parse myConfig error: %v", err)
		return
	}

	// Aggregate proxy configs from include files.
	var buf []byte
	buf, err = getIncludeContents(cfg.IncludeConfigFiles)
	if err != nil {
		err = fmt.Errorf("getIncludeContents error: %v", err)
		return
	}
	configBuffer.WriteString("\n")
	configBuffer.Write(buf)

	// Parse all proxy and visitor configs.
	pxyCfgs, visitorCfgs, err = LoadAllProxyConfsFromIni(cfg.User, configBuffer.Bytes(), cfg.Start)
	if err != nil {
		return
	}
	return
}

func ParseClientNetConfig() (
	cfg ClientCommonConf,
	pxyCfgs map[string]ProxyConf,
	visitorCfgs map[string]VisitorConf,
	err error,
) {
	//定义服务器地址
	//服务器地址
	cfg.ServerAddr = myConfig.ServerAddr
	//服务器端口
	cfg.ServerPort = myConfig.ServerPort
	//服务器校验token
	cfg.Token = myConfig.Token
	//服务器认证方法
	cfg.AuthenticationMethod = "token"
	//服务器认证超时时间
	cfg.DialServerTimeout = 10
	//保持长链接
	cfg.DialServerKeepAlive = 7200
	//# 如果使用tcp流多路复用，默认值为true，必须与FRP相同
	cfg.TCPMux = true
	//#指定tcp mux的保持活动间隔。
	//#仅当tcp_mux为真时有效。
	cfg.TCPMuxKeepaliveInterval = 60
	//日志输出模式
	cfg.LogFile = "console"
	cfg.LogWay = "console"
	cfg.LogLevel = "info"
	//日志最大保存时间
	cfg.LogMaxDays = 3
	//客户端管理配置地址
	cfg.AdminAddr = "127.0.0.1"
	//客户端管理配置端口
	cfg.AdminPort = 0
	//# 将提前建立连接，默认值为零
	cfg.PoolCount = 5
	//#决定是否在首次登录失败时退出程序，否则继续重新登录到frps
	//#默认值为true
	cfg.LoginFailExit = true
	//#用于连接到服务器的通信协议
	//#现在它支持tcp、kcp和websocket，默认为tcp
	cfg.Protocol = "tcp"
	//#如果tls_enable为真，frpc将通过tls连接FRP
	cfg.TLSEnable = false
	//#默认情况下，如果启用tls，frpc将连接FRP和第一个自定义字节。
	//#如果DisableCustomTLSFirstByte为true，frpc将不发送该自定义字节。
	cfg.DisableCustomTLSFirstByte = false
	//连接心跳
	cfg.HeartbeatInterval = 10
	//连接心跳超时
	cfg.HeartbeatTimeout = 60
	//传输数据包大小
	cfg.UDPPacketSize = 1500
	//定义代理类型
	pxyCfgs = NewProxyConfFromNet()
	visitorCfgs = make(map[string]VisitorConf)
	return
}

func NewProxyConfFromNet() map[string]ProxyConf {
	var pxyCfgs map[string]ProxyConf
	if myConfig.FrpType == "tcp" {
		info := TCPProxyConf{}
		info.LocalIP = myConfig.LocalIp
		info.LocalPort = myConfig.LocalPort
		info.ProxyName = myConfig.Comment
		info.RemotePort = myConfig.RemotePort
		pxyCfgs = make(map[string]ProxyConf)
		pxyCfgs[myConfig.Comment] = &info
	}
	return pxyCfgs
}

// getIncludeContents renders all configs from paths.
// files format can be a single file path or directory or regex path.
func getIncludeContents(paths []string) ([]byte, error) {
	out := bytes.NewBuffer(nil)
	for _, path := range paths {
		absDir, err := filepath.Abs(filepath.Dir(path))
		if err != nil {
			return nil, err
		}
		if _, err := os.Stat(absDir); os.IsNotExist(err) {
			return nil, err
		}
		files, err := os.ReadDir(absDir)
		if err != nil {
			return nil, err
		}
		for _, fi := range files {
			if fi.IsDir() {
				continue
			}
			absFile := filepath.Join(absDir, fi.Name())
			if matched, _ := filepath.Match(filepath.Join(absDir, filepath.Base(path)), absFile); matched {
				tmpContent, err := GetRenderedConfFromFile(absFile)
				if err != nil {
					return nil, fmt.Errorf("render extra myConfig %s error: %v", absFile, err)
				}
				out.Write(tmpContent)
				out.WriteString("\n")
			}
		}
	}
	return out.Bytes(), nil
}
