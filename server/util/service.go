package util

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	log "github.com/sirupsen/logrus" // 用logrus第三方开源库来替换标准库log包，logrus兼容标准库log包的所有API
	"github.com/spf13/viper"
)

// 加载配置
func LoadConfig(configName string) {
	viper.SetConfigName(configName) 
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Fatal error config file: %+v", err)
	}

	level, _ := log.ParseLevel(viper.GetString("log.level")) // 将从配置文件中获取到的日志从字符串转为log.level类型
	log.SetLevel(level)
}

// GetServerHost 尝试替换原本从配置文件中设置的 服务器地址的端口号
func GetServerHost(sh string) string {
	// 获取环境变量 PORT_HTTP 的值，通常用于指定服务器监听的端口号
	p := os.Getenv("PORT_HTTP")
	// 使用 net.SplitHostPort 函数将传入的地址 sh 分割为主机名（h）和端口号（_）
	h, _, err := net.SplitHostPort(sh)
	if err != nil || len(p) == 0 {
		return sh
	}
	return fmt.Sprintf("%s:%s", h, p)
}

// 启动服务器
func RunServer(srv *http.Server) {
	log.Printf("Starting server at %s", srv.Addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %+v", err)
	}
	log.Println("Server stopped")
}

// ShutdownServer 优雅地关闭服务器，等待当前活动连接完成
func ShutdownServer(srv *http.Server) {
	// 创建一个 5 秒超时的上下文，如果5秒内所有连接没有关闭，强制关闭服务器
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// Shutdown 尝试关闭服务器
	// 优雅地停止服务器服务
	// 确保正在处理的请求尽可能地完成
	// 提供一个超时机制，防止无限期等待
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown: %+v", err)
	}
	log.Println("Server exiting")
}