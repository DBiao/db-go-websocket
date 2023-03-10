package router

import (
	"db-go-websocket/internal/dwebsocket"
	"db-go-websocket/internal/global"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/snowlyg/multi"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

const (
	// 最大的消息大小
	maxMessageSize = 8192
)

type WsController struct {
}

func StartWebSocket(router *gin.Engine) {
	websocketHandler := &WsController{}
	router.GET("/dwebsocket", websocketHandler.WebSocket)
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
		global.LOG.Error("upgrade error", zap.Error(err))
		http.NotFound(ctx.Writer, ctx.Request)
		return
	}

	//设置读取消息大小上线
	conn.SetReadLimit(maxMessageSize)

	// 验证token是否存在
	token := ctx.Param("token")
	_, err = multi.AuthDriver.GetMultiClaims(token)
	if err != nil {
		return
	}

	clientId := ctx.Param("clientId")
	cId, _ := strconv.Atoi(clientId)
	appId := ctx.Param("appId")
	aId, _ := strconv.Atoi(appId)
	client := dwebsocket.NewClient(uint64(cId), uint64(aId), conn)

	dwebsocket.Manager.AddClient(client)

	client.Close()

	go client.Run(conn)

	for {
		_, b, err := conn.ReadMessage()
		if err != nil {
			return
		}
		client.InChan <- b
	}
}
