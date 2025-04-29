package server

import (
	"embed"
	"net/http"

	"github.com/Mxmilu666/nya-bird-lg-go/source/server/handles"
	"github.com/Mxmilu666/nya-bird-lg-go/source/server/middleware"

	"github.com/gin-gonic/gin"
)

func initRouter(r *gin.Engine, frontendFS embed.FS) *gin.Engine {
	r.Use(middleware.Serve("/", middleware.EmbedFolder(frontendFS, "frontend/dist")))
	r.NoRoute(func(c *gin.Context) {
		data, err := frontendFS.ReadFile("frontend/dist/index.html")
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.Data(http.StatusOK, "text/html; charset=utf-8", data)
	})
	// API路由
	api := r.Group("/api")
	{
		// bird 相关路由
		bird := api.Group("/bird")
		{
			bird.GET("/summary", handles.GetBirdSummary)
			bird.GET("/detail", handles.GetBirdDetail)
		}

		// 服务器信息路由
		api.GET("/servers", handles.GetServerList)

		// traceroute 路由
		api.GET("/traceroute", handles.GetTraceroute)
	}

	return r
}
