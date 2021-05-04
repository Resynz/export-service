package exter

import (
	"encoding/base64"
	"export-service/code"
	"export-service/common"
	"export-service/config"
	"export-service/db"
	"fmt"
	model_export "github.com/Resynz/model-export"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

func Download(ctx *gin.Context) {
	alias := ctx.Param("alias")
	b, err := base64.StdEncoding.DecodeString(alias)
	if err != nil {
		log.Printf("[Download] error:%s\n", err.Error())
		common.HandleResponse(ctx, code.InvalidRequest, nil)
		return
	}
	aliasArr := strings.Split(string(b), "|")
	if len(aliasArr) != 2 {
		common.HandleResponse(ctx, code.InvalidRequest, nil)
		return
	}
	sn := aliasArr[0]
	operator, err := strconv.Atoi(aliasArr[1])
	if err != nil {
		log.Printf("[Download] error:%s\n", err.Error())
		common.HandleResponse(ctx, code.InvalidRequest, nil)
		return
	}
	var task model_export.Task
	has, err := db.ExportHandler.GetOne(&task, task.GetTableName(), "sn", sn)
	if err != nil {
		log.Printf("[Download] error:%s\n", err.Error())
		common.HandleResponse(ctx, code.BadRequest, nil)
		return
	}
	if !has || task.Status != model_export.ExportTaskStatusSuccess {
		common.HandleResponse(ctx, code.InvalidRequest, nil)
		return
	}
	var taskDetail model_export.TaskDetail
	has, err = db.ExportHandler.GetOne(&taskDetail, taskDetail.GetTableName(), "sn", sn)
	if err != nil {
		log.Printf("[Download] error:%s\n", err.Error())
		common.HandleResponse(ctx, code.BadRequest, nil)
		return
	}
	if !has || taskDetail.ResultPath == "" {
		common.HandleResponse(ctx, code.InvalidRequest, nil)
		return
	}
	filePath := fmt.Sprintf("%s/%s", config.Conf.ResultPath, taskDetail.ResultPath)
	f, err := os.Stat(filePath)
	if err != nil {
		log.Printf("[Download] error:%s\n", err.Error())
		common.HandleResponse(ctx, code.BadRequest, nil)
		return
	}

	if f.IsDir() {
		common.HandleResponse(ctx, code.BadRequest, nil)
		return
	}

	// 记录下载日志
	downloadLog := &model_export.DownloadLog{
		Id:         0,
		Sn:         sn,
		Downloader: int64(operator),
		CreateTime: time.Now().Unix(),
		Ip:         ctx.ClientIP(),
		UserAgent:  ctx.Request.UserAgent(),
	}

	if err = db.ExportHandler.Save(downloadLog, downloadLog.GetTableName()); err != nil {
		log.Printf("[Download] error:%s\n", err.Error())
		common.HandleResponse(ctx, code.BadRequest, nil)
		return
	}

	attachName := taskDetail.FileName
	fileExt := path.Ext(filePath)
	if !strings.HasSuffix(attachName, fileExt) {
		attachName = fmt.Sprintf("%s%s", attachName, fileExt)
	}

	ctx.FileAttachment(filePath, attachName)
}
