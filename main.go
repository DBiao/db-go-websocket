package main

import "db-go-websocket/internal/app"

//go:generate go env -w GO111MODULE=on
//go:generate go env -w GOPROXY=https://goproxy.cn
//go:generate go mod tidy
//go:generate go mod download

func main() {
	// 初始化服务
	app.Start()

	// 优雅关闭连接
	defer app.Close()
}
