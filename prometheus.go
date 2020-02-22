package main

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
	"net/http"
	"strings"
)

const (
	namespace       = "atlan_presto"
	workerNamespace = namespace + "_worker"
)

var (
	stackNameVarLabel = []string{"prestoStackName"}
	workersVarLabel   = []string{"prestoStackName", "prestoWorkerId"}
)

var (
	// cluster level metrics
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

	// worker level metrics
	processors = prometheus.NewDesc(
		prometheus.BuildFQName(workerNamespace, "", "processors"),
		"total processors a worker has",
		workersVarLabel, nil,
	)
	heapAvailable = prometheus.NewDesc(
		prometheus.BuildFQName(workerNamespace, "", "heap_available"),
		"heap available for worker",
		workersVarLabel, nil,
	)
	heapUsed = prometheus.NewDesc(
		prometheus.BuildFQName(workerNamespace, "", "heap_used"),
		"heap used for worker",
		workersVarLabel, nil,
	)
	nonHeapUsed = prometheus.NewDesc(
		prometheus.BuildFQName(workerNamespace, "", "non_heap_used"),
		"non heap used for worker",
		workersVarLabel, nil,
	)
	processCpuLoad = prometheus.NewDesc(
		prometheus.BuildFQName(workerNamespace, "", "process_cpu_load"),
		"user cpu load",
		workersVarLabel, nil,
	)
	systemCpuLoad = prometheus.NewDesc(
		prometheus.BuildFQName(workerNamespace, "", "system_cpu_load"),
		"system cpu load",
		workersVarLabel, nil,
	)
	generalPoolFreeBytes = prometheus.NewDesc(
		prometheus.BuildFQName(workerNamespace, "", "general_pool_free_bytes"),
		"general pool free bytes",
		workersVarLabel, nil,
	)
	generalPoolMaxBytes = prometheus.NewDesc(
		prometheus.BuildFQName(workerNamespace, "", "general_pool_max_bytes"),
		"general pool max bytes",
		workersVarLabel, nil,
	)
	generalPoolReservedBytes = prometheus.NewDesc(
		prometheus.BuildFQName(workerNamespace, "", "general_pool_reserved_bytes"),
		"general pool reserved bytes",
		workersVarLabel, nil,
	)
	generalPoolReservedRevocableBytes = prometheus.NewDesc(
		prometheus.BuildFQName(workerNamespace, "", "general_pool_reserved_revocable_bytes"),
		"general pool reserved revocable bytes",
		workersVarLabel, nil,
	)
	totalNodeMemory = prometheus.NewDesc(
		prometheus.BuildFQName(workerNamespace, "", "total_node_memory"),
		"total node memory",
		workersVarLabel, nil,
	)
)

type Config struct {
	host      string
	port      string
	stackName string
}

type clusterMetrics struct {
	config         Config
	ClusterMetrics ClusterMetrics
}

func (e *clusterMetrics) Describe(chan<- *prometheus.Desc) {

}

type workersMetrics struct {
	config        Config
	workerMetrics workerMetrics
}

func (e *clusterMetrics) Collect(ch chan<- prometheus.Metric) {
	// cluster level metrics
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

func (e *workersMetrics) Describe(chan<- *prometheus.Desc) {

}

func (e *workersMetrics) Collect(ch chan<- prometheus.Metric) {
	// worker level metrics
	workers := workers{}
	workers, _ = workers.collect(e.config.host, e.config.port)
	for k, _ := range workers {
		wm := workerMetrics{}
		workerId := strings.Split(k, " ")[0]
		wm, err := wm.collect(e.config.host, e.config.port, workerId)
		if err != nil {
			fmt.Println(err)
			return
		}
		ch <- prometheus.MustNewConstMetric(processors, prometheus.GaugeValue, float64(wm.Processors), e.config.stackName, workerId)
		ch <- prometheus.MustNewConstMetric(heapAvailable, prometheus.GaugeValue, float64(wm.HeapAvailable), e.config.stackName, workerId)
		ch <- prometheus.MustNewConstMetric(heapUsed, prometheus.GaugeValue, float64(wm.HeapUsed), e.config.stackName, workerId)
		ch <- prometheus.MustNewConstMetric(nonHeapUsed, prometheus.GaugeValue, float64(wm.NonHeapUsed), e.config.stackName, workerId)
		ch <- prometheus.MustNewConstMetric(processCpuLoad, prometheus.GaugeValue, wm.ProcessCpuLoad, e.config.stackName, workerId)
		ch <- prometheus.MustNewConstMetric(systemCpuLoad, prometheus.GaugeValue, wm.SystemCpuLoad, e.config.stackName, workerId)
		ch <- prometheus.MustNewConstMetric(generalPoolFreeBytes, prometheus.GaugeValue, float64(wm.MemoryInfo.Pools.General.FreeBytes), e.config.stackName, workerId)
		ch <- prometheus.MustNewConstMetric(generalPoolMaxBytes, prometheus.GaugeValue, float64(wm.MemoryInfo.Pools.General.MaxBytes), e.config.stackName, workerId)
		ch <- prometheus.MustNewConstMetric(generalPoolReservedBytes, prometheus.GaugeValue, float64(wm.MemoryInfo.Pools.General.ReservedBytes), e.config.stackName, workerId)
		ch <- prometheus.MustNewConstMetric(generalPoolReservedRevocableBytes, prometheus.GaugeValue, float64(wm.MemoryInfo.Pools.General.ReservedRevocableBytes), e.config.stackName, workerId)
		ch <- prometheus.MustNewConstMetric(totalNodeMemory, prometheus.GaugeValue, float64(wm.TotalNodeMemory), e.config.stackName, workerId)
	}
}

func prometheusExporterStart(host string, port string, stackName string, listenAddress string) {
	http.Handle("/metrics", promhttp.Handler())
	log.Info("Listening on ", listenAddress)
	config := Config{
		host:      host,
		port:      port,
		stackName: stackName,
	}
	prometheus.MustRegister(&clusterMetrics{
		config: config,
	})
	prometheus.MustRegister(&workersMetrics{
		config:        config,
		workerMetrics: workerMetrics{},
	})
	log.Fatal(http.ListenAndServe(listenAddress, nil))
}
