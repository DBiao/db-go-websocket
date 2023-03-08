package ws

import (
	"db-go-websocket/internal/global"
	"db-go-websocket/internal/global/statuscode"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"net/http"
)

const (
	// 最大的消息大小
	maxMessageSize = 8192
)

type WsController struct {
}

type renderData struct {
	ClientId string `json:"clientId"`
}

func (ws *WsController) WebSocket(ctx *gin.Context) {
	conn, err := (&websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		// 允许所有CORS跨域请求
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}).Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Errorf("upgrade error: %v", err)
		http.NotFound(ctx.Writer, ctx.Request)
		return
	}

	//设置读取消息大小上线
	conn.SetReadLimit(maxMessageSize)

	//解析参数
	systemId := global.CONFIG.System.SystemId
	if systemId == 0 {
		_ = Render(conn, "", "", statuscode.ERROR, "系统ID不能为空", []string{})
		_ = conn.Close()
		return
	}

	clientId := "util.GenClientId()"

	clientSocket := NewClient(clientId, string(systemId), conn)

	Manager.AddClient2SystemClient(string(systemId), clientSocket)

	//读取客户端消息
	clientSocket.Read()

	//if err = api.ConnRender(conn, renderData{ClientId: clientId}); err != nil {
	//	_ = conn.Close()
	//	return
	//}

	// 用户连接事件
	Manager.Connect <- clientSocket
}
