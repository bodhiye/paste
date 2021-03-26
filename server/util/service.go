package util

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func LoadConfig(configName string) {
	viper.SetConfigName(configName)
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Fatal error config file: %+v", err)
	}

	level, _ := log.ParseLevel(viper.GetString("log.level"))
	log.SetLevel(level)
}

func GetServerHost(sh string) string {
	p := os.Getenv("PORT_HTTP")
	h, _, err := net.SplitHostPort(sh)
	if err != nil || len(p) == 0 {
		return sh
	}
	return fmt.Sprintf("%s:%s", h, p)
}

func RunServer(srv *http.Server) {
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %+v", err)
		}
	}()

	quit := make(chan os.Signal, 5)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")
}

func ShutdownServer(srv *http.Server) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown: %+v", err)
	}
	log.Println("Server exiting")
}
