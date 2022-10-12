package mem

import (
	"time"

	"github.com/fatedier/frp/pkg/util/metric"
)

const (
	ReserveDays = 7
)

type ServerStats struct {
	TotalTrafficIn  int64
	TotalTrafficOut int64
	CurConns        int64
	ClientCounts    int64
	ProxyTypeCounts map[string]int64
}

type ProxyStats struct {
	Name            string
	Type            string
	TodayTrafficIn  int64
	TodayTrafficOut int64
	LastStartTime   string
	LastCloseTime   string
	CurConns        int64
}

type ProxyTrafficInfo struct {
	Name       string
	TrafficIn  []int64
	TrafficOut []int64
}

type ProxyStatistics struct {
	Name          string
	ProxyType     string
	TrafficIn     metric.DateCounter
	TrafficOut    metric.DateCounter
	CurConns      metric.Counter
	LastStartTime time.Time
	LastCloseTime time.Time
}

type ServerStatistics struct {
	TotalTrafficIn  metric.DateCounter
	TotalTrafficOut metric.DateCounter
	CurConns        metric.Counter

	// counter for clients
	ClientCounts metric.Counter

	// counter for proxy types
	ProxyTypeCounts map[string]metric.Counter

	// statistics for different proxies
	// key is proxy name
	ProxyStatistics map[string]*ProxyStatistics
}

type Collector interface {
	GetServer() *ServerStats
	GetProxiesByType(proxyType string) []*ProxyStats
	GetProxiesByTypeAndName(proxyType string, proxyName string) *ProxyStats
	GetProxyTraffic(name string) *ProxyTrafficInfo
}
