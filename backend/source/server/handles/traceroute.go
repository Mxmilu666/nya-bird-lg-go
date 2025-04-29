package handles

import (
	"regexp"
	"strings"

	"github.com/Mxmilu666/nya-bird-lg-go/source"
	"github.com/Mxmilu666/nya-bird-lg-go/source/helper"
	"github.com/gin-gonic/gin"
)

// TracerouteHop 表示traceroute的每一跳信息
type TracerouteHop struct {
	HopNum  int    `json:"hop_num"` // 跳数
	Address string `json:"address"` // IP地址
	RTT     string `json:"rtt"`     // 响应时间
	Host    string `json:"host"`    // 主机名
}

// parseTracerouteOutput 解析traceroute命令的输出并返回
func parseTracerouteOutput(output string) ([]TracerouteHop, string, string) {
	var hops []TracerouteHop
	var description string
	var errorMsg string
	lines := strings.Split(output, "\n")

	// 检查是否存在错误信息
	if strings.Contains(output, "Error executing traceroute") ||
		strings.Contains(output, "Name does not resolve") ||
		strings.Contains(output, "Cannot handle") {
		// 提取错误信息
		for _, line := range lines {
			if strings.Contains(line, "Name does not resolve") ||
				strings.Contains(line, "Cannot handle") {
				errorMsg = strings.TrimSpace(line)
				break
			}
		}
		if errorMsg == "" && strings.Contains(output, "Error executing traceroute") {
			errorMsg = "Error executing traceroute"
		}
		// 返回完整的错误信息
		return hops, description, output
	}

	// 寻找 hops not responding
	for _, line := range lines {
		if strings.Contains(line, "hops not responding") {
			description = strings.TrimSpace(line)
			break
		}
	}

	// 跳过标题行
	startIndex := 0
	for i, line := range lines {
		if strings.HasPrefix(line, "traceroute to") {
			startIndex = i + 1
			break
		}
	}

	// 正则表达式匹配标准traceroute输出格式
	hopRegex := regexp.MustCompile(`^\s*(\d+)\s+([^\s(]+)(?:\s+\(([^)]+)\))?\s+(\d+\.\d+\s*ms)`)

	// 匹配超时的行
	timeoutRegex := regexp.MustCompile(`^\s*(\d+)(?:\s+\*)+`)

	lastHopNum := 0

	for i := startIndex; i < len(lines); i++ {
		line := lines[i]

		// 跳过空行和总结行
		if len(strings.TrimSpace(line)) == 0 ||
			strings.Contains(line, "not responding") ||
			strings.Contains(line, "hops max") {
			continue
		}

		var hopNum int
		if matches := hopRegex.FindStringSubmatch(line); len(matches) > 0 {
			hopNum = 0
			for _, c := range matches[1] {
				hopNum = hopNum*10 + int(c-'0')
			}

			// 检查是否有跳数缺失
			if lastHopNum > 0 && hopNum > lastHopNum+1 {
				// 填充缺失的跳数
				for j := lastHopNum + 1; j < hopNum; j++ {
					hops = append(hops, TracerouteHop{
						HopNum:  j,
						Address: "*",
						RTT:     "timeout",
					})
				}
			}

			lastHopNum = hopNum

			hop := TracerouteHop{
				HopNum: hopNum,
			}

			// 主机名和IP地址
			hostname := matches[2]
			if hostname == "*" {
				hop.Address = "*"
				hop.RTT = "timeout"
			} else {
				hop.Host = hostname

				// IP地址
				if matches[3] != "" {
					hop.Address = matches[3]
				} else {
					hop.Address = hostname
					hop.Host = ""
				}

				// 响应时间
				if matches[4] != "" {
					hop.RTT = matches[4]
				}
			}

			hops = append(hops, hop)
		} else if timeoutMatches := timeoutRegex.FindStringSubmatch(line); len(timeoutMatches) > 0 {
			// 处理超时的情况
			hopNum = 0
			for _, c := range timeoutMatches[1] {
				hopNum = hopNum*10 + int(c-'0')
			}

			// 检查是否有跳数缺失
			if lastHopNum > 0 && hopNum > lastHopNum+1 {
				// 填充缺失的跳数
				for j := lastHopNum + 1; j < hopNum; j++ {
					hops = append(hops, TracerouteHop{
						HopNum:  j,
						Address: "*",
						RTT:     "timeout",
					})
				}
			}

			lastHopNum = hopNum

			hop := TracerouteHop{
				HopNum:  hopNum,
				Address: "*",
				RTT:     "timeout",
			}

			hops = append(hops, hop)
		}
	}

	return hops, description, errorMsg
}

// GetTraceroute 执行traceroute命令并返回结果
func GetTraceroute(c *gin.Context) {
	// 获取目标IP或域名
	targetParam := c.Query("target")
	if targetParam == "" {
		SendResponse(c, 400, "error", "Missing required parameter: target")
		return
	}

	// 获取要查询的服务器ID列表
	serverParam := c.Query("server")
	var serverIDs []string

	if serverParam != "" {
		// 指定了服务器，验证每个服务器ID的有效性
		for _, serverID := range strings.Split(serverParam, ",") {
			serverID = strings.TrimSpace(serverID)
			if serverID == "" {
				continue
			}

			found := false
			for _, server := range source.AppConfig.LG.Servers {
				if server.ID == serverID {
					serverIDs = append(serverIDs, serverID)
					found = true
					break
				}
			}

			if !found {
				SendResponse(c, 404, "error", "Server "+serverID+" not found")
				return
			}
		}
	} else {
		// 未指定服务器，使用所有可用服务器
		for _, server := range source.AppConfig.LG.Servers {
			serverIDs = append(serverIDs, server.ID)
		}
	}

	// 向所有节点发送traceroute请求
	responses, errors := helper.BatchRequest(c, serverIDs, "/traceroute", targetParam)

	// 处理结果
	result := make(map[string]map[string]interface{})
	for i, serverId := range serverIDs {
		var server source.ServerInfo
		for _, s := range source.AppConfig.LG.Servers {
			if s.ID == serverId {
				server = s
				break
			}
		}

		serverResult := map[string]interface{}{
			"id":          server.ID,
			"displayName": server.DisplayName,
		}

		if errors[i] != nil {
			serverResult["error"] = errors[i].Error()
		} else {
			hops, description, errorMsg := parseTracerouteOutput(responses[i])
			serverResult["hops"] = hops

			if description != "" {
				serverResult["description"] = description
			}

			if errorMsg != "" {
				serverResult["errorMsg"] = errorMsg
			}

			serverResult["rawOutput"] = responses[i]
		}

		result[server.ID] = serverResult
	}

	SendResponse(c, 200, "success", result)
}
