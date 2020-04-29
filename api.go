package boulle

import (
	"context"
	"github.com/coreos/etcd/clientv3"
	"github.com/golang/glog"
	jsoniter "github.com/json-iterator/go"
	"time"
)

var (
	appName = "boulle"
)

type Client struct {
	idc      string
	app      string
	client   *clientv3.Client
	nic      string
	key      string
	ip       string
	interval int
	Expand   interface{}
}

func (n *Client) RegisterTicker(shopCh <-chan struct{}) {
	n.Register()
	timeTimer := time.NewTimer(time.Duration(n.interval) * time.Second)
	for {
		select {
		case <-timeTimer.C:
			n.Register()
			timeTimer.Reset(time.Duration(n.interval) * time.Second)
		case <-shopCh:
			glog.Infof("stop register...")
			return
		}
	}
}

func (n *Client) Register() {
	v := Data{
		Idc:            n.idc,
		App:            n.app,
		Ip:             n.ip,
		LastUpdateTime: time.Now(),
		Expand:         n.Expand,
	}
	data, err := jsoniter.Marshal(v)
	if err != nil {
		glog.Errorf("etcd data %#v json marshal error:%s", v, err.Error())
	} else {
		_, err = n.client.Put(context.Background(), n.key, string(data))
		if err != nil {
			glog.Errorf("etcd data %s  put error:%s", string(data), err.Error())
		} else {
			glog.V(6).Infof("put etcd  key:%s value:%s success", n.key, string(data))
		}
	}
}
func (n *Client) RegisterRemove() error {
	_, err := n.client.Delete(context.Background(), n.key)
	if err != nil {
		glog.Errorf("etcd key: %s  delete error:%s", n.key, err.Error())
		return err
	} else {
		glog.V(10).Infof("delete etcd  key:%s success", n.key)
		return nil
	}
}

//func Initialization(idc,  string, nics, endpoints []string, username, password, prefix string, timeout,
//	interval int,  expand interface{}) (*New, error) {
//	client, err := NewEtcdClient(endpoints, username, password, time.Duration(timeout)*time.Second)
//	if err != nil {
//		return nil, err
//	}
//
//	n := &Client{
//		idc:      idc,
//		app:  project,
//		client:   client,
//		nic:      nic,
//		ip:       ip,
//		key:      etcdKey(prefix, project, ip, viper.GetString("boulle.id")),
//		interval: interval,
//		Expand: expand,
//	}
//	return n, nil
//}

func NewWithEtcdClient(prefix, cir, idc, app, ip, id string, interval int, etcdClient *clientv3.Client,
	expand interface{}) (*Client,
	error) {
	n := &Client{
		idc:      idc,
		app:      app,
		client:   etcdClient,
		key:      EtcdKey(prefix, cir, idc, app, ip, id),
		ip:       ip,
		interval: interval,
		Expand:   expand,
	}
	return n, nil
}

//
//func InitializationWithViper() (*New, error) {
//	//nics := strings.Split(viper.GetString("boulle.nic"), ",")
//	//metricsPath := viper.GetString("boulle.metricsPath")
//	//interval := viper.GetInt("boulle.interval")
//	//port := viper.GetInt("boulle.port")
//	//labels := viper.GetString("bolle.labels")
//	err := viper.UnmarshalKey("boulle", &config)
//	if err != nil {
//		return nil, errors.New("boulle 获取配置失败")
//	}
//	glog.V(20).Infof("%s  config:%#v", appName, config)
//	return Initialization(config.Idc, config.Project, config.Nics, config.Etcd.Endpoints, config.Etcd.Username, config.Etcd.Password, config.Etcd.Prefix, config.Etcd.Timeout, config.Interval)
//}
