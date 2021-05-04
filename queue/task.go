package queue

import (
	"export-service/db"
	"export-service/lib/executor"
	db_handler "github.com/Resynz/db-handler"
	model_export "github.com/Resynz/model-export"
	"log"
)

func PushQueueTask(sn string) {
	taskQueue <- sn
}

func initTaskQueue() {
	var task model_export.Task
	var list []*model_export.Task
	err := db.ExportHandler.List(&list, task.GetTableName(), &db_handler.Condition{
		Where:  "status in (?,?)",
		Params: []interface{}{model_export.ExportTaskStatusPending, model_export.ExportTaskStatusProcessing},
	})
	if err != nil {
		log.Printf("初始化任务队列失败！ error:%s\n", err.Error())
		return
	}
	for _, v := range list {
		taskQueue <- v.Sn
	}
}

func startTaskQueue() {
	go initTaskQueue()
	for !StopQueue {
		sn := <-taskQueue
		if sn == "" {
			continue
		}
		dealTask(sn)
	}
	StopQueueFlag <- true
}

func dealTask(sn string) {
	ex := &executor.TaskExecutor{
		Sn: sn,
	}
	ex.Exec()
}
