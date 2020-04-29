package boulle

import "time"

type Data struct {
	Idc            string      `json:"idc"`
	App            string      `json:"project"`
	Ip             string      `json:"ip"`
	LastUpdateTime time.Time   `json:"lastUpdateTime"`
	Expand         interface{} `json:"expand,omitempty"`
}

type Etcd struct {
	Endpoints []string `json:"endpoints"`
	Username  string   `json:"username"`
	Password  string   `json:"password"`
	Prefix    string   `json:"prefix"`
	Timeout   int      `json:"timeout"`
}

type Config struct {
	Cir      string `json:"cir"`
	Idc      string `json:"idc"`
	App      string `json:"project"`
	Interval int    `json:"interval"`
	Etcd     Etcd   `json:"etcd"`
}
