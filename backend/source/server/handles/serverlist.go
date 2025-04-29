package handles

import (
	"github.com/Mxmilu666/nya-bird-lg-go/source"
	"github.com/gin-gonic/gin"
)

// ServerListResponse 表示服务器列表响应的结构
type ServerListResponse struct {
	ID          string `json:"id"`           // 服务器ID
	DisplayName string `json:"display_name"` // 服务器显示名称
}

// GetServerList 返回所有配置的服务器列表
func GetServerList(c *gin.Context) {
	var serverList []ServerListResponse

	// 从配置中获取所有服务器并转换为响应格式
	for _, server := range source.AppConfig.LG.Servers {
		serverList = append(serverList, ServerListResponse{
			ID:          server.ID,
			DisplayName: server.DisplayName,
		})
	}

	SendResponse(c, 200, "success", serverList)
}
