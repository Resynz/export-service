package server

import (
	"export-service/controller/exter"
	"github.com/gin-gonic/gin"
)

func RegisterExternalRoute(route *gin.RouterGroup) {
	downloadGroup := route.Group("/download")
	downloadGroup.GET("/:alias", exter.Download)
}
