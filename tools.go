package boulle

import (
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/golang/glog"
	"math/rand"
	"net"
	"time"
)

func GetCurrentTime() string {
	return time.Now().Format("2006-01-02T15:04:05Z07:00")
}

func NewEtcdClient(endpoints []string, username, password string, timeout time.Duration) (*clientv3.Client, error) {
	config := clientv3.Config{
		Endpoints:   endpoints,
		Username:    username,
		Password:    password,
		DialTimeout: timeout,
	}
	return clientv3.New(config)
}

func getIp(nicList []string) (nic, ip string) {
	for _, n := range nicList {
		if ip := getNicIp(n); ip != "" {
			return n, ip
		}
	}
	return "", ""
}

func getNicIp(name string) string {
	interfaceAddr, err := net.InterfaceByName(name)
	if err != nil {
		glog.V(20).Infof("fail to get net specific interface addrss name:%s err: %v", name, err)
	}
	addresses, err := interfaceAddr.Addrs()
	if err != nil {
		glog.V(20).Infof("fail to get net specific interface addrss name:%s err: %v", name, err)

	}
	for _, addr := range addresses {
		if ipnet, ok := addr.(*net.IPNet); !ok {
			glog.V(20).Infof("fail to get net specific interface addrss name:%s err: %v", name, err)

		} else {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}

	}
	return ""
}

func etcdKey(prefix, project, ip string, id string) string {
	if id == "" {
		id = "TSBPBVHGMC"
	}
	return fmt.Sprintf("%s/%s/%s-%s", prefix, project, ip, id)
}

var r *rand.Rand

func init() {
	r = rand.New(rand.NewSource(time.Now().Unix()))
}

// RandString 生成随机字符串
func RandString(len int) string {
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		b := r.Intn(26) + 65
		bytes[i] = byte(b)
	}
	return string(bytes)
}
