package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/ibmstorage/ibm-storageinsights-ecosystem-prometheus/config"
	"github.com/ibmstorage/ibm-storageinsights-ecosystem-prometheus/simetrics"
)

type MetricCollector struct {
	metrics map[string]*prometheus.Desc
}

func NewMetricCollector() *MetricCollector {
	metrics := make(map[string]*prometheus.Desc)

	// Dynamically create Prometheus metrics from the list of names configured
	for _, name := range config.AppConfig.Metrics {
		metrics[name] = prometheus.NewDesc(
			"si_" + name,               			// Metric name
			"Storage Insights metric " + name, 		// Description
			[]string{"device_name"},     			// Labels (if applicable)
			nil,                 					// Optional: No extra labels or values
		)
	}

	return &MetricCollector{
		metrics: metrics,
	}
}

// Describe implements prometheus.Collector interface
func (collector *MetricCollector) Describe(ch chan<- *prometheus.Desc) {
	// Send all metric descriptions to the channel
	for _, metric := range collector.metrics {
		ch <- metric
	}
}

// Collect implements prometheus.Collector interface
func (collector *MetricCollector) Collect(ch chan<- prometheus.Metric) {
	metricsList, err := simetrics.FetchData()
	if err != nil {
		panic(err)
	}
	
	fmt.Printf("Data received-->: %v\n", metricsList)

	// Loop through each entry in the metrics list
	for _, item := range metricsList {
		if data, ok := item.(map[string]interface{}); ok {

			// Get the name and metrics
			compname := data["name"]
			metrics := data["metrics"].([]interface{}) // Get the metrics array

			if len(metrics) > 0 {
				// Extract the latest metric (first element in the metrics list)
				latestMetric := metrics[0].(map[string]interface{})

				// Type assert compname to string before using it
				compnameStr, ok := compname.(string)
				if !ok {
					log.Println("compname is not a string")
					continue
				}

				// Loop through each metric name in the collector's metrics map
				for name, metric := range collector.metrics {
					// Ensure the metric exists and is of the correct type (float64)
					if val, exists := latestMetric[name]; exists {
						if floatVal, ok := val.(float64); ok {
							ch <- prometheus.MustNewConstMetric(
								metric,                // The Prometheus metric descriptor
								prometheus.GaugeValue, // The type of metric
								floatVal,              // The value of the metric
								compnameStr,              // The "name" label for the metric
							)
							log.Printf("Collected metric: %s", name)
						} else {
							log.Printf("Invalid type for metric %s: expected float64, got %T", name, val)
						}
					} else {
						log.Printf("Metric %s not found in the latest data", name)
					}
				}
			} else {
				log.Println("No metrics found for", compname)
			}
		}
	}	
}

func main() {
	listenAddr := flag.String("listen-address", ":8085", "The address to listen on for metrics requests.")
	configFilePath := flag.String("config", "config.json", "Path to the configuration file")
	flag.Parse()

	configBytes, err := os.ReadFile(*configFilePath)
	if err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}

	if err := config.LoadConfig(configBytes); err != nil {
		log.Fatalf("Failed to parse config file: %v", err)
	}

	if config.AppConfig.Debug {
		fmt.Printf("SI URL %s\n", config.AppConfig.Siurl)
		fmt.Printf("IBM ID %s\n", config.AppConfig.Ibmid)
		fmt.Printf("ApiKey: %s\n", config.AppConfig.ApiKey)
		fmt.Printf("Tenant ID: %s\n", config.AppConfig.TenantId)

		fmt.Println("Metric types")
		for _, item := range config.AppConfig.Metrics {
			fmt.Printf("- Metric: %s\n", item)
		}
	}

	collector := NewMetricCollector()
	prometheus.MustRegister(collector)

	http.Handle("/metrics", promhttp.Handler())
	log.Println("Listening on", *listenAddr)
	log.Fatal(http.ListenAndServe(*listenAddr, nil))
}
