package main

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
)

func prometheusExporterStart(listenAddress string) {
	http.Handle("/metrics", promhttp.Handler())
	log.Info("Listening on ", listenAddress)
	log.Fatal(http.ListenAndServe(listenAddress, nil))
}
