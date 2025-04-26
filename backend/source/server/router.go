package server

import (
	"github.com/Mxmilu666/nya-bird-lg-go/source/server/handles"

	"github.com/gin-gonic/gin"
)

func initRouter(r *gin.Engine) *gin.Engine {
	// API路由
	api := r.Group("/api")
	{
		// bird 相关路由
		bird := api.Group("/bird")
		{
			bird.GET("/summary", handles.GetBirdSummary)
			bird.GET("/detail", handles.GetBirdDetail)
		}
	}
	return r
}
