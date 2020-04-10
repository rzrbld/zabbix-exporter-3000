package handlers

import (
	"time"
  "encoding/json"
  "strconv"
  "strings"
  "log"
  // "fmt"
  // "reflect"
	"github.com/prometheus/client_golang/prometheus"
  zbx "github.com/rzrbld/zabbix-exporter-3000/zabbix"
  cnf "github.com/rzrbld/zabbix-exporter-3000/config"
)


func RecordMetrics() {

  var sourceRefreshSec, _ = strconv.Atoi(cnf.SourceRefresh)
  var labelsSliceRaw = strings.Split(cnf.MetricLabels, ",")
  // labels in prom format
  var labelsSlicePrometheus []string
  // labels with path "a>b"
  var labelsSliceComplex []string
  // labels average
  var labelsSliceAvg []string

  for _, vl := range labelsSliceRaw{
    if strings.Contains(vl, ">") {
      labelsSlicePrometheus = append(labelsSlicePrometheus, strings.Replace(vl, ">", "_", -1))
      labelsSliceComplex = append(labelsSliceComplex, vl)
    }else{
      labelsSlicePrometheus = append(labelsSlicePrometheus, vl)
      labelsSliceAvg = append(labelsSliceAvg, vl)
    }
  }

  log.Print("labels_prom   :", labelsSlicePrometheus)
  log.Print("labels_complex:", labelsSliceComplex)
  log.Print("labels_average:", labelsSliceAvg)

  var itemsMetric = prometheus.NewGaugeVec(prometheus.GaugeOpts{
    Namespace: "our_company",
  	Subsystem: "blob_storage",
  	Name: cnf.MetricName,
  	Help: "bucket size in kbytes",
  }, labelsSlicePrometheus)

  prometheus.MustRegister(itemsMetric)

	go func() {
		for {
      items, err := zbx.Session.Do(zbx.Query)
      if err != nil {
        log.Print("ERROR While Do request: ",err)
      }else{
        var results []map[string]interface{}
        json.Unmarshal(items.Body, &results)

        for _, result := range results {

          // fmt.Println(reflect.TypeOf(labelsSlice))
          // fmt.Println(reflect.TypeOf(result))
          // fmt.Println(reflect.TypeOf(result["name"]))
          // fmt.Println(reflect.TypeOf(result["interfaces"]))

          labelsWithValues := make(map[string]string)


          if len(labelsSliceAvg) > 0 {
            for _, vAvg := range labelsSliceAvg{
              labelsWithValues[vAvg] = result[vAvg].(string)
            }
          }

          if len(labelsSliceComplex) > 0 {
            for _, vCplx := range labelsSliceComplex{

              var promLabel = strings.Replace(vCplx, ">", "_", -1)
              var path = strings.Split(vCplx, ">")

              log.Print("LENGTH >>>> ", len(result[path[0]].([]interface{})) )

              if(len(result[path[0]].([]interface{})) > 0){
                for _,cplx := range result[path[0]].([]interface{}) {
                  subCplx := cplx.(map[string]interface{})
      	          labelsWithValues[promLabel] = subCplx[path[1]].(string)
      		      }
              } else {
                  labelsWithValues[promLabel] = "NA"
              }
            }
          }

          f, _ := strconv.ParseFloat(result["lastvalue"].(string), 64)
          itemsMetric.With(labelsWithValues).Set(float64(f))
      	}
      }
      time.Sleep(time.Duration(sourceRefreshSec) * time.Second)
		}
	}()
}
