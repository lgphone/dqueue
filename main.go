package main

import (
	"dqueue/api"
	"dqueue/deley_queue"
	"dqueue/setting"
	"dqueue/store"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

var confPath string

func Setup() {
	setting.Setup(confPath)
	store.RedisStoreSetup()
}

func Run() {
	// 接收参数
	flag.StringVar(&confPath, "c", "./app.ini", "配置文件路径.")
	flag.Parse()
	// 初始化
	Setup()
	// api http服务
	address := fmt.Sprintf(":%v", setting.AppSetting.Port)
	s := &http.Server{
		Addr:           address,
		Handler:        nil,
		ReadTimeout:    time.Second * 10,
		WriteTimeout:   time.Second * 10,
		MaxHeaderBytes: 1 << 20,
	}

	http.HandleFunc("/push", api.PushJob)
	http.HandleFunc("/delete", api.DeleteJob)
	http.HandleFunc("/get", api.GetJob)
	log.Println("INFO: start http api server ...")
	log.Printf("INFO: listen %s\n", address)
	// 开始协程轮询过期job并放入ready队列中
	log.Println("INFO: start queue server ...")
	go deley_queue.LoopPushExpiredJobReadyQueue(setting.AppSetting.PushToReadyQueueSize)
	// 启动http 服务
	go func() {
		if err := s.ListenAndServe(); err != nil {
			log.Fatalln(err)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("exit ...")
}
