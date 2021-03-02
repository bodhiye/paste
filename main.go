package main

import (
	"context"
	"net/http"
	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	"paste.org.cn/paste/middleware"
	"paste.org.cn/paste/router"
	"paste.org.cn/paste/util"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	util.LoadConfig("config")

	paste := gin.New()
	paste.Use(gin.Recovery())
	paste.Use(middleware.LogInfo)
	paste.Use(middleware.ReqID)
	router.Init(paste)

	srv := &http.Server{
		Addr:    util.GetServerHost(viper.GetString("server.host")),
		Handler: paste,
	}
	util.RunServer(srv)
	cancel()
	util.ShutdownServer(srv)
}
