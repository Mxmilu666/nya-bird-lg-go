package main

import (
	"encoding/json"

	"github.com/Mxmilu666/nya-bird-lg-go/source"
	logger "github.com/Mxmilu666/nya-bird-lg-go/source/logger"
	"github.com/Mxmilu666/nya-bird-lg-go/source/server"
)

func main() {
	// 加载配置
	err := source.LoadConfig()
	if err != nil {
		panic("Failed to load config: " + err.Error())
	}

	// 初始化日志
	l := logger.InitLogger()
	if l == nil {
		panic("Failed to initialize logger")
	}

	l.Info("Nya!,Nya-Bird-LG")

	// 预览配置
	configJSON, err := json.Marshal(source.AppConfig)
	if err != nil {
		l.Error("Failed to convert configuration to JSON", "error", err)
		return
	}
	l.Info("Configuration", "config", string(configJSON))

	// 启动服务器
	server.Setupserver()
}
