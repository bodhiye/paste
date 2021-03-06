package main

import (
	"context"
	"net/http"
	"runtime"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"paste.org.cn/paste/server/db"
	"paste.org.cn/paste/server/middleware"
	"paste.org.cn/paste/server/router"
	"paste.org.cn/paste/server/util"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	util.LoadConfig("config")

	paste := gin.New()
	paste.Use(gin.Recovery())
	paste.Use(middleware.LogInfo)
	paste.Use(middleware.ReqID)

	pasteDB, err := db.NewPaste(ctx, viper.Sub("paste.mgo"))
	if err != nil {
		log.Errorf("init paste db failed: %+v", err)
		return
	}
	router.Init(paste, pasteDB)

	srv := &http.Server{
		Addr:    util.GetServerHost(viper.GetString("server.host")),
		Handler: paste,
	}
	util.RunServer(srv)
	cancel()
	util.ShutdownServer(srv)
}
