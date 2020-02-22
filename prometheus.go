package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
)

const (
	namespace = "atlan_presto"
)

var (
	stackNameVarLabel = []string{"prestoStackName"}
)

var (
	runningQueries = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "running_queries"),
		"Running requests of the presto cluster.",
		stackNameVarLabel, nil,
	)
	blockedQueries = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "blocked_queries"),
		"Blocked queries of the presto cluster.",
		stackNameVarLabel, nil,
	)
	queuedQueries = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "queued_queries"),
		"Queued queries of the presto cluster.",
		stackNameVarLabel, nil,
	)
	activeWorkers = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "active_workers"),
		"Active workers of the presto cluster.",
		stackNameVarLabel, nil,
	)
	runningDrivers = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "running_drivers"),
		"Running drivers of the presto cluster.",
		stackNameVarLabel, nil,
	)
	reservedMemory = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "reserved_memory"),
		"Reserved memory of the presto cluster.",
		stackNameVarLabel, nil,
	)
	totalInputRows = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "total_input_rows"),
		"Total input rows of the presto cluster.",
		stackNameVarLabel, nil,
	)
	totalInputBytes = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "total_input_bytes"),
		"Total input bytes of the presto cluster.",
		stackNameVarLabel, nil,
	)
	totalCpuTimeSecs = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "total_cpu_time_secs"),
		"Total cpu time of the presto cluster.",
		stackNameVarLabel, nil,
	)
)

type Config struct {
	host      string
	port      string
	stackName string
}

type Metrics struct {
	config         Config
	ClusterMetrics ClusterMetrics
}

func (e Metrics) Describe(ch chan<- *prometheus.Desc) {
	ch <- runningQueries
}

func (e *Metrics) Collect(ch chan<- prometheus.Metric) {
	e.ClusterMetrics, _ = ClusterMetrics{}.collect(e.config.host, e.config.port)
	ch <- prometheus.MustNewConstMetric(runningQueries, prometheus.GaugeValue, e.ClusterMetrics.RunningQueries, e.config.stackName)
	ch <- prometheus.MustNewConstMetric(blockedQueries, prometheus.GaugeValue, e.ClusterMetrics.BlockedQueries, e.config.stackName)
	ch <- prometheus.MustNewConstMetric(queuedQueries, prometheus.GaugeValue, e.ClusterMetrics.QueuedQueries, e.config.stackName)
	ch <- prometheus.MustNewConstMetric(activeWorkers, prometheus.GaugeValue, e.ClusterMetrics.ActiveWorkers, e.config.stackName)
	ch <- prometheus.MustNewConstMetric(runningDrivers, prometheus.GaugeValue, e.ClusterMetrics.RunningQueries, e.config.stackName)
	ch <- prometheus.MustNewConstMetric(reservedMemory, prometheus.GaugeValue, e.ClusterMetrics.ReservedMemory, e.config.stackName)
	ch <- prometheus.MustNewConstMetric(totalInputRows, prometheus.GaugeValue, e.ClusterMetrics.TotalInputRows, e.config.stackName)
	ch <- prometheus.MustNewConstMetric(totalInputBytes, prometheus.GaugeValue, e.ClusterMetrics.TotalInputBytes, e.config.stackName)
	ch <- prometheus.MustNewConstMetric(totalCpuTimeSecs, prometheus.GaugeValue, e.ClusterMetrics.TotalCpuTimeSecs, e.config.stackName)
}

func prometheusExporterStart(host string, port string, stackName string, listenAddress string) {
	http.Handle("/metrics", promhttp.Handler())
	log.Info("Listening on ", listenAddress)
	prometheus.MustRegister(&Metrics{
		config: struct {
			host      string
			port      string
			stackName string
		}{host: host, port: port, stackName: stackName},
	})
	log.Fatal(http.ListenAndServe(listenAddress, nil))
}
