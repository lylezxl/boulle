package boulle

type Data struct {
	Idc            string `json:"idc"`
	Project        string `json:"project"`
	Nic            string `json:"nic"`
	Ip             string `json:"ip"`
	Prometheus     *P     `json:"prometheus,omitempty"`
	LastUpdateTime string `json:"lastUpdateTime"`
}

type Prometheus struct {
	Metrics string  `json:"metrics"`
	Port    int     `json:"port"`
	Labels  []Label `json:"labels"`
}

type Label struct {
	Key   string `json:"key"`
	Value string `json:"value"`
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

type P struct {
	Metrics string            `json:"metrics"`
	Port    int               `json:"port"`
	Labels  map[string]string `json:"labels"`
}
