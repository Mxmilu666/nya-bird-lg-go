package helper

import (
	"sync"

	"github.com/Mxmilu666/nya-bird-lg-go/source"
)

var (
	setting helperConfig
	once    sync.Once
)

// initConfig 初始化配置，确保只初始化一次
func initConfig() {
	once.Do(func() {
		setting = helperConfig{
			servers:       getServerIDs(),
			proxyPort:     source.AppConfig.LG.ProxyPort,
			singleTimeout: source.AppConfig.LG.Timeout,
		}
	})
}

// GetConfig 获取配置，确保配置已经初始化
func GetConfig() *helperConfig {
	initConfig()
	return &setting
}

// getServerIDs 从配置中获取所有服务器 ID
func getServerIDs() []string {
	if source.AppConfig == nil {
		return nil
	}
	var ids []string
	for _, server := range source.AppConfig.LG.Servers {
		ids = append(ids, server.ID)
	}
	return ids
}
