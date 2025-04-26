package source

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	logger "github.com/Mxmilu666/nya-bird-lg-go/source/logger"
	"github.com/joho/godotenv"
)

// ServerConfig 结构体定义服务器配置项
type ServerConfig struct {
	Host string
	Port int
}

type LookingGlassConfig struct {
	Servers   []ServerInfo
	Domain    string
	ProxyPort int
	SSL       bool
	Timeout   int
}

// Config 结构体定义配置项
type Config struct {
	Server ServerConfig
	LG     LookingGlassConfig
}

type ServerInfo struct {
	ID          string
	DisplayName string
}

// 全局变量保存配置
var AppConfig *Config

// getEnvAs 通用的环境变量获取函数
func getEnvAs[T string | int | bool](key string, defaultValue T) T {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}

	var result T
	var err error

	switch any(defaultValue).(type) {
	case string:
		return any(valueStr).(T)
	case int:
		val, e := strconv.Atoi(valueStr)
		result, err = any(val).(T), e
	case bool:
		val, e := strconv.ParseBool(valueStr)
		result, err = any(val).(T), e
	}

	if err != nil {
		logger.Warn(fmt.Sprintf("Invalid format for environment variable %s, using default value: %v", key, defaultValue))
		return defaultValue
	}

	return result
}

// LoadConfig 加载配置
func LoadConfig() error {
	// 尝试加载.env文件
	loadEnvFile()

	serverHost := getEnvAs("BIRDLG_HOST", "0.0.0.0")
	serverPort := getEnvAs("BIRDLG_LISTEN", 5000)

	lg_servers := getEnvAs("BIRDLG_SERVERS", "")
	lg_servers_list := ParseServerList(lg_servers)
	lg_domain := getEnvAs("BIRDLG_DOMAIN", "")
	lg_proxy_port := getEnvAs("BIRDLG_PROXY_PORT", 8000)
	lg_ssl := getEnvAs("BIRDLG_SSL", false)
	lg_timeout := getEnvAs("BIRDLG_TIMEOUT", 10)

	AppConfig = &Config{
		Server: ServerConfig{
			Host: serverHost,
			Port: serverPort,
		},
		LG: LookingGlassConfig{
			Servers:   lg_servers_list,
			Domain:    lg_domain,
			ProxyPort: lg_proxy_port,
			SSL:       lg_ssl,
			Timeout:   lg_timeout,
		},
	}
	return nil
}

// loadEnvFile 尝试从当前目录和上级目录加载.env文件
func loadEnvFile() {
	// 尝试直接加载.env
	err := godotenv.Load()
	if err == nil {
		return
	}
}

// ParseServerList 解析服务器列表字符串
func ParseServerList(servers string) []ServerInfo {
	if servers == "" {
		return nil
	}

	// 按逗号分割服务器列表
	serverList := strings.Split(servers, ",")
	result := make([]ServerInfo, 0, len(serverList))

	for _, server := range serverList {
		// 去除空白字符
		server = strings.TrimSpace(server)
		if server == "" {
			continue
		}

		// 检查是否包含显示名称（使用<>标记）
		if strings.Contains(server, "<") && strings.Contains(server, ">") {
			// 提取显示名称和ID
			parts := strings.Split(server, "<")
			displayName := parts[0]
			id := strings.TrimSuffix(parts[1], ">")
			result = append(result, ServerInfo{
				ID:          id,
				DisplayName: displayName,
			})
		} else {
			// 简单格式，ID即为显示名称
			result = append(result, ServerInfo{
				ID:          server,
				DisplayName: server,
			})
		}
	}

	return result
}
