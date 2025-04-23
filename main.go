package main

import (
	"net/http"

	"k8s.io/klog/v2"
)

const (
	filterPrefix    = "/scheduler/filter"
	scorePrefix     = "/scheduler/score"
	healthCheckPath = "/healthz"
)

func main() {
	http.HandleFunc(healthCheckPath, healthCheckHandler)
	http.HandleFunc(filterPrefix, filterHandler)
	http.HandleFunc(scorePrefix, scoreHandler)

	klog.Info("Starting scheduler extender on :8010")
	if err := http.ListenAndServe(":8010", nil); err != nil {
		klog.Fatalf("Failed to start server: %v", err)
	}
}
