package task

import (
	"db-go-websocket/internal/dwebsocket"
	"db-go-websocket/utils"
	"fmt"
	"time"
)

type Task1 struct{}

func init() {
	task := &Task1{}
	Timer(0*time.Second, 1*time.Second, task)
}

// Start 注册
func (t *Task1) Start() bool {
	defer utils.PrintPanic()
	dwebsocket.ClearTimeoutConnections()
	return true
}

// Stop 下线
func (t *Task1) Stop() bool {
	fmt.Println("stop")
	return true
}
