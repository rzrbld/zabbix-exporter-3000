package config

import (
	"os"
	strconv "strconv"
)

var (
	Server         = getEnv("ZABBIX_API_ENDPOINT", "http://zabbix/api_jsonrpc.php")
	User           = getEnv("ZABBIX_USER", "admin")
	Password       = getEnv("ZABBIX_PASSWORD", "admin")
	SslSkip, _     = strconv.ParseBool(getEnv("ZABBIX_SKIP_SSL", "true"))
	MainHostPort   = getEnv("ZE3000_HOST_PORT", "localhost:8080")
)

func getEnv(key, fallback string) string {
	value, exist := os.LookupEnv(key)

	if !exist {
		return fallback
	}

	return value
}
