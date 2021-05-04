package queue

import "export-service/config"

var (
	taskQueue     chan string
	StopQueue     = false
	StopQueueFlag = make(chan bool, 1)
)

func StartQueue() {
	taskQueue = make(chan string, config.Conf.QueueSize)
	StopQueue = false
	go startTaskQueue()
}
