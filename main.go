package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func getInstanceName() string {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://metadata.google.internal/computeMetadata/v1/instance/name", nil)
	if err != nil {
		log.Fatal(err)
		return "unknown"
	}
	req.Header.Set("Metadata-Flavor", "Google")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
		return "unknown"
	}
	if resp.StatusCode == 200 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
			return "unknown"
		}
		return string(body)
	}
	fmt.Printf("failed to get instance name got %s", resp.Status)
	return "unknown"
}

func getMaintenanceState() float64 {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://metadata.google.internal/computeMetadata/v1/instance/maintenance-event", nil)
	if err != nil {
		return -1
	}
	req.Header.Set("Metadata-Flavor", "Google")

	// loop to retry if there is API maintenance
	for i := 1; i < 10; i++ {
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
			return -1
		}
		if resp.StatusCode == 200 {
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println(err)
				return -1
			}
			switch string(body) {
			case "NONE":
				return 0
			case "MIGRATE_ON_HOST_MAINTENANCE":
				return 1
			case "TERMINATE_ON_HOST_MAINTENANCE":
				return 2
			}
		} else {
			fmt.Printf("Not status code 200 got %s", resp.Status)
		}
		time.Sleep(10 * time.Second)
	}
	return -1
}

func recordMetrics() {
	go func() {
		for {
			opsProcessed.Set(getMaintenanceState())
			time.Sleep(10 * time.Second)
		}
	}()
}

var (
	opsProcessed = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "gcp_maintenance_state",
		Help: "Report if a maintenance is planned on a GCE instance.",
		ConstLabels: prometheus.Labels{
			"host": getInstanceName(),
		},
	})
)

func main() {
	recordMetrics()

	registry := prometheus.NewRegistry()
	registry.MustRegister(opsProcessed)
	handler := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})

	http.Handle("/metrics", handler)
	http.ListenAndServe(":9723", nil)
}
