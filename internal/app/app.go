package app

import (
	"context"
	"db-go-websocket/internal/dgrpc/dserver"
	"db-go-websocket/internal/global"
	"db-go-websocket/internal/router"
	"db-go-websocket/pkg/etcd"
	"db-go-websocket/pkg/logger"
	"db-go-websocket/pkg/shutdown"
	"db-go-websocket/pkg/viper"
	"log"
	"time"

	"go.uber.org/zap"
)

// Start 初始化服务
func Start() {
	var err error

	// 初始化配置
	if global.VIPER, err = viper.InitViper(); err != nil {
		log.Fatal(err)
	}

	// 初始化日志
	global.LOG = logger.InitZap()

	// 初始化grpc服务
	if err = dserver.InitGrpcServer(); err != nil {
		log.Fatal(err)
	}

	// 初始化etcd
	if err = etcd.InitEtcd(); err != nil {
		log.Fatal(err)
	}

	// 初始化Gin
	router.InitRouter()
}

// Close 优雅关闭
func Close() {
	shutdown.NewHook().Close(
		// 关闭http server
		func() {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
			defer cancel()

			if err := global.SERVER.Shutdown(ctx); err != nil {
				global.LOG.Error("dserver shutdown err", zap.Error(err))
			}
		},

		// 关闭grpc server
		func() {
			if global.GRPCSERVER != nil {
				global.GRPCSERVER.Stop()
			}
		})
}
