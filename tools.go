package boulle

import (
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"time"
)

func NewEtcdClient(endpoints []string, username, password string, timeout time.Duration) (*clientv3.Client, error) {
	config := clientv3.Config{
		Endpoints:   endpoints,
		Username:    username,
		Password:    password,
		DialTimeout: timeout,
	}
	return clientv3.New(config)
}

func EtcdKey(prefix, cir, idc, app, ip string, id string) string {
	if id == "" {
		id = "TSBPBVHGMC"
	}
	return fmt.Sprintf("%s/%s/%s/%s/%s-%s", prefix, app, cir, idc, ip, id)
}
