/**
 * @Author: Resynz
 * @Date: 2021/4/22 11:19
 */
package main

import (
	"export-service/action"
	"export-service/db"
	"export-service/handler"
	"export-service/queue"
	"export-service/server"
	"log"
)

func main() {
	log.Println("starting service exit listener ...")
	handler.SetSignalHandler()
	log.Println("initializing db handler ...")
	if err := db.InitDBHandler(); err != nil {
		log.Fatalf("init db handler failed! error:%s\n", err.Error())
	}
	log.Println("initializing actions ...")
	action.InitActions()
	log.Println("starting queue ...")
	go queue.StartQueue()
	log.Println("\033[42;30m DONE \033[0m[ExportService] Start Success!")
	server.StartServer()
}
