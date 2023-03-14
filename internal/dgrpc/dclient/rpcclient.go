package dclient

import (
	"db-go-websocket/pkg/etcd"
	"fmt"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"github.com/shimingyah/pool"
	"go.etcd.io/etcd/client/v3/naming/resolver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
	"sync"
	"time"
)

func InitGrpcClient() (*grpc.ClientConn, error) {
	etcdResolver, err := resolver.NewBuilder(etcd.GetInstance())

	//从输入的证书文件中为客户端构造TLS凭证
	creds, err := credentials.NewClientTLSFromFile("../pkg/tls/server.pem", "go-grpc-example")
	if err != nil {
		panic(err)
	}

	retryOps := []grpc_retry.CallOption{
		grpc_retry.WithMax(2),
		grpc_retry.WithPerRetryTimeout(time.Second * 2),
		grpc_retry.WithBackoff(grpc_retry.BackoffLinearWithJitter(time.Second/2, 0.2)),
	}
	retryInterceptor := grpc_retry.UnaryClientInterceptor(retryOps...)

	opts := []grpc.DialOption{
		// 重试机制
		grpc.WithUnaryInterceptor(retryInterceptor),

		// 心跳检查
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                20 * time.Second, // 每隔20秒ping一次服务器
			Timeout:             5 * time.Second,  // 若回包在5s内返回则认为正常，否则连接将被回收
			PermitWithoutStream: false,
		}),

		// TLS加密
		grpc.WithTransportCredentials(creds),

		// 负载均衡策略
		grpc.WithResolvers(etcdResolver),
		grpc.WithDefaultServiceConfig(roundrobin.Name),
	}

	conn, err := grpc.Dial(fmt.Sprintf("etcd:///%s", "sericeName"), opts...)
	if err != nil {
		return nil, err
	}

	return conn, err
}

func InitGrpcClientPool(addrs ...string) (*sync.Map, error) {
	pools := &sync.Map{}

	if len(addrs) == 0 {
		return pools, nil
	}

	for _, addr := range addrs {
		pool, err := pool.New(addr, pool.DefaultOptions)
		if err != nil {
			return nil, err
		}
		pools.Store(addr, pool)
	}

	return pools, nil
}

//func grpcConn(addr string) (*grpc.ClientConn, error) {
//	p, ok := global.POOLS.Load(addr)
//	if !ok {
//		return nil, errors.New("addr not exist")
//	}
//
//	conn, err := p.(pool.Pool).Get()
//	if err != nil {
//		return nil, err
//	}
//
//	return conn.Value(), nil
//}
