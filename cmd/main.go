package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"power-supply-sys/internal/app"
	"syscall"
	"time"
)

func main() {
	// 创建应用实例
	application, err := app.New()
	if err != nil {
		log.Fatalf("创建应用实例失败: %v", err)
	}

	// 启动服务器(在独立的 goroutine 中)
	go func() {
		if err := application.Run(); err != nil {
			log.Fatalf("服务器启动失败: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("收到关闭信号...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := application.Shutdown(ctx); err != nil {
		log.Fatalf("强制关闭应用: %v", err)
	}

	log.Println("应用已安全退出")
}
