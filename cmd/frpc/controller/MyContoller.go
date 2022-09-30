package controller

import (
	"fmt"
	"net/http"
)

func hello(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprint(w, "hello你好")
	if err != nil {
		fmt.Print(w, "hello你好")
		return
	}
}
