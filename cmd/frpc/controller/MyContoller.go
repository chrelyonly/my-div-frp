package controller

import (
	"fmt"
	"github.com/fatedier/frp/pkg/util/log"
	"net/http"
)

func initWeb() {
	http.HandleFunc("/hello", helloWord)
	log.Info("成功启动web服务器:8889")
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
