package handlers

import (
	"time"
  "encoding/json"
  "strconv"
  "strings"
  "log"
  "regexp"
	"github.com/prometheus/client_golang/prometheus"
  zbx "github.com/rzrbld/zabbix-exporter-3000/zabbix"
  cnf "github.com/rzrbld/zabbix-exporter-3000/config"
)

var sourceRefreshSec, _ = strconv.Atoi(cnf.SourceRefresh)
var labelsSliceRaw = strings.Split(cnf.MetricLabels, ",")
// labels in prom format
var labelsSlicePrometheus []string
// labels with path "a>b"
var labelsSliceComplex []string
// labels average
var labelsSliceAvg []string

// var itemsMetric *prometheus.GaugeVec
var metricsSlice []*prometheus.GaugeVec

func cleanUpName(name string)(string){
  reg, err := regexp.Compile("[^a-zA-Z0-9]+")
  if err != nil {
     log.Fatal(err)
  }
  cleanName := reg.ReplaceAllString(strings.ToLower(name), "")
  return cleanName
}

func buildMetrics() {
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

  var results = queryZabbix()

  for k, result := range results {
    //clean up metric name
    cleanName := cleanUpName(result[cnf.MetricNameField].(string))
    fullName := cnf.MetricNamePrefix+"_"+cleanName
    if cnf.RandomizeNames {
      fullName = cnf.MetricNamePrefix+"_"+cleanName+"_"+strconv.Itoa(k)
    }


    var itemsMetric = prometheus.NewGaugeVec(prometheus.GaugeOpts{
      Namespace: cnf.MetricNamespace,
    	Subsystem: cnf.MetricSubsystem,
    	Name: fullName,
    	Help: result[cnf.MetricHelpField].(string),
    }, labelsSlicePrometheus)
    metricsSlice = append(metricsSlice,itemsMetric)
    if cnf.StrictRegister {
      prometheus.MustRegister(itemsMetric)
    }else{
      prometheus.Register(itemsMetric)
    }
  }
}

func queryZabbix()([]map[string]interface{}) {
  items, err := zbx.Session.Do(zbx.Query)
  if err != nil {
    log.Fatal("ERROR While Do request: ",err)
  }

  var results []map[string]interface{}
  json.Unmarshal(items.Body, &results)
  return results
}


func RecordMetrics() {
  buildMetrics()
	go func() {
		for {
      var results = queryZabbix()
      for resKey, result := range results {

        labelsWithValues := make(map[string]string)


        if len(labelsSliceAvg) > 0 {
          for _, vAvg := range labelsSliceAvg{
            if result[vAvg] != nil {
              labelsWithValues[vAvg] = result[vAvg].(string)
            }else{
              labelsWithValues[vAvg] = "NA"
            }
          }
        }

        if len(labelsSliceComplex) > 0 {
          for _, vCplx := range labelsSliceComplex{

            var promLabel = strings.Replace(vCplx, ">", "_", -1)
            var path = strings.Split(vCplx, ">")
            if result[path[0]] != nil {
              if(len(result[path[0]].([]interface{})) > 0){
                for _,cplx := range result[path[0]].([]interface{}) {
                  subCplx := cplx.(map[string]interface{})
                  if subCplx[path[1]] != nil {
      	             labelsWithValues[promLabel] = subCplx[path[1]].(string)
                  } else {
                     labelsWithValues[promLabel] = "NA"
                  }
      		      }
              } else {
                  labelsWithValues[promLabel] = "NA"
              }
            } else {
              labelsWithValues[promLabel] = "NA"
            }
          }
        }

        f, _ := strconv.ParseFloat(result[cnf.MetricValue].(string), 64)
        metricsSlice[resKey].With(labelsWithValues).Set(float64(f))
    	}

      time.Sleep(time.Duration(sourceRefreshSec) * time.Second)
		}
	}()
}
