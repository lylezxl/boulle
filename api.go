package boulle

import (
	"context"
	"encoding/json"
	"github.com/coreos/etcd/clientv3"
	"github.com/golang/glog"
	jsoniter "github.com/json-iterator/go"
	"github.com/spf13/viper"
	"strings"
	"time"
)

type New struct {
	client      *clientv3.Client
	nic         string
	port        int
	metricsPath string
	key         string
	ip          string
	labels      map[string]string
	interval    int
}

func (n *New) RegisterTicker(shopCh <-chan struct{}) {
	timeTimer := time.NewTimer(time.Duration(n.interval) * time.Second)
	for {
		select {
		case <-timeTimer.C:
			v := Data{
				Nic:            n.nic,
				Ip:             n.ip,
				Port:           n.port,
				MetricsPath:    n.metricsPath,
				LastUpdateTime: getCurrentTime(),
			}
			data, err := jsoniter.Marshal(v)
			if err != nil {
				glog.Errorf("etcd data %#v json marshal error:%s", v, err.Error())
			} else {
				_, err = n.client.Put(context.Background(), n.key, string(data))
				if err != nil {
					glog.Errorf("etcd data %s  put error:%s", string(data), err.Error())
				} else {
					glog.V(10).Infof("put etcd  key:%s value:%s success", n.key, string(data))
				}
			}
			timeTimer.Reset(time.Minute * 1)
		case <-shopCh:
			glog.Infof("stop register:%s")
		}
	}
}
func (n *New) RegisterRemove() error {
	_, err := n.client.Delete(context.Background(), n.key)
	if err != nil {
		glog.Errorf("etcd key: %s  delete error:%s", n.key, err.Error())
		return err
	} else {
		glog.V(10).Infof("delete etcd  key:%s success", n.key)
		return nil
	}
}
func Initialization(nics, endpoints []string, username, password, metricsPath, prefix string, timeout, interval, port int, labels map[string]string) (*New, error) {
	nic, ip := getIp(nics)
	if nic == "" || ip == "" {
		return nil, NotIpaddress
	}
	c, err := newEtcdClient(endpoints, username, password, time.Duration(timeout)*time.Second)
	if err != nil {
		return nil, err
	}
	if metricsPath == "" {
		metricsPath = "/metrics"
	}
	return &New{
		client:      c,
		nic:         nic,
		port:        port,
		metricsPath: metricsPath,
		ip:          ip,
		key:         etcdKey(prefix, ip, port),
		labels:      labels,
		interval:    interval,
	}, nil
}

func InitializationWithViper() (*New, error) {
	nics := strings.Split(viper.GetString("boulle.nic"), ",")
	metricsPath := viper.GetString("boulle.metricsPath")
	interval := viper.GetInt("boulle.interval")
	port := viper.GetInt("boulle.port")
	labels := viper.GetString("bolle.labels")
	var m map[string]string
	if labels != "" {
		err := json.Unmarshal([]byte(labels), m)
		if err != nil {
			glog.V(15).Infof("pasre prometheus label data:%s  error:%s", labels, err)
		}
	} else {
		glog.V(15).Infof("prometheus no label")

	}
	endpoints := strings.Split(viper.GetString("boulle.etcd.endpoints"), ",")
	username := viper.GetString("boulle.etcd.username")
	password := viper.GetString("boulle.etcd.password")
	prefix := viper.GetString("boulle.etcd.prefix")
	timeout := viper.GetInt("boulle.etcd.timeout")

	//
	glog.V(20).Infof("viper config nic:%#v", nics)
	glog.V(20).Infof("viper config metricsPath:%s", metricsPath)
	glog.V(20).Infof("viper config labels:%s", labels)
	glog.V(20).Infof("viper config interval:%d", interval)
	glog.V(20).Infof("viper config port:%d", port)
	glog.V(20).Infof("viper config etcd endpoints:%#v", endpoints)
	glog.V(20).Infof("viper config etcd username:\"%s\"", username)
	glog.V(20).Infof("viper config etcd prefix:%s", prefix)
	glog.V(20).Infof("viper config etcd password:\"%s\"", password)
	glog.V(20).Infof("viper config etcd timeout:%d", timeout)
	return Initialization(nics, endpoints, username, password, metricsPath, prefix, timeout, interval, port, m)
}
