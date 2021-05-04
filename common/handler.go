/**
 * @Author: Resynz
 * @Date: 2021/4/22 11:26
 */
package common

import (
	"export-service/code"
	"github.com/gin-gonic/gin"
	"net/http"
)

func HandleResponse(ctx *gin.Context, c code.ResponseCode, d interface{}, msg ...string) {
	m := code.GetCodeMsg(c)
	if len(msg) > 0 {
		m = msg[0]
	}
	data := map[string]interface{}{
		"code":    c,
		"message": m,
	}
	if d != nil {
		data["data"] = d
	}
	ctx.JSON(http.StatusOK, data)
}
