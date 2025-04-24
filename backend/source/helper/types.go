package helper

import (
	"net/http"
	"time"
)

// channelData 用于在通道中传输数据的结构
type channelData struct {
	id   int
	data string
}

// 全局配置结构
type helperConfig struct {
	servers       []string // 允许的节点列表
	proxyPort     int      // 代理端口
	totalTimeout  int      // 批量请求总超时（秒）
	singleTimeout int      // 单次请求超时（秒）
}

var transport = &http.Transport{
	Proxy:                 http.ProxyFromEnvironment,
	ForceAttemptHTTP2:     true,
	MaxIdleConns:          100,
	IdleConnTimeout:       90 * time.Second,
	TLSHandshakeTimeout:   10 * time.Second,
	ExpectContinueTimeout: 1 * time.Second,
}
