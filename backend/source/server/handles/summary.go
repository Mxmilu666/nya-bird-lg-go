package handles

import (
	"strings"

	"github.com/Mxmilu666/nya-bird-lg-go/source"
	"github.com/Mxmilu666/nya-bird-lg-go/source/helper"
	"github.com/gin-gonic/gin"
)

// Protocol 表示BIRD协议的结构
type Protocol struct {
	Name  string `json:"name"`
	Proto string `json:"proto"`
	Table string `json:"table"`
	State string `json:"state"`
	Since string `json:"since"`
	Info  string `json:"info"`
}

func parseSummaryOutput(output string) []Protocol {
	var protocols []Protocol
	lines := strings.Split(output, "\n")

	// 跳过标题行
	for _, line := range lines[1:] {
		if len(strings.TrimSpace(line)) == 0 {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 5 {
			continue
		}

		protocol := Protocol{
			Name:  fields[0],
			Proto: fields[1],
			Table: fields[2],
			State: fields[3],
			Since: fields[4],
		}

		// 如果有额外的Info字段
		if len(fields) > 5 {
			protocol.Info = strings.Join(fields[5:], " ")
		}

		protocols = append(protocols, protocol)
	}
	return protocols
}

func GetBirdSummary(c *gin.Context) {
	serverParam := c.Query("server")
	var targetServers []string

	if serverParam != "" {
		// 将包含+号的参数分割成多个服务器ID
		targetServers = strings.Split(serverParam, ",")
	}

	// 从配置中获取所有服务器ID
	var serverIDs []string
	if len(targetServers) > 0 {
		// 如果指定了目标服务器，添加所有匹配的服务器
		for _, targetServer := range targetServers {
			targetServer = strings.TrimSpace(targetServer)
			if targetServer == "" {
				continue
			}
			found := false
			for _, server := range source.AppConfig.LG.Servers {
				if server.ID == targetServer {
					serverIDs = append(serverIDs, server.ID)
					found = true
					break
				}
			}
			if !found {
				c.JSON(404, gin.H{
					"error": "服务器 " + targetServer + " 不存在",
				})
				return
			}
		}
	} else {
		// 如果没有指定服务器，返回所有服务器
		for _, server := range source.AppConfig.LG.Servers {
			serverIDs = append(serverIDs, server.ID)
		}
	}

	// 获取summary命令
	command := helper.GetBirdCommand("summary", "")

	// 使用batchRequest获取服务器的汇总信息
	responses, err := helper.BatchRequest(c, serverIDs, "/bird", command)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 构造响应
	result := make(map[string]map[string]interface{})
	for i, serverId := range serverIDs {
		var server source.ServerInfo
		// 查找对应的服务器配置
		for _, s := range source.AppConfig.LG.Servers {
			if s.ID == serverId {
				server = s
				break
			}
		}

		serverKey := map[string]interface{}{
			"id":          server.ID,
			"displayName": server.DisplayName,
			"protocols":   parseSummaryOutput(responses[i]),
		}
		result[server.ID] = serverKey
	}

	c.JSON(200, result)
}
