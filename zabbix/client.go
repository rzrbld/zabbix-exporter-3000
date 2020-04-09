package zabbix

import (
  "crypto/tls"
	"fmt"
	"log"
	"net/http"
	"github.com/cavaliercoder/go-zabbix"
  cnf "github.com/rzrbld/zabbix-exporter-3000/config"
)

func Connect() (*zabbix.Session, error) {
  client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: cnf.SslSkip}}}

	cache := zabbix.NewSessionFileCache().SetFilePath("./zabbix_session")
	session, err := zabbix.CreateClient(cnf.Server).
		WithCache(cache).
		WithHTTPClient(client).
		WithCredentials(cnf.User, cnf.Password).
		Connect()
	if err != nil {
		log.Fatalf("%v\n", err)
	}

	version, err := session.GetVersion()

	if err != nil {
		panic(err)
	}

	fmt.Printf("Connected to Zabbix API v%s \r\n", version)
  return session, err
}
