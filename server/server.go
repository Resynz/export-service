/**
 * @Author: Resynz
 * @Date: 2021/4/22 11:48
 */
package server

import (
	"export-service/config"
	"export-service/controller"
	"fmt"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"log"
)

func StartServer() {
	gin.SetMode(config.Conf.Mode)
	app := gin.New()
	app.MaxMultipartMemory = 8 << 20 // 8mb
	app.Use(gzip.Gzip(gzip.DefaultCompression))
	// 添加recovery中间件
	app.Use(gin.Recovery())
	app.GET("/ping", controller.Ping)

	RegisterInternalRoute(app.Group("/internal"))

	RegisterExternalRoute(app.Group("/external"))

	if err := app.Run(fmt.Sprintf(":%d", config.Conf.AppPort)); err != nil {
		log.Fatalf("start server failed! error:%s\n", err.Error())
	}
}
