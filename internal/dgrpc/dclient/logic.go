package dclient

import (
	"context"
	"db-go-websocket/internal/global"
	"db-go-websocket/internal/proto"
	"db-go-websocket/pkg/etcd"
	"fmt"
	log "github.com/sirupsen/logrus"
	"sync"
)

func SendRpc2Client(addr string, messageId, sendUserId, clientId string, code int, message string, data *string) {
	conn := global.POOLS

	log.WithFields(log.Fields{
		"host":     global.CONFIG.Grpc.Port,
		"port":     global.CONFIG.Http.Port,
		"add":      addr,
		"clientId": clientId,
		"msg":      data,
	}).Info("发送到服务器")

	c := proto.NewCommonServiceClient(conn)
	_, err := c.Send2Client(context.Background(), &proto.Send2ClientReq{
		MessageId:  messageId,
		SendUserId: sendUserId,
		ClientId:   clientId,
		Code:       int32(code),
		Message:    message,
		Data:       *data,
	})
	if err != nil {
		log.Errorf("failed to call: %v", err)
	}
}

func CloseRpcClient(addr string, clientId, systemId string) {
	conn := global.POOLS

	log.WithFields(log.Fields{
		"host":     global.CONFIG.Grpc.Port,
		"port":     global.CONFIG.Http.Port,
		"add":      addr,
		"clientId": clientId,
	}).Info("发送关闭连接到服务器")

	c := proto.NewCommonServiceClient(conn)
	_, err := c.CloseClient(context.Background(), &proto.CloseClientReq{
		SystemId: systemId,
		ClientId: clientId,
	})
	if err != nil {
		log.Errorf("failed to call: %v", err)
	}
}

// SendRpcBindGroup 绑定分组
func SendRpcBindGroup(addr string, systemId string, groupName string, clientId string, userId string, extend string) {
	conn := global.POOLS

	c := proto.NewCommonServiceClient(conn)
	_, err := c.BindGroup(context.Background(), &proto.BindGroupReq{
		SystemId:  systemId,
		GroupName: groupName,
		ClientId:  clientId,
		UserId:    userId,
		Extend:    extend,
	})
	if err != nil {
		log.Errorf("failed to call: %v", err)
	}
}

// SendGroupBroadcast 发送分组消息
func SendGroupBroadcast(messageId, sendUserId, groupName string, code int, message string, data *string) {
	etcd.GlobalSetting.ServerListLock.Lock()
	defer etcd.GlobalSetting.ServerListLock.Unlock()
	for _, addr := range etcd.GlobalSetting.ServerList {
		conn := global.POOLS
		fmt.Println(addr)
		c := proto.NewCommonServiceClient(conn)
		_, err := c.Send2Group(context.Background(), &proto.Send2GroupReq{
			MessageId:  messageId,
			SendUserId: sendUserId,
			GroupName:  groupName,
			Code:       int32(code),
			Message:    message,
			Data:       *data,
		})
		if err != nil {
			log.Errorf("failed to call: %v", err)
		}
	}
}

// SendSystemBroadcast 发送系统信息
func SendSystemBroadcast(systemId string, messageId, sendUserId string, code int, message string, data *string) {
	etcd.GlobalSetting.ServerListLock.Lock()
	defer etcd.GlobalSetting.ServerListLock.Unlock()
	for _, addr := range etcd.GlobalSetting.ServerList {
		conn := global.POOLS
		fmt.Println(addr)
		c := proto.NewCommonServiceClient(conn)
		_, err := c.Send2System(context.Background(), &proto.Send2SystemReq{
			SystemId:   systemId,
			MessageId:  messageId,
			SendUserId: sendUserId,
			Code:       int32(code),
			Message:    message,
			Data:       *data,
		})
		if err != nil {
			log.Errorf("failed to call: %v", err)
		}
	}
}

func GetOnlineListBroadcast(systemId *string, groupName *string) (clientIdList []string) {
	etcd.GlobalSetting.ServerListLock.Lock()
	defer etcd.GlobalSetting.ServerListLock.Unlock()

	serverCount := len(etcd.GlobalSetting.ServerList)

	onlineListChan := make(chan []string, serverCount)
	var wg sync.WaitGroup

	wg.Add(serverCount)
	for _, addr := range etcd.GlobalSetting.ServerList {
		go func(addr string) {
			conn := global.POOLS
			fmt.Println(addr)
			c := proto.NewCommonServiceClient(conn)
			response, err := c.GetGroupClients(context.Background(), &proto.GetGroupClientsReq{
				SystemId:  *systemId,
				GroupName: *groupName,
			})
			if err != nil {
				log.Errorf("failed to call: %v", err)
			} else {
				onlineListChan <- response.List
			}
			wg.Done()

		}(addr)
	}

	wg.Wait()

	for i := 1; i <= serverCount; i++ {
		list, ok := <-onlineListChan
		if ok {
			clientIdList = append(clientIdList, list...)
		} else {
			return
		}
	}
	close(onlineListChan)

	return
}
