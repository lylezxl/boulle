package boulle

type Data struct {
	Nic            string `json:"nic"`
	Ip             string `json:"ip"`
	Port           int    `json:"port"`
	MetricsPath    string `json:"metricsPath"`
	LastUpdateTime string `json:"lastUpdateTime"`
}
