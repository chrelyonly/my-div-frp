package myConfig

import "github.com/fatedier/frp/cmd/frpc/entity"

// 是否网络模式
var IsNet bool

// frp服务器配置信息
var ServerAddr string
var ServerPort int
var Token string

// frp代理配置信息
var FrpcConfig []entity.ListMap
