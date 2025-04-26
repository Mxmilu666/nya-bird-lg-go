package handles

import (
	"regexp"
	"strings"

	"github.com/Mxmilu666/nya-bird-lg-go/source"
	"github.com/Mxmilu666/nya-bird-lg-go/source/helper"
	"github.com/gin-gonic/gin"
)

// DetailInfo 表示协议详细信息的结构
type DetailInfo struct {
	State        string        `json:"state,omitempty"`
	NeighborAddr string        `json:"neighbor_address,omitempty"`
	NeighborAS   string        `json:"neighbor_as,omitempty"`
	LocalAS      string        `json:"local_as,omitempty"`
	NeighborID   string        `json:"neighbor_id,omitempty"`
	Channels     []ChannelInfo `json:"channels,omitempty"`
}

// ChannelInfo 表示协议中通道的详细信息
type ChannelInfo struct {
	Name       string         `json:"name"`         // 通道名称 (ipv4, ipv6等)
	State      string         `json:"state"`        // 通道状态
	RouteStats RouteStatsInfo `json:"route_stats"`  // 路由变更统计
	BGPNextHop string         `json:"bgp_next_hop"` // BGP下一跳
}

// RouteStatsInfo 表示路由变更统计信息
type RouteStatsInfo struct {
	ImportUpdates   map[string]string `json:"import_updates"`   // 导入更新统计
	ImportWithdraws map[string]string `json:"import_withdraws"` // 导入撤回统计
	ExportUpdates   map[string]string `json:"export_updates"`   // 导出更新统计
	ExportWithdraws map[string]string `json:"export_withdraws"` // 导出撤回统计
}

// parseDetailOutput 解析detail命令的输出
func parseDetailOutput(output string) DetailInfo {
	var detail DetailInfo
	lines := strings.Split(output, "\n")

	var currentChannel *ChannelInfo
	var inRouteStats bool

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}

		// 新通道开始
		if m := regexp.MustCompile(`^\s*Channel\s+(\S+)`).FindStringSubmatch(line); len(m) > 0 {
			// 重置路由统计标志
			inRouteStats = false
			detail.Channels = append(detail.Channels, ChannelInfo{Name: m[1]})
			currentChannel = &detail.Channels[len(detail.Channels)-1]
			continue
		}

		// 基本 value 处理
		if kv := regexp.MustCompile(`^\s*(\S[^:]+):\s+(.+)`).FindStringSubmatch(line); len(kv) > 0 {
			key, val := strings.TrimSpace(kv[1]), strings.TrimSpace(kv[2])
			switch key {
			case "BGP state":
				detail.State = val
			case "Neighbor address":
				detail.NeighborAddr = val
			case "Neighbor AS":
				detail.NeighborAS = val
			case "Local AS":
				detail.LocalAS = val
			case "Neighbor ID":
				detail.NeighborID = val
			}

			if currentChannel != nil {
				switch key {
				case "State":
					currentChannel.State = val
				case "BGP Next hop":
					currentChannel.BGPNextHop = val
				case "Route change stats":
					// 进入统计区
					inRouteStats = true
					currentChannel.RouteStats = RouteStatsInfo{
						ImportUpdates:   make(map[string]string),
						ImportWithdraws: make(map[string]string),
						ExportUpdates:   make(map[string]string),
						ExportWithdraws: make(map[string]string),
					}
				}
			}
		}

		// 独立解析 Import updates Export withdraws
		if inRouteStats && currentChannel != nil {
			if stats := regexp.MustCompile(`^\s+(\S[^:]+):\s+(.+)`).FindStringSubmatch(line); len(stats) > 0 {
				statsType := strings.TrimSpace(stats[1])
				parts := strings.Fields(stats[2])
				names := []string{"received", "rejected", "filtered", "ignored", "accepted"}
				m := make(map[string]string)
				for i, n := range names {
					if i < len(parts) {
						m[n] = parts[i]
					} else {
						m[n] = ""
					}
				}
				switch statsType {
				case "Import updates":
					currentChannel.RouteStats.ImportUpdates = m
				case "Import withdraws":
					currentChannel.RouteStats.ImportWithdraws = m
				case "Export updates":
					currentChannel.RouteStats.ExportUpdates = m
				case "Export withdraws":
					currentChannel.RouteStats.ExportWithdraws = m
				}
			}
		}
	}

	return detail
}

// GetBirdDetail 获取指定协议的详细信息
func GetBirdDetail(c *gin.Context) {
	// 获取查询参数
	serverParam := c.Query("server")
	protocolParam := c.Query("protocol")

	// 协议参数是必需的
	if protocolParam == "" {
		SendResponse(c, 400, "error", "Missing required parameter: protocol")
		return
	}

	var targetServers []string
	if serverParam != "" {
		// 将包含,号的参数分割成多个服务器ID
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
				SendResponse(c, 404, "error", "Server "+targetServer+" not found")
				return
			}
		}
	} else {
		// 如果没有指定服务器，返回所有服务器
		for _, server := range source.AppConfig.LG.Servers {
			serverIDs = append(serverIDs, server.ID)
		}
	}

	// 获取detail命令
	command := helper.GetBirdCommand("detail", protocolParam)

	// 使用batchRequest获取服务器的详细信息
	responses, err := helper.BatchRequest(c, serverIDs, "/bird", command)
	if err != nil {
		SendResponse(c, 500, "error", err.Error())
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

		// 解析详细输出
		parsedDetail := parseDetailOutput(responses[i])

		serverKey := map[string]interface{}{
			"id":          server.ID,
			"displayName": server.DisplayName,
			"detail":      parsedDetail,
			"rawOutput":   responses[i], // 保留原始输出
		}
		result[server.ID] = serverKey
	}

	SendResponse(c, 200, "success", result)
}
