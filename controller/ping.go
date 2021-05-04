package controller

import (
	"export-service/code"
	"export-service/common"
	"github.com/gin-gonic/gin"
)

func Ping(ctx *gin.Context) {
	common.HandleResponse(ctx, code.SuccessCode, nil)
}
