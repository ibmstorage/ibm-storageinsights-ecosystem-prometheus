package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/shirou/gopsutil/mem"
	"none.rks/config/config"
	"none.rks/simetrics"
)

var (
	// Exporter Metrics
	memTotalDesc     = prometheus.NewDesc("memexporter_memory_total_bytes", "The amount of total memory in bytes.", []string{}, nil)
	memUsedDesc      = prometheus.NewDesc("memexporter_memory_used_bytes", "The amount of used memory in bytes.", []string{}, nil)
	memAvailableDesc = prometheus.NewDesc("memexporter_memory_available_bytes", "The amount of available memory in bytes.", []string{}, nil)
	memFreeDesc      = prometheus.NewDesc("memexporter_memory_free_bytes", "The amount of free memory in bytes.", []string{}, nil)

	// SI Metrics
	diskTotalDataRate     = prometheus.NewDesc("si_disk_total_data_rate", "Total Backend Data Rate.", []string{}, nil)
	diskTotalResponseTime = prometheus.NewDesc("si_disk_total_response_time", "Total Backend Response Time.", []string{}, nil)
)

type memCollector struct{}

func (mc memCollector) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(mc, ch)
}

func (mc memCollector) Collect(ch chan<- prometheus.Metric) {

	stats, err := mem.VirtualMemory()
	if err != nil {
		panic(err)
	}

	metricsList, err := simetrics.FetchData()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Data received: %v\n", metricsList)
	latest := metricsList[0].(map[string]interface{})

	ch <- prometheus.MustNewConstMetric(memTotalDesc, prometheus.GaugeValue, float64(stats.Total))
	ch <- prometheus.MustNewConstMetric(memUsedDesc, prometheus.GaugeValue, float64(stats.Used))
	ch <- prometheus.MustNewConstMetric(memAvailableDesc, prometheus.GaugeValue, float64(stats.Available))
	ch <- prometheus.MustNewConstMetric(memFreeDesc, prometheus.GaugeValue, float64(stats.Free))

	ch <- prometheus.MustNewConstMetric(diskTotalDataRate, prometheus.GaugeValue, float64(latest["disk_total_data_rate"].(float64)))
	ch <- prometheus.MustNewConstMetric(diskTotalResponseTime, prometheus.GaugeValue, float64(latest["disk_total_response_time"].(float64)))
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

		fmt.Println("System IDs:")
		for _, item := range config.AppConfig.SystemIds {
			fmt.Printf("- UUID: %s\n", item)
		}

		fmt.Println("System IDs:")
		for _, item := range config.AppConfig.Metrics {
			fmt.Printf("- Metric: %s\n", item)
		}
	}

	prometheus.MustRegister(&memCollector{})

	http.Handle("/metrics", promhttp.Handler())
	log.Println("Listening on", *listenAddr)
	log.Fatal(http.ListenAndServe(*listenAddr, nil))
}
