package handlers

import (
	"encoding/json"
	"github.com/prometheus/client_golang/prometheus"
	cnf "github.com/rzrbld/zabbix-exporter-3000/config"
	zbx "github.com/rzrbld/zabbix-exporter-3000/zabbix"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var sourceRefreshSec, _ = strconv.Atoi(cnf.SourceRefresh)
var labelsSliceRaw = strings.Split(cnf.MetricLabels, ",")

// labels in prom format
var labelsSlicePrometheus []string

// labels with path "a>b"
var labelsSliceComplex []string

// labels average
var labelsSliceAvg []string
var rawMetricNames []string
var uniqMetricNames []string
var rawMetricDesc []string
var uniqMetricDesc []string
var itemsMetric *prometheus.GaugeVec
var metricsMap = make(map[string]*prometheus.GaugeVec, 1000)

//helpers
func cleanUpName(name string) string {
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		log.Fatal(err)
	}
	cleanName := reg.ReplaceAllString(strings.ToLower(name), "")
	return cleanName
}

func uniqueSlice(intSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func registerMetric(metric *prometheus.GaugeVec) {
	if cnf.StrictRegister {
		prometheus.MustRegister(metric)
	} else {
		prometheus.Register(metric)
	}
}

func buildMetrics() {
	for _, vl := range labelsSliceRaw {
		if strings.Contains(vl, ">") {
			labelsSlicePrometheus = append(labelsSlicePrometheus, strings.Replace(vl, ">", "_", -1))
			labelsSliceComplex = append(labelsSliceComplex, vl)
		} else {
			labelsSlicePrometheus = append(labelsSlicePrometheus, vl)
			labelsSliceAvg = append(labelsSliceAvg, vl)
		}
	}

	log.Print("Labels that will be produced      :", labelsSlicePrometheus)
	log.Print("Complex labels that will be parsed:", labelsSliceComplex)
	log.Print("Plain labels that will be parsed  :", labelsSliceAvg)

	var results = queryZabbix()

	if cnf.MetricNameField != "" {
		for k, result := range results {
			cleanName := cleanUpName(result[cnf.MetricNameField].(string))
			rawMetricNames = append(rawMetricNames, cleanName)
			if result[cnf.MetricHelpField] != nil {
				if result[cnf.MetricHelpField].(string) != "" {
					rawMetricDesc = append(rawMetricDesc, result[cnf.MetricHelpField].(string))
				} else {
					rawMetricDesc = append(rawMetricDesc, "NA_"+strconv.Itoa(k))
				}
			} else {
				rawMetricDesc = append(rawMetricDesc, "NA")
			}
		}
	}

	log.Println("Raw Metrics    : ", rawMetricNames)
	log.Println("Raw Description: ", rawMetricDesc)
	uniqMetricNames := uniqueSlice(rawMetricNames)
	uniqMetricDesc := uniqueSlice(rawMetricDesc)

	log.Println("Uniq Metrics    : ", uniqMetricNames)
	log.Println("Uniq Description: ", uniqMetricDesc)

	if len(uniqMetricNames) != len(uniqMetricDesc) {
		log.Print("WARNING: Number of Metrics and Description not equal")

		if len(uniqMetricNames) < len(uniqMetricDesc) {
			log.Fatal("ERROR: Insufficient uniq Metrics. Try to use more unique ZE3000_METRIC_NAME_FIELD, or use ZE3000_SINGLE_METRIC_NAME=true")
		} else {
			log.Print("WARNING: I try to heal this by populating NA")
			for k, _ := range uniqMetricNames {
				uniqMetricDesc = append(uniqMetricDesc, "NA_"+strconv.Itoa(k))
			}
		}
	}

	if cnf.SingleMetric {
		fullName := cnf.MetricNamePrefix

		itemsMetric = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: cnf.MetricNamespace,
			Subsystem: cnf.MetricSubsystem,
			Name:      fullName,
			Help:      cnf.SingleMetricHelp,
		}, labelsSlicePrometheus)

		registerMetric(itemsMetric)
	} else {
		for k, name := range uniqMetricNames {
			//clean up metric name
			fullName := cnf.MetricNamePrefix
			if cnf.MetricNameField != "" {
				fullName = cnf.MetricNamePrefix + "_" + name
			}

			itemsMetric = prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Namespace: cnf.MetricNamespace,
				Subsystem: cnf.MetricSubsystem,
				Name:      fullName,
				Help:      uniqMetricDesc[k],
			}, labelsSlicePrometheus)

			metricsMap[name] = itemsMetric
		}

		for _, metric := range metricsMap {
			registerMetric(metric)
		}
	}

	log.Print("Number of bject getting from Zabbix    : ", len(results))
	if cnf.SingleMetric {
		log.Print("Number of metrics that will be produced: ", 1)
	} else {
		log.Print("Number of metrics that will be produced: ", len(metricsMap))
	}
}

func queryZabbix() []map[string]interface{} {
	items, err := zbx.Session.Do(zbx.Query)
	if err != nil {
		log.Fatal("ERROR While Do request: ", err)
	}

	var results []map[string]interface{}
	json.Unmarshal(items.Body, &results)
	if len(results) == 0 {
		log.Fatal("Empty response from Zabbix. Check query at ZE3000_ZABBIX_QUERY")
	}
	return results
}

func RecordMetrics() {
	buildMetrics()
	go func() {
		for {
			var results = queryZabbix()
			for _, result := range results {

				labelsWithValues := make(map[string]string)

				if len(labelsSliceAvg) > 0 {
					for _, vAvg := range labelsSliceAvg {
						if result[vAvg] != nil {
							labelsWithValues[vAvg] = result[vAvg].(string)
						} else {
							labelsWithValues[vAvg] = "NA"
						}
					}
				}

				if len(labelsSliceComplex) > 0 {
					for _, vCplx := range labelsSliceComplex {

						var promLabel = strings.Replace(vCplx, ">", "_", -1)
						var path = strings.Split(vCplx, ">")
						if result[path[0]] != nil {
							if len(result[path[0]].([]interface{})) > 0 {
								for _, cplx := range result[path[0]].([]interface{}) {
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

				var f float64
				f = float64(0)
				if result[cnf.MetricValue] != nil {
					f, _ = strconv.ParseFloat(result[cnf.MetricValue].(string), 64)
				}

				if cnf.SingleMetric {
					itemsMetric.With(labelsWithValues).Set(f)
				} else {
					cleanName := cleanUpName(result[cnf.MetricNameField].(string))
					metricsMap[cleanName].With(labelsWithValues).Set(f)
				}
			}

			time.Sleep(time.Duration(sourceRefreshSec) * time.Second)
		}
	}()
}
