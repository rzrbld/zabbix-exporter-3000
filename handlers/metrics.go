package handlers

import (
	// "log"
	// "time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// var bucketsSizes = promauto.NewGaugeVec(prometheus.GaugeOpts{
// 	Name: "bucket_size_current",
// 	Help: "bucket size in kbytes",
// }, []string{"bucket"})

var objectsCount = promauto.NewGauge(prometheus.GaugeOpts{
	Name: "objects_count_current",
	Help: "number of objects on cluster",
})

func RecordMetrics() {
	go func() {
		for {
			// du, err := madmClnt.DataUsageInfo()
			// if err != nil {
			// 	log.Print("Error while getting bucket size metrics from server")
			// } else {
			// 	if len(du.BucketsSizes) != 0 {
			// 		for k, v := range du.BucketsSizes {
			// 			bucketsSizes.WithLabelValues(string(k)).Set(float64(v))
			// 		}
			// 	}
			// 	objectsCount.Set(float64(du.ObjectsCount))
			// objectsSize.Set(float64(1111))
			objectsCount.Set(float64(2222))
			// }
			// time.Sleep(2 * time.Minute)
		}
	}()
}
