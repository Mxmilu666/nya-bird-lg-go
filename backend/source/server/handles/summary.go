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
		targetServers = strings.Split(serverParam, ",")
	}

	var serverIDs []string
	if len(targetServers) > 0 {
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
				SendResponse(c, 404, "error", "Server "+targetServer+" not found")
				return
			}
		}
	} else {
		for _, server := range source.AppConfig.LG.Servers {
			serverIDs = append(serverIDs, server.ID)
		}
	}

	command := helper.GetBirdCommand("summary", "")

	responses, errors := helper.BatchRequest(c, serverIDs, "/bird", command)

	result := make(map[string]map[string]interface{})
	for i, serverId := range serverIDs {
		var server source.ServerInfo
		for _, s := range source.AppConfig.LG.Servers {
			if s.ID == serverId {
				server = s
				break
			}
		}

		serverKey := map[string]interface{}{
			"id":          server.ID,
			"displayName": server.DisplayName,
		}

		if errors[i] != nil {
			serverKey["error"] = errors[i].Error()
		} else {
			serverKey["protocols"] = parseSummaryOutput(responses[i])
		}

		result[server.ID] = serverKey
	}

	response := result

	SendResponse(c, 200, "success", response)
}
