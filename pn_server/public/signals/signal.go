package signals

import (
	"dc3/public/logs"
	"os"
	"os/signal"
	"syscall"
)

// WaitSignal 等待信号, 收到信号后执行回调函数
func WaitWith(stops ...func()) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logs.Infof("优雅退出服务...")

	// 优雅退出
	for _, stop := range stops {
		stop()
	}
}
