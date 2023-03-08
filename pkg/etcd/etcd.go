package etcd

import (
	"context"
	"db-go-websocket/internal/global"
	"github.com/coreos/etcd/clientv3"
	log "github.com/sirupsen/logrus"
	"net"
	"sync"
	"time"
)

const (
	//ETCD服务列表路径
	ETCD_SERVER_LIST = "/wsServers/"
	//账号信息前缀
	ETCD_PREFIX_ACCOUNT_INFO = "ws/account/"
)

var etcdKvClient *clientv3.Client
var mu sync.Mutex

// InitEtcd ETCD注册发现服务 将服务器地址、端口注册到etcd中
func InitEtcd() error {
	if global.CONFIG.System.IsCluster {
		// 注册租约
		ser, err := NewServiceReg(global.CONFIG.Etcd.Endpoints, 5)
		if err != nil {
			return err
		}

		hostPort := net.JoinHostPort(GlobalSetting.LocalHost, "etcd.GlobalSetting.CommonSetting.RPCPort")
		// 添加key
		err = ser.PutService(ETCD_SERVER_LIST+hostPort, hostPort)
		if err != nil {
			return err
		}

		cli, err := NewClientDis(global.CONFIG.Etcd.Endpoints)
		if err != nil {
			return err
		}
		_, err = cli.GetService(ETCD_SERVER_LIST)
		if err != nil {
			return err
		}
	}

	return nil
}

func GetInstance() *clientv3.Client {
	if etcdKvClient == nil {
		if client, err := clientv3.New(clientv3.Config{
			Endpoints:   global.CONFIG.Etcd.Endpoints,
			DialTimeout: 5 * time.Second,
		}); err != nil {
			log.Error(err)
			return nil
		} else {
			//创建时才加锁
			mu.Lock()
			defer mu.Unlock()
			etcdKvClient = client
			return etcdKvClient
		}

	}
	return etcdKvClient
}

func Put(key, value string) error {
	_, err := GetInstance().Put(context.Background(), key, value)
	return err
}

func Get(key string) (resp *clientv3.GetResponse, err error) {
	resp, err = GetInstance().Get(context.Background(), key)
	return resp, err
}
