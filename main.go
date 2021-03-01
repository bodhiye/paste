package main

import (
	"context"
	"net/http"
	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	"paste.org.cn/paste/middleware"
	"paste.org.cn/paste/util"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	util.LoadConfig("config")

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.LogInfo)
	router.Use(middleware.ReqID)
	router.Any("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "paste ok!")
	})

	srv := &http.Server{
		Addr:    util.GetServerHost(viper.GetString("server.host")),
		Handler: router,
	}
	util.RunServer(srv)
	cancel()
	util.ShutdownServer(srv)
}
