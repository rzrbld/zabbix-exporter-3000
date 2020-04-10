package config

import (
	"os"
	strconv "strconv"
  "strings"
)

var (
	Server         = getEnv("ZABBIX_API_ENDPOINT", "http://zabbix/api_jsonrpc.php")
	User           = getEnv("ZABBIX_USER", "admin")
	Password       = getEnv("ZABBIX_PASSWORD", "admin")
	SslSkip, _     = strconv.ParseBool(getEnv("ZABBIX_SKIP_SSL", "true"))
	MainHostPort   = getEnv("ZE3000_HOST_PORT", "localhost:8080")
	MetricName     = getEnv("ZE3000_ZABBIX_METRIC_NAME", "zabbix_exporter_metric")
	SourceRefresh  = getEnv("ZE3000_ZABBIX_REFRESH_DELAY_SEC", "10")
	MetricLabels   = strings.TrimSpace(getEnv("ZE3000_ZABBIX_METRIC_LABELS", "name,key_,hosts>name"))
	// MetricLabels   = strings.TrimSpace(getEnv("ZE3000_ZABBIX_METRIC_LABELS", "name,key_,itemid"))
	Query          = getEnv("ZE3000_ZABBIX_QUERY_PARAMS", `{     "jsonrpc": "2.0",     "method": "item.get",     "params": {     	"itemids":["330254","329514","178909"],         "output": ["name","key_","description","lastvalue"],         "selectHosts": ["name","status","host"],         "selectInterfaces": ["ip","dns"],         "sortfield":"key_"     },     "auth": "%auth-token%",     "id": 1 }`)
)

func getEnv(key, fallback string) string {
	value, exist := os.LookupEnv(key)

	if !exist {
		return fallback
	}

	return value
}
