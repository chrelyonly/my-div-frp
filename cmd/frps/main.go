package main

import (
	"github.com/fatedier/frp/pkg/util/log"
	"math/rand"
	"time"

	"github.com/fatedier/golib/crypto"

	_ "github.com/fatedier/frp/assets/frps"
	_ "github.com/fatedier/frp/pkg/metrics"
)

func main() {
	crypto.DefaultSalt = "frp"
	rand.Seed(time.Now().UnixNano())
	log.Info("启动中...")
	Execute()
}
