package dserver

import (
	"db-go-websocket/internal/global"
	"db-go-websocket/internal/proto"
	"fmt"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/shimingyah/pool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
	"net"
	"time"
)

func InitGrpcServer() error {
	// 如果是集群，则启用RPC进行通讯
	if global.CONFIG.System.IsCluster {
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", global.CONFIG.Grpc.Port))
		if err != nil {
			return err
		}

		// 从输入证书文件和密钥文件为服务端构造TLS凭证
		creds, err := credentials.NewServerTLSFromFile("../pkg/tls/server.pem", "../pkg/tls/server.key")
		if err != nil {
			panic(err)
		}

		global.GRPCSServer = grpc.NewServer(
			// 配置
			grpc.InitialWindowSize(pool.InitialWindowSize),
			grpc.InitialConnWindowSize(pool.InitialConnWindowSize),
			grpc.MaxSendMsgSize(pool.MaxSendMsgSize),
			grpc.MaxRecvMsgSize(pool.MaxRecvMsgSize),

			// 健康检查
			grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
				MinTime:             5 * time.Second, // 如果客户端两次 ping 的间隔小于此值，则关闭连接
				PermitWithoutStream: true,            // 即使没有 active stream, 也允许 ping
			}),
			grpc.KeepaliveParams(keepalive.ServerParameters{
				MaxConnectionIdle:     15 * time.Second,      // 如果一个 client 空闲超过该值, 发送一个 GOAWAY, 为了防止同一时间发送大量 GOAWAY, 会在此时间间隔上下浮动 10%, 例如设置为15s，即 15+1.5 或者 15-1.5
				MaxConnectionAge:      30 * time.Second,      // 如果任意连接存活时间超过该值, 发送一个 GOAWAY
				MaxConnectionAgeGrace: 5 * time.Second,       // 在强制关闭连接之间, 允许有该值的时间完成 pending 的 rpc 请求
				Time:                  pool.KeepAliveTime,    // 每隔10秒ping一次客户端
				Timeout:               pool.KeepAliveTimeout, // 若回包在3s内返回则认为正常，否则连接将被回收
			}),

			// 流拦截器
			grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
				grpc_ctxtags.StreamServerInterceptor(),
				grpc_opentracing.StreamServerInterceptor(),
				grpc_prometheus.StreamServerInterceptor,
				grpc_zap.StreamServerInterceptor(global.LOG),
				grpc_auth.StreamServerInterceptor(AuthInterceptor),
				grpc_recovery.StreamServerInterceptor(),
			)),

			// 单元拦截器
			grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
				grpc_ctxtags.UnaryServerInterceptor(),
				grpc_opentracing.UnaryServerInterceptor(),
				grpc_prometheus.UnaryServerInterceptor,
				grpc_zap.UnaryServerInterceptor(global.LOG),
				grpc_auth.UnaryServerInterceptor(AuthInterceptor),
				grpc_recovery.UnaryServerInterceptor(),
			)),

			// 新建gRPC服务器实例,并开启TLS认证
			grpc.Creds(creds),
		)

		proto.RegisterCommonServiceServer(global.GRPCSServer, &CommonServiceServer{})

		err = global.GRPCSServer.Serve(lis)
		if err != nil {
			return err
		}
	}

	return nil
}
