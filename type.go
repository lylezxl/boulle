package boulle

type Data struct {
	Idc            string      `json:"idc"`
	Project        string      `json:"project"`
	Nic            string      `json:"nic"`
	Ip             string      `json:"ip"`
	Prometheus     *Prometheus `json:"prometheus,omitempty"`
	LastUpdateTime string      `json:"lastUpdateTime"`
}

type Prometheus struct {
	Metrics string `json:"metrics"`
	Port    int    `json:"port"`
	Labels  string `json:"labels"`
}
type Etcd struct {
	Endpoints []string `json:"endpoints"`
	Username  string   `json:"username"`
	Password  string   `json:"password"`
	Prefix    string   `json:"prefix"`
	Timeout   int      `json:"timeout"`
}
type Config struct {
	Idc               string     `json:"idc"`
	Project           string     `json:"project"`
	Nics              []string   `json:"nics"`
	Interval          int        `json:"interval"`
	Enable_promethues bool       `json:"enable_promethues"`
	Prometheus        Prometheus `json:"prometheus"`
	Etcd              Etcd       `json:"etcd"`
}
