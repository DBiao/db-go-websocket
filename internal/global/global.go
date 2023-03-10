package global

import (
	config "db-go-websocket/conf"
	"github.com/Shopify/sarama"
	"google.golang.org/grpc"
	"net/http"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	CONFIG     *config.Config
	VIPER      *viper.Viper
	LOG        *zap.Logger
	SERVER     *http.Server
	GRPCSERVER *grpc.Server
	KAFKA      sarama.Consumer
)
