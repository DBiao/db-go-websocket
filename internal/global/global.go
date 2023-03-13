package global

import (
	config "db-go-websocket/conf"
	"github.com/Shopify/sarama"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net/http"
)

var (
	CONFIG      *config.Config
	VIPER       *viper.Viper
	LOG         *zap.Logger
	SERVER      *http.Server
	GRPCSServer *grpc.Server
	KAFKA       sarama.Consumer
	GRPCClient  *grpc.ClientConn
)
