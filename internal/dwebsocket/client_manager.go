package dwebsocket

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

var Manager = newClientManager() // 管理者

// ClientManager 连接管理
type ClientManager struct {
	Clients           map[uint64]*Client // 全部的连接
	ClientsLock       sync.RWMutex       // 读写锁
	Connect           chan *Client       // 连接处理
	DisConnect        chan *Client       // 断开连接处理
	SystemClientsLock sync.RWMutex
}

func newClientManager() (clientManager *ClientManager) {
	clientManager = &ClientManager{
		Clients:    make(map[uint64]*Client),
		Connect:    make(chan *Client, 10000),
		DisConnect: make(chan *Client, 10000),
	}

	return
}

// AddClient 添加客户端
func (manager *ClientManager) AddClient(client *Client) {
	manager.ClientsLock.Lock()
	defer manager.ClientsLock.Unlock()
	manager.Clients[client.UserId] = client
}

// GetAllClient 获取所有的客户端
func (manager *ClientManager) GetAllClient() map[uint64]*Client {
	manager.ClientsLock.RLock()
	defer manager.ClientsLock.RUnlock()

	return manager.Clients
}

func (manager *ClientManager) GetAll() []*Client {
	manager.ClientsLock.RLock()
	defer manager.ClientsLock.RUnlock()

	var client []*Client

	for _, value := range manager.Clients {
		client = append(client, value)
	}

	return client
}

// GetCount 客户端数量
func (manager *ClientManager) GetCount() int {
	manager.ClientsLock.RLock()
	defer manager.ClientsLock.RUnlock()
	return len(manager.Clients)
}

// DelClient 删除用户
func (manager *ClientManager) DelClient(clientId uint64) {
	manager.ClientsLock.Lock()
	defer manager.ClientsLock.Unlock()

	delete(manager.Clients, clientId)
}

// GetClientById 通过clientId获取
func (manager *ClientManager) GetClientById(clientId uint64) (*Client, error) {
	manager.ClientsLock.RLock()
	defer manager.ClientsLock.RUnlock()

	if client, ok := manager.Clients[clientId]; !ok {
		return nil, errors.New("客户端不存在")
	} else {
		return client, nil
	}
}

// ClearTimeoutConnections 定时清理超时连接
func ClearTimeoutConnections() {
	currentTime := uint64(time.Now().Unix())
	clients := Manager.GetAll()
	for _, client := range clients {
		if client.IsHeartbeatTimeout(currentTime) {
			fmt.Println("心跳时间超时 关闭连接", client.Addr, client.UserId, client.FirstTime, client.HeartbeatTime)
			client.Conn.Close()
		}
	}
}
