package server

import (
	"export-service/controller/inter"
	"github.com/gin-gonic/gin"
)

func RegisterInternalRoute(route *gin.RouterGroup) {
	taskGroup := route.Group("/task")
	taskGroup.POST("/", inter.AddTask)
	taskGroup.GET("/download-url/:sn", inter.GenerateDownloadUrl)
}
