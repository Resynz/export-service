package inter

import (
	"export-service/action"
	"export-service/code"
	"export-service/common"
	"export-service/db"
	"export-service/queue"
	"export-service/tools"
	"fmt"
	model_export "github.com/Resynz/model-export"
	"github.com/gin-gonic/gin"
	"time"
)

type AddTaskParam struct {
	FileName  string `json:"file_name" binding:"required"`
	Creator   int64  `json:"creator"`
	Action    string `json:"action" binding:"required"`
	Condition string `json:"condition" binding:"required"`
	NotifyUrl string `json:"notify_url"`
}

func AddTask(ctx *gin.Context) {
	var form AddTaskParam
	if err := ctx.ShouldBind(&form); err != nil {
		common.HandleResponse(ctx, code.InvalidParams, nil, err.Error())
		return
	}
	_, ok := action.ActionMap[form.Action]
	if !ok {
		common.HandleResponse(ctx, code.InvalidRequest, nil, fmt.Sprintf("action:%s not found", form.Action))
		return
	}

	sn, err := tools.GenerateSn()
	if err != nil {
		common.HandleResponse(ctx, code.BadRequest, nil, err.Error())
		return
	}

	session := db.ExportHandler.DB.NewSession()
	if err = session.Begin(); err != nil {
		common.HandleResponse(ctx, code.BadRequest, nil, err.Error())
		return
	}
	defer session.Close()

	task := &model_export.Task{
		Id:          0,
		Sn:          sn,
		Status:      model_export.ExportTaskStatusPending,
		Creator:     form.Creator,
		CreateTime:  time.Now().Unix(),
		ExecuteTime: 0,
		FinishTime:  0,
		CostTime:    0,
		UpdateTime:  time.Now().Unix(),
	}
	if _, err = session.Insert(task); err != nil {
		common.HandleResponse(ctx, code.BadRequest, nil, err.Error())
		return
	}

	taskDetail := &model_export.TaskDetail{
		Id:         0,
		Sn:         sn,
		Action:     form.Action,
		Condition:  form.Condition,
		ResultPath: "",
		NotifyUrl:  form.NotifyUrl,
		FileName:   form.FileName,
	}
	if _, err = session.Insert(taskDetail); err != nil {
		common.HandleResponse(ctx, code.BadRequest, nil, err.Error())
		return
	}

	taskLog := &model_export.TaskLog{
		Id:         0,
		Sn:         sn,
		Type:       model_export.ExportTaskLogTypeCreate,
		Operator:   form.Creator,
		Remark:     "",
		CreateTime: time.Now().Unix(),
	}

	if _, err = session.Insert(taskLog); err != nil {
		common.HandleResponse(ctx, code.BadRequest, nil, err.Error())
		return
	}

	if err = session.Commit(); err != nil {
		common.HandleResponse(ctx, code.BadRequest, nil, err.Error())
		return
	}

	// push to queue
	go queue.PushQueueTask(sn)

	data := map[string]string{
		"sn": sn,
	}

	common.HandleResponse(ctx, code.SuccessCode, data)
}
