package config

import "encoding/json"

var AppConfig Config

type Config struct {
	Siurl     string   `json:"siurl"`
	Ibmid     string   `json:"ibmid"`
	ApiKey    string   `json:"apiKey"`
	Debug     bool     `json:"debug"`
	TenantId  string   `json:"tenantId"`
	SystemIds []string `json:"systemIDs"`
	Metrics   []string `json:"metrics"`
}

func LoadConfig(data []byte) error {
	return json.Unmarshal(data, &AppConfig)
}
