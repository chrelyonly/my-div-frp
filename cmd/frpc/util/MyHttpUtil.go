package util

import (
	"encoding/json"
	"github.com/fatedier/frp/cmd/frpc/entity"
	"github.com/fatedier/frp/pkg/util/log"
	"io/ioutil"
	"net/http"
	"strconv"
)

func MyHttpUtil(url string, method string) (res entity.R) {
	get, err := http.Get(url)
	if method == "get" {

	} else if method == "post" {

	} else {
		return res
	}
	if err != nil {
		log.Error("链接服务器超时")
	}
	status := get.StatusCode
	if status != 200 {
		log.Info("链接服务器超时")
		log.Info("返回状态" + strconv.Itoa(status))
	}
	body, err := ioutil.ReadAll(get.Body)
	if err != nil {
		log.Info("读取错误")
		return
	}
	_ = json.Unmarshal(body, &res)
	return res
}
