package dclient

import (
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"github.com/shimingyah/pool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"sync"
	"time"
)

func InitGrpcClient(addr string) (*grpc.ClientConn, error) {
	retryOps := []grpc_retry.CallOption{
		grpc_retry.WithMax(2),
		grpc_retry.WithPerRetryTimeout(time.Second * 2),
		grpc_retry.WithBackoff(grpc_retry.BackoffLinearWithJitter(time.Second/2, 0.2)),
	}
	retryInterceptor := grpc_retry.UnaryClientInterceptor(retryOps...)

	opts := []grpc.DialOption{grpc.WithUnaryInterceptor(retryInterceptor),
		// 心跳检查
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                20 * time.Second, // 每隔20秒ping一次服务器
			Timeout:             5 * time.Second,  // 若回包在5s内返回则认为正常，否则连接将被回收
			PermitWithoutStream: false,
		})}

	conn, err := grpc.Dial(addr, opts...)
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
