package dserver

import (
	"context"
	"db-go-websocket/internal/global"
	"db-go-websocket/internal/proto"
	"fmt"
	"google.golang.org/grpc"
	"net"
)

func InitGrpcServer() error {
	// 如果是集群，则启用RPC进行通讯
	if global.CONFIG.System.IsCluster {
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", global.CONFIG.Grpc.Port))
		if err != nil {
			return err
		}

		global.GRPCSERVER = grpc.NewServer()
		proto.RegisterCommonServiceServer(global.GRPCSERVER, &CommonServiceServer{})

		err = global.GRPCSERVER.Serve(lis)
		if err != nil {
			return err
		}
	}

	return nil
}

type CommonServiceServer struct{}

func (this *CommonServiceServer) Send2Client(ctx context.Context, req *proto.Send2ClientReq) (*proto.Send2ClientReply, error) {
	return &proto.Send2ClientReply{}, nil
}

func (this *CommonServiceServer) CloseClient(ctx context.Context, req *proto.CloseClientReq) (*proto.CloseClientReply, error) {
	return &proto.CloseClientReply{}, nil
}

// BindGroup 添加分组到group
func (this *CommonServiceServer) BindGroup(ctx context.Context, req *proto.BindGroupReq) (*proto.BindGroupReply, error) {
	return &proto.BindGroupReply{}, nil
}

func (this *CommonServiceServer) Send2Group(ctx context.Context, req *proto.Send2GroupReq) (*proto.Send2GroupReply, error) {
	return &proto.Send2GroupReply{}, nil
}

func (this *CommonServiceServer) Send2System(ctx context.Context, req *proto.Send2SystemReq) (*proto.Send2SystemReply, error) {
	return &proto.Send2SystemReply{}, nil
}

// GetGroupClients 获取分组在线用户列表
func (this *CommonServiceServer) GetGroupClients(ctx context.Context, req *proto.GetGroupClientsReq) (*proto.GetGroupClientsReply, error) {
	return &proto.GetGroupClientsReply{}, nil
}
