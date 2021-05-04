package inter

import (
	"encoding/base64"
	"export-service/code"
	"export-service/common"
	"export-service/config"
	"export-service/db"
	"fmt"
	model_export "github.com/Resynz/model-export"
	"github.com/gin-gonic/gin"
)

func GenerateDownloadUrl(ctx *gin.Context) {
	sn := ctx.Param("sn")
	var task model_export.Task
	has, err := db.ExportHandler.GetOne(&task, task.GetTableName(), "sn", sn)
	if err != nil {
		common.HandleResponse(ctx, code.BadRequest, nil, err.Error())
		return
	}
	if !has || task.Status != model_export.ExportTaskStatusSuccess {
		common.HandleResponse(ctx, code.InvalidRequest, nil)
		return
	}

	type formValidate struct {
		Operator int64 `form:"operator" binding:"required"`
	}
	var form formValidate
	if err = ctx.ShouldBind(&form); err != nil {
		common.HandleResponse(ctx, code.InvalidParams, nil, err.Error())
		return
	}

	alias := fmt.Sprintf("%s|%d", task.Sn, form.Operator)

	downloadUrl := fmt.Sprintf("%s/external/download/%s", config.Conf.Host, base64.StdEncoding.EncodeToString([]byte(alias)))

	data := map[string]string{
		"download_url": downloadUrl,
	}
	common.HandleResponse(ctx, code.SuccessCode, data)
}
