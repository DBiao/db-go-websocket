package app

import (
	"context"
	"db-go-websocket/internal/dgrpc/dclient"
	"db-go-websocket/internal/dgrpc/dserver"
	"db-go-websocket/internal/global"
	"db-go-websocket/internal/router"
	"db-go-websocket/pkg/etcd"
	"db-go-websocket/pkg/kafka"
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

	// 初始化grpc client
	if global.GRPCClient, err = dclient.InitGrpcClient(); err != nil {
		log.Fatal(err)
	}

	// 初始化grpc服务
	if err = dserver.InitGrpcServer(); err != nil {
		log.Fatal(err)
	}

	// 初始化etcd
	if err = etcd.InitEtcd(); err != nil {
		log.Fatal(err)
	}

	// 初始化kafka
	if global.KAFKA, err = kafka.InitKafka(); err != nil {
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
			if global.GRPCSServer != nil {
				global.GRPCSServer.Stop()
			}
		},

		// 关闭kafka
		func() {
			if global.KAFKA != nil {
				global.KAFKA.Close()
			}
		},

		// 关闭grpc client
		func() {
			global.GRPCClient.Close()
		})
}
