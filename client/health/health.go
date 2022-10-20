package health

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/fatedier/frp/pkg/util/xlog"
)

var ErrHealthCheckType = errors.New("error health check type")

type Monitor struct {
	checkType      string
	interval       time.Duration
	timeout        time.Duration
	maxFailedTimes int

	// For tcp
	addr string

	// For http
	url string

	failedTimes    uint64
	statusOK       bool
	statusNormalFn func()
	statusFailedFn func()

	ctx    context.Context
	cancel context.CancelFunc
}

func NewMonitor(ctx context.Context, checkType string,
	intervalS int, timeoutS int, maxFailedTimes int,
	addr string, url string,
	statusNormalFn func(), statusFailedFn func(),
) *Monitor {
	if intervalS <= 0 {
		intervalS = 10
	}
	if timeoutS <= 0 {
		timeoutS = 3
	}
	if maxFailedTimes <= 0 {
		maxFailedTimes = 1
	}
	newctx, cancel := context.WithCancel(ctx)
	return &Monitor{
		checkType:      checkType,
		interval:       time.Duration(intervalS) * time.Second,
		timeout:        time.Duration(timeoutS) * time.Second,
		maxFailedTimes: maxFailedTimes,
		addr:           addr,
		url:            url,
		statusOK:       false,
		statusNormalFn: statusNormalFn,
		statusFailedFn: statusFailedFn,
		ctx:            newctx,
		cancel:         cancel,
	}
}

func (monitor *Monitor) Start() {
	go monitor.checkWorker()
}

func (monitor *Monitor) Stop() {
	monitor.cancel()
}

func (monitor *Monitor) checkWorker() {
	xl := xlog.FromContextSafe(monitor.ctx)
	for {
		doCtx, cancel := context.WithDeadline(monitor.ctx, time.Now().Add(monitor.timeout))
		err := monitor.doCheck(doCtx)

		// check if this monitor has been closed
		select {
		case <-monitor.ctx.Done():
			cancel()
			return
		default:
			cancel()
		}

		if err == nil {
			xl.Trace("心跳正常")
			if !monitor.statusOK && monitor.statusNormalFn != nil {
				xl.Info("运行状况检查状态更改为成功")
				monitor.statusOK = true
				monitor.statusNormalFn()
			}
		} else {
			xl.Warn("连接状态心跳获取失败: %v", err)
			monitor.failedTimes++
			if monitor.statusOK && int(monitor.failedTimes) >= monitor.maxFailedTimes && monitor.statusFailedFn != nil {
				xl.Warn("运行状况检查状态更改为失败")
				monitor.statusOK = false
				monitor.statusFailedFn()
			}
		}

		time.Sleep(monitor.interval)
	}
}

func (monitor *Monitor) doCheck(ctx context.Context) error {
	switch monitor.checkType {
	case "tcp":
		return monitor.doTCPCheck(ctx)
	case "http":
		return monitor.doHTTPCheck(ctx)
	default:
		return ErrHealthCheckType
	}
}

func (monitor *Monitor) doTCPCheck(ctx context.Context) error {
	// if tcp address is not specified, always return nil
	if monitor.addr == "" {
		return nil
	}

	var d net.Dialer
	conn, err := d.DialContext(ctx, "tcp", monitor.addr)
	if err != nil {
		return err
	}
	conn.Close()
	return nil
}

func (monitor *Monitor) doHTTPCheck(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, "GET", monitor.url, nil)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, resp.Body)

	if resp.StatusCode/100 != 2 {
		return fmt.Errorf("do http health check, StatusCode is [%d] not 2xx", resp.StatusCode)
	}
	return nil
}
