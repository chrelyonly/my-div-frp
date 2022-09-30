// Copyright 2016 fatedier, fatedier@gmail.com
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

package main

import (
	"fmt"
	_ "github.com/fatedier/frp/assets/frpc"
	"github.com/fatedier/frp/cmd/frpc/sub"
	"github.com/fatedier/frp/pkg/util/log"
	"net/http"
)

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
	sub.Execute()
	//执行完之后启动web页面
	initWeb()
}

func initWeb() {
	http.HandleFunc("/hello", helloWord)
	err := http.ListenAndServe(":8889", nil)
	if err != nil {
		log.Error("启动服务器失败", err)
		return
	}
}

func helloWord(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprint(w, "hello.word")
	if err != nil {
		return
	}
}
