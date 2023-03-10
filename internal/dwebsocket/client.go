package dwebsocket

import (
	"db-go-websocket/utils"
	"github.com/gorilla/websocket"
	"sync/atomic"
	"time"
)

const (
	// 用户连接超时时间
	heartbeatExpirationTime = 60
)

type Client struct {
	UserId        uint64          // 业务端标识用户ID
	AppId         uint64          // 登录的平台Id app/web/ios
	Conn          *websocket.Conn // 用户连接
	FirstTime     uint64          // 首次连接时间
	HeartbeatTime uint64          // 用户上次心跳时间
	Addr          string          // 客户端地址
	OutChan       chan []byte
	InChan        chan []byte
	QuitChan      chan struct{}
	CloseFlag     int32
}

type SendData struct {
	Code int
	Msg  string
	Data *interface{}
}

func NewClient(clientId, appId uint64, conn *websocket.Conn) *Client {
	return &Client{
		UserId:        clientId,
		AppId:         appId,
		Conn:          conn,
		FirstTime:     uint64(time.Now().Unix()),
		HeartbeatTime: uint64(time.Now().Unix()),
		Addr:          conn.RemoteAddr().String(),
		OutChan:       make(chan []byte, 1),
		InChan:        make(chan []byte, 1),
		QuitChan:      make(chan struct{}, 1),
	}
}

func (c *Client) Run(ws *websocket.Conn) {
	defer utils.PrintPanic()

	for {
		select {
		case in, ok := <-c.InChan:
			if !ok {
				return
			}
			Handle(c, in)
		case out, ok := <-c.OutChan:
			if !ok {
				return
			}
			err := ws.WriteMessage(websocket.TextMessage, out)
			if err != nil {
				return
			}
		case <-c.QuitChan:
			return
		}
	}
}

// Close 关闭监听
func (c *Client) Close() {
	c.Conn.SetCloseHandler(func(code int, text string) error {
		c.stop()
		c.Conn.Close()
		return nil
	})
}

func (c *Client) stop() {
	if atomic.CompareAndSwapInt32(&c.CloseFlag, 0, 1) {
		select {
		case c.QuitChan <- struct{}{}:
		default:
		}
	}
}

// SendMsg 发送数据
func (c *Client) SendMsg(msg []byte) {
	if c == nil {
		return
	}

	c.OutChan <- msg
}

// IsHeartbeatTimeout 心跳超时
func (c *Client) IsHeartbeatTimeout(currentTime uint64) (timeout bool) {
	if c.HeartbeatTime+heartbeatExpirationTime <= currentTime {
		timeout = true
	}

	return
}
