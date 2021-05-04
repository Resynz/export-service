package executor

import (
	"encoding/json"
	action2 "export-service/action"
	"export-service/db"
	"export-service/lib/notify"
	"fmt"
	db_handler "github.com/Resynz/db-handler"
	model_export "github.com/Resynz/model-export"
	"log"
	"time"
)

type TaskExecutor struct {
	Sn string
}

func (s *TaskExecutor) log(txt string) {
	log.Printf("[TaskExecutor] task(sn:%s) %s\n", s.Sn, txt)
}

func (s *TaskExecutor) Exec() {
	var err error
	var has bool
	notifyParam := notify.NotifyParam{
		NotifyUrl: "",
		Result:    notify.NotifyResultSuccess,
		Message:   "",
	}
	defer func() {
		var task model_export.Task
		has, e := db.ExportHandler.GetOne(&task, task.GetTableName(), "sn", s.Sn)
		if e != nil {
			s.log(e.Error())
			return
		}
		if has {
			status := model_export.ExportTaskStatusSuccess

			taskLog := &model_export.TaskLog{
				Id:         0,
				Sn:         s.Sn,
				Type:       model_export.ExportTaskLogTypeExecuteSuccess,
				Operator:   model_export.System,
				Remark:     "",
				CreateTime: time.Now().Unix(),
			}

			if err != nil {
				s.log(err.Error())
				status = model_export.ExportTaskStatusFailed
				taskLog.Type = model_export.ExportTaskLogTypeExecuteFailed
				taskLog.Remark = err.Error()
				notifyParam.Result = notify.NotifyResultFailed
				notifyParam.Message = err.Error()
			}
			task.Status = status
			task.UpdateTime = time.Now().Unix()
			task.CostTime = time.Now().Unix() - task.ExecuteTime
			task.FinishTime = time.Now().Unix()
			_ = db.ExportHandler.Save(&task, task.GetTableName())
			_ = db.ExportHandler.Save(taskLog, taskLog.GetTableName())
			// todo send notify
		}
	}()

	s.log("executing start ...")
	var task model_export.Task
	has, err = db.ExportHandler.GetOne(&task, task.GetTableName(), "sn", s.Sn)
	if err != nil {
		return
	}
	if !has {
		err = fmt.Errorf("task not found")
		return
	}
	if task.Status == model_export.ExportTaskStatusSuccess {
		err = fmt.Errorf("invalid status:%d", task.Status)
		return
	}

	task.Status = model_export.ExportTaskStatusProcessing
	task.ExecuteTime = time.Now().Unix()
	task.CostTime = 0
	task.FinishTime = 0
	if err = db.ExportHandler.Save(&task, task.GetTableName()); err != nil {
		return
	}

	taskLog := &model_export.TaskLog{
		Id:         0,
		Sn:         s.Sn,
		Type:       model_export.ExportTaskLogTypeExecuteStart,
		Operator:   model_export.System,
		Remark:     "",
		CreateTime: time.Now().Unix(),
	}
	if err = db.ExportHandler.Save(taskLog, taskLog.GetTableName()); err != nil {
		return
	}

	var taskDetail model_export.TaskDetail
	has, err = db.ExportHandler.GetOne(&taskDetail, taskDetail.GetTableName(), "sn", s.Sn)
	if err != nil {
		return
	}

	if !has {
		err = fmt.Errorf("task_detail not found")
		return
	}
	notifyParam.NotifyUrl = taskDetail.NotifyUrl

	// todo 1.找到action
	action, ok := action2.ActionMap[taskDetail.Action]
	if !ok {
		err = fmt.Errorf("action:%s not found", taskDetail.Action)
		return
	}
	var condition db_handler.Condition
	if err = json.Unmarshal([]byte(taskDetail.Condition), &condition); err != nil {
		return
	}
	res, err := action(&condition)
	if err != nil {
		return
	}
	res.Creator = task.Creator
	if err = res.Exec(); err != nil {
		return
	}
	taskDetail.ResultPath = res.ResultPath
	if err = db.ExportHandler.Save(&taskDetail, taskDetail.GetTableName()); err != nil {
		return
	}
}
