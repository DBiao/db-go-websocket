package handle

import (
	"db-go-websocket/internal/dwebsocket"
	"time"
)

func init() {
	// Websocket 路由
	dwebsocket.Register(10000, HeartbeatController)
}

func HeartbeatController(client *dwebsocket.Client, message []byte) ([]byte, error) {
	client.HeartbeatTime = uint64(time.Now().Unix())
	return nil, nil
}
