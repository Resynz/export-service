/**
 * @Author: Resynz
 * @Date: 2021/4/22 11:50
 */
package handler

import (
	"export-service/queue"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func SetSignalHandler() {
	sign := make(chan os.Signal, 1)
	signal.Notify(sign, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGSEGV)
	go func() {
		defer func() {
			log.Println("程序退出.")
			os.Exit(0)
		}()
		s := <-sign
		log.Printf("接收到退出信号:%v\n", s)
		log.Println("正在停止任务队列 ...")
		queue.StopQueue = true
		queue.PushQueueTask("")
		<-queue.StopQueueFlag
		log.Println("停止任务队列完毕!")
	}()
}
