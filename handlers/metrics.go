package handlers

import (
	"time"
  "encoding/json"
  "strconv"
  "log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
  zbx "github.com/rzrbld/zabbix-exporter-3000/zabbix"
  cnf "github.com/rzrbld/zabbix-exporter-3000/config"
)

var sourceRefreshSec, _ = strconv.Atoi(cnf.SourceRefresh)



var objectsCount = promauto.NewGauge(prometheus.GaugeOpts{
	Name: "objects_count_current",
	Help: "number of objects on cluster",
})

var itemsMetric = promauto.NewGaugeVec(prometheus.GaugeOpts{
	Name: cnf.MetricName,
	Help: "bucket size in kbytes",
}, []string{"name","key_","itemid"})

func RecordMetrics() {
	go func() {
		for {
      items, err := zbx.Session.Do(zbx.Query)
      if err != nil {
        log.Print("ERROR While Do request: ",err)
      }else{
        var results []map[string]interface{}
        json.Unmarshal(items.Body, &results)
        if err != nil {
      		log.Print("ERROR While convert response to JSON: ",err)
        }

        for key, result := range results {

      		log.Print("Reading Value for Key :", key)
      		// //Reading each value by its key
      		// fmt.Println("Id :", result["itemid"],
      		// 	"- Name :", result["name"],
      		// 	"- Department :", result["key_"],
      		// 	"- Designation :", result["hosts"])
          f, _ := strconv.ParseFloat(result["lastvalue"].(string), 64)
          itemsMetric.WithLabelValues(string(result["name"].(string)),string(result["key_"].(string)),string(result["itemid"].(string))).Set(float64(f))
      	}
      }
      time.Sleep(time.Duration(sourceRefreshSec) * time.Second)
		}
	}()
}
