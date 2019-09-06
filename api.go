package boulle

import (
	"context"
	"github.com/coreos/etcd/clientv3"
	"github.com/golang/glog"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"time"
)

var (
	appName = "boulle"
	config  = Config{}
)

type New struct {
	project    string
	client     *clientv3.Client
	nic        string
	key        string
	ip         string
	interval   int
	Prometheus *Prometheus
}

func (n *New) RegisterTicker(shopCh <-chan struct{}) {
	timeTimer := time.NewTimer(time.Duration(n.interval) * time.Second)
	for {
		select {
		case <-timeTimer.C:
			v := Data{
				Project:        n.project,
				Nic:            n.nic,
				Ip:             n.ip,
				Prometheus:     n.Prometheus,
				LastUpdateTime: GetCurrentTime(),
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
func Initialization(project string, nics, endpoints []string, username, password, prefix string, timeout, interval int, enable_prometheus bool, prometheus *Prometheus) (*New, error) {
	nic, ip := getIp(nics)
	if nic == "" || ip == "" {
		return nil, NotIpaddress
	}
	c, err := NewEtcdClient(endpoints, username, password, time.Duration(timeout)*time.Second)
	if err != nil {
		return nil, err
	}

	n := &New{
		project:  project,
		client:   c,
		nic:      nic,
		ip:       ip,
		key:      etcdKey(prefix, project, ip, RandString(10)),
		interval: interval,
	}
	if enable_prometheus {
		if prometheus.Metrics == "" {
			prometheus.Metrics = "/metrics"
		}
		n.Prometheus = prometheus
	}
	return n, nil
}

func InitializationWithViper() (*New, error) {
	//nics := strings.Split(viper.GetString("boulle.nic"), ",")
	//metricsPath := viper.GetString("boulle.metricsPath")
	//interval := viper.GetInt("boulle.interval")
	//port := viper.GetInt("boulle.port")
	//labels := viper.GetString("bolle.labels")
	err := viper.UnmarshalKey("boulle", &config)
	if err != nil {
		return nil, errors.New("boulle 获取配置失败")
	}
	glog.V(20).Infof("%s  config:%#v", appName, config)
	return Initialization(config.Project, config.Nics, config.Etcd.Endpoints, config.Etcd.Username, config.Etcd.Password, config.Etcd.Prefix, config.Etcd.Timeout, config.Interval, config.Enable_promethues, &config.Prometheus)
}
