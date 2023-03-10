package dwebsocket

import (
	"db-go-websocket/internal/global"
	"db-go-websocket/utils"
	"encoding/json"
	"go.uber.org/zap"
	"sync"
)

// WsRequest 通用请求数据格式
type WsRequest struct {
	Code uint32      `json:"seq"`            // 消息的唯一Id
	Data interface{} `json:"data,omitempty"` // 数据 json
}

// WsResponse 通用返回数据格式
type WsResponse struct {
	Code uint32      `json:"seq"`            // 消息的唯一Id
	Data interface{} `json:"data,omitempty"` // 数据 json
}

type HandleFunc func(client *Client, message []byte) ([]byte, error)

var (
	handlers        = make(map[uint32]HandleFunc)
	handlersRWMutex sync.RWMutex
)

// Register 注册
func Register(key uint32, value HandleFunc) {
	handlersRWMutex.Lock()
	defer handlersRWMutex.Unlock()
	handlers[key] = value

	return
}

func getHandlers(code uint32) (value HandleFunc, ok bool) {
	handlersRWMutex.RLock()
	defer handlersRWMutex.RUnlock()

	value, ok = handlers[code]

	return
}

// Handle 处理数据
func Handle(client *Client, message []byte) {
	defer utils.PrintPanic()

	request := &WsRequest{}
	err := json.Unmarshal(message, request)
	if err != nil {
		global.LOG.Error("", zap.Error(err))
	}

	data, err := json.Marshal(request.Data)
	if err != nil {
		global.LOG.Error("", zap.Error(err))
		return
	}

	// 采用 map 注册的方式
	value, ok := getHandlers(request.Code)
	if !ok {
		global.LOG.Error("", zap.Error(err))
		return
	}

	resp, err := value(client, data)
	if err != nil {
		global.LOG.Error("", zap.Error(err))
		return
	}

	client.SendMsg(resp)
}
