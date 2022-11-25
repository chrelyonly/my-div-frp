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
	httpRes, err := http.Get(url)
	if err != nil {
		log.Error("链接服务器超时")
	}
	status := httpRes.StatusCode
	if status != 200 {
		log.Info("链接服务器超时")
		log.Info("返回状态" + strconv.Itoa(status))
	}
	body, err := ioutil.ReadAll(httpRes.Body)
	if err != nil {
		log.Info("读取错误")
		return
	}
	_ = json.Unmarshal(body, &res)
	return res
}
