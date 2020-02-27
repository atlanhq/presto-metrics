package main

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
	"net/http"
)

const (
	namespace       = "atlan_presto"
	workerNamespace = namespace + "_worker"
	queryNamespace  = namespace + "_query"
)

var (
	stackNameVarLabel = []string{"prestoStackName"}
	workersVarLabel   = []string{"prestoStackName", "prestoWorkerId"}
	queryVarLabel     = []string{"prestoStackName", "queryId"}
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

	// memory

	clusterGeneralPoolFreeMemory = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "cluster_general_pool_free_memory"),
		"total free general pool memory of cluster.",
		stackNameVarLabel, nil,
	)

	clusterGeneralPoolTotalMemory = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "cluster_general_pool_total_memory"),
		"total general pool memory of cluster.",
		stackNameVarLabel, nil,
	)

	clusterGeneralPoolReservedMemory = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "cluster_general_pool_reserved_memory"),
		"total general pool reserved memory of cluster.",
		stackNameVarLabel, nil,
	)

	clusterGeneralPoolRevocableMemory = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "cluster_general_pool_revocable_memory"),
		"total general pool revocable memory of cluster.",
		stackNameVarLabel, nil,
	)

	medianWorkersGeneralPoolFreeMemory = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "median_workers_general_pool_memory"),
		"median workers general pool memory",
		stackNameVarLabel, nil,
	)

	meanWorkerGeneralFreePoolMemory = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "mean_workers_general_pool_memory"),
		"mean workers general pool memory",
		stackNameVarLabel, nil,
	)

	// cpu

	clusterUserCPUUtilisation = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "cluster_user_cpu_utilisation"),
		"cluster user cpu utlisation",
		stackNameVarLabel, nil,
	)

	clusterSystemCPUUtilisation = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "cluster_system_cpu_utilisation"),
		"cluster system cpu utlisation",
		stackNameVarLabel, nil,
	)

	medianWorkerUserCPUUtilisation = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "median_workers_user_cpu_utilisation"),
		"median workers user cpu utlisation",
		stackNameVarLabel, nil,
	)

	medianWorkerSystemCPUUtilisation = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "median_workers_system_cpu_utilisation"),
		"median workers system cpu utlisation",
		stackNameVarLabel, nil,
	)

	meanWorkerUserCPUUtilisation = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "mean_workers_user_cpu_utilisation"),
		"mean workers user cpu utlisation",
		stackNameVarLabel, nil,
	)

	meanWorkerSystemCPUUtilisation = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "mean_workers_system_cpu_utilisation"),
		"mean workers system cpu utlisation",
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

	// query level metrics

	completedDrivers = prometheus.NewDesc(
		prometheus.BuildFQName(queryNamespace, "", "completed_drivers"),
		"completed drivers",
		queryVarLabel, nil,
	)

	cumulativeUserMemory = prometheus.NewDesc(
		prometheus.BuildFQName(queryNamespace, "", "cumulative_user_memory"),
		"cumulative user memory",
		queryVarLabel, nil,
	)

	elapsedTime = prometheus.NewDesc(
		prometheus.BuildFQName(queryNamespace, "", "elapsed_time"),
		"elapsed time",
		queryVarLabel, nil,
	)

	executionTime = prometheus.NewDesc(
		prometheus.BuildFQName(queryNamespace, "", "execution_time"),
		"execution time",
		queryVarLabel, nil,
	)
	peakTotalMemoryReservation = prometheus.NewDesc(
		prometheus.BuildFQName(queryNamespace, "", "peak_total_memory_reservation"),
		"peak total memory reservation",
		queryVarLabel, nil,
	)

	peakUserMemoryReservation = prometheus.NewDesc(
		prometheus.BuildFQName(queryNamespace, "", "peak_user_memory_reservation"),
		"peak user memory reservation",
		queryVarLabel, nil,
	)

	queuedDrivers = prometheus.NewDesc(
		prometheus.BuildFQName(queryNamespace, "", "queued_drivers"),
		"queued drivers",
		queryVarLabel, nil,
	)

	queuedTime = prometheus.NewDesc(
		prometheus.BuildFQName(queryNamespace, "", "queued_time"),
		"queued time",
		queryVarLabel, nil,
	)

	rawInputDataSize = prometheus.NewDesc(
		prometheus.BuildFQName(queryNamespace, "", "raw_input_data_size"),
		"raw input data size",
		queryVarLabel, nil,
	)

	runningDriversCount = prometheus.NewDesc(
		prometheus.BuildFQName(queryNamespace, "", "running_drivers"),
		"running drivers",
		queryVarLabel, nil,
	)

	totalCpuTime = prometheus.NewDesc(
		prometheus.BuildFQName(queryNamespace, "", "total_cpu_time"),
		"total cpu time",
		queryVarLabel, nil,
	)

	totalDrivers = prometheus.NewDesc(
		prometheus.BuildFQName(queryNamespace, "", "total_drivers"),
		"total drivers",
		queryVarLabel, nil,
	)

	totalMemoryReservation = prometheus.NewDesc(
		prometheus.BuildFQName(queryNamespace, "", "total_memory_reservation"),
		"total memory reservation",
		queryVarLabel, nil,
	)

	totalScheduledTime = prometheus.NewDesc(
		prometheus.BuildFQName(queryNamespace, "", "total_scheduled_time"),
		"total scheduled time",
		queryVarLabel, nil,
	)

	userMemoryReservation = prometheus.NewDesc(
		prometheus.BuildFQName(queryNamespace, "", "user_memory_reservation"),
		"total user memory reservation",
		queryVarLabel, nil,
	)

	fullyBlocked = prometheus.NewDesc(
		prometheus.BuildFQName(queryNamespace, "", "fully_blocked"),
		"full blocked",
		queryVarLabel, nil,
	)
)

type clusterMetrics struct {
	config         Config
	ClusterMetrics ClusterMetrics
}

func (e *clusterMetrics) Describe(chan<- *prometheus.Desc) {

}

func (e *clusterMetrics) Collect(ch chan<- prometheus.Metric) {
	// cluster level metrics
	e.ClusterMetrics, _ = ClusterMetrics{}.collect(e.config.host, e.config.port, e.config.apiPrefix)
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

type workersMetrics struct {
	config               Config
	workerMetrics        workerMetrics
	ClusterMemoryMetrics ClusterMemoryMetrics
	ClusterCPUMetrics    ClusterCPUMetrics
}

func (e *workersMetrics) Describe(chan<- *prometheus.Desc) {

}

func (e *workersMetrics) Collect(ch chan<- prometheus.Metric) {
	// worker level metrics
	workers := workers{}
	workers, _ = workers.collect(e.config.host, e.config.port)

	var clusterGeneralPoolTotalMemoryArray []float64
	var clusterGeneralPoolFreeMemoryArray []float64
	var clusterGeneralPoolReservedMemoryArray []float64
	var clusterGeneralPoolRevocableMemoryArray []float64
	var clusterUserCPUUtilisationArray []float64
	var clusterSystemCPUUtilisationArray []float64

	for k := range workers {
		wm := workerMetrics{}
		workerId := k
		wm, err := wm.collect(e.config.host, e.config.port, workerId, e.config.apiPrefix)
		if err != nil {
			fmt.Println(err)
			return
		}

		clusterGeneralPoolFreeMemoryArray = append(clusterGeneralPoolFreeMemoryArray, float64(wm.MemoryInfo.Pools.General.FreeBytes))
		clusterGeneralPoolTotalMemoryArray = append(clusterGeneralPoolTotalMemoryArray, float64(wm.MemoryInfo.Pools.General.MaxBytes))
		clusterGeneralPoolReservedMemoryArray = append(clusterGeneralPoolReservedMemoryArray, float64(wm.MemoryInfo.Pools.General.ReservedBytes))
		clusterGeneralPoolRevocableMemoryArray = append(clusterGeneralPoolRevocableMemoryArray, float64(wm.MemoryInfo.Pools.General.ReservedRevocableBytes))
		clusterUserCPUUtilisationArray = append(clusterUserCPUUtilisationArray, wm.ProcessCpuLoad)
		clusterSystemCPUUtilisationArray = append(clusterSystemCPUUtilisationArray, wm.SystemCpuLoad)

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

	e.ClusterMemoryMetrics = ClusterMemoryMetrics{
		Sum(clusterGeneralPoolFreeMemoryArray),
		Sum(clusterGeneralPoolTotalMemoryArray),
		Sum(clusterGeneralPoolReservedMemoryArray),
		Sum(clusterGeneralPoolRevocableMemoryArray),
		Median(clusterGeneralPoolFreeMemoryArray),
		Sum(clusterGeneralPoolFreeMemoryArray) / float64(len(workers)),
	}
	e.ClusterCPUMetrics = ClusterCPUMetrics{
		ClusterUserCPUUtilisation:        Sum(clusterUserCPUUtilisationArray),
		ClusterSystemCPUUtilisation:      Sum(clusterSystemCPUUtilisationArray),
		MedianWorkerUserCPUUtilisation:   Median(clusterUserCPUUtilisationArray),
		MedianWorkerSystemCPUUtilisation: Median(clusterSystemCPUUtilisationArray),
		MeanWorkerUserCPUUtilisation:     Sum(clusterUserCPUUtilisationArray) / float64(len(workers)),
		MeanWorkerSystemCPUUtilisation:   Sum(clusterSystemCPUUtilisationArray) / float64(len(workers)),
	}

	ch <- prometheus.MustNewConstMetric(clusterGeneralPoolFreeMemory, prometheus.GaugeValue, e.ClusterMemoryMetrics.ClusterGeneralPoolFreeMemory, e.config.stackName)
	ch <- prometheus.MustNewConstMetric(clusterGeneralPoolTotalMemory, prometheus.GaugeValue, e.ClusterMemoryMetrics.ClusterGeneralPoolTotalMemory, e.config.stackName)
	ch <- prometheus.MustNewConstMetric(clusterGeneralPoolReservedMemory, prometheus.GaugeValue, e.ClusterMemoryMetrics.ClusterGeneralPoolReservedMemory, e.config.stackName)
	ch <- prometheus.MustNewConstMetric(clusterGeneralPoolRevocableMemory, prometheus.GaugeValue, e.ClusterMemoryMetrics.ClusterGeneralPoolRevocableMemory, e.config.stackName)
	ch <- prometheus.MustNewConstMetric(medianWorkersGeneralPoolFreeMemory, prometheus.GaugeValue, e.ClusterMemoryMetrics.MedianWorkersGeneralPoolFreeMemory, e.config.stackName)
	ch <- prometheus.MustNewConstMetric(meanWorkerGeneralFreePoolMemory, prometheus.GaugeValue, e.ClusterMemoryMetrics.MeanWorkerGeneralFreePoolMemory, e.config.stackName)
	ch <- prometheus.MustNewConstMetric(clusterUserCPUUtilisation, prometheus.GaugeValue, e.ClusterCPUMetrics.ClusterUserCPUUtilisation, e.config.stackName)
	ch <- prometheus.MustNewConstMetric(clusterSystemCPUUtilisation, prometheus.GaugeValue, e.ClusterCPUMetrics.ClusterSystemCPUUtilisation, e.config.stackName)
	ch <- prometheus.MustNewConstMetric(medianWorkerUserCPUUtilisation, prometheus.GaugeValue, e.ClusterCPUMetrics.MedianWorkerUserCPUUtilisation, e.config.stackName)
	ch <- prometheus.MustNewConstMetric(medianWorkerSystemCPUUtilisation, prometheus.GaugeValue, e.ClusterCPUMetrics.MedianWorkerSystemCPUUtilisation, e.config.stackName)
	ch <- prometheus.MustNewConstMetric(meanWorkerUserCPUUtilisation, prometheus.GaugeValue, e.ClusterCPUMetrics.MeanWorkerUserCPUUtilisation, e.config.stackName)
	ch <- prometheus.MustNewConstMetric(meanWorkerSystemCPUUtilisation, prometheus.GaugeValue, e.ClusterCPUMetrics.MeanWorkerSystemCPUUtilisation, e.config.stackName)

}

type queryMetrics struct {
	config         Config
	clusterQueries clusterQueries
}

func (e *queryMetrics) Describe(chan<- *prometheus.Desc) {

}

func (e *queryMetrics) Collect(ch chan<- prometheus.Metric) {
	e.clusterQueries, _ = clusterQueries{}.collect(e.config.host, e.config.port, e.config.apiPrefix)
	for _, query := range e.clusterQueries {
		queryStats := query.QueryStats
		ch <- prometheus.MustNewConstMetric(completedDrivers, prometheus.GaugeValue, float64(queryStats.CompletedDrivers), e.config.stackName, query.QueryId)
		ch <- prometheus.MustNewConstMetric(cumulativeUserMemory, prometheus.GaugeValue, queryStats.CumulativeUserMemory, e.config.stackName, query.QueryId)
		ch <- prometheus.MustNewConstMetric(elapsedTime, prometheus.GaugeValue, queryStats.ElapsedTimeParsed, e.config.stackName, query.QueryId)
		ch <- prometheus.MustNewConstMetric(executionTime, prometheus.GaugeValue, queryStats.ExecutionTimeParsed, e.config.stackName, query.QueryId)
		ch <- prometheus.MustNewConstMetric(peakTotalMemoryReservation, prometheus.GaugeValue, queryStats.PeakTotalMemoryReservationParsed, e.config.stackName, query.QueryId)
		ch <- prometheus.MustNewConstMetric(peakUserMemoryReservation, prometheus.GaugeValue, queryStats.PeakUserMemoryReservationParsed, e.config.stackName, query.QueryId)
		ch <- prometheus.MustNewConstMetric(queuedDrivers, prometheus.GaugeValue, float64(queryStats.QueuedDrivers), e.config.stackName, query.QueryId)
		ch <- prometheus.MustNewConstMetric(queuedTime, prometheus.GaugeValue, queryStats.QueuedTimeParsed, e.config.stackName, query.QueryId)
		ch <- prometheus.MustNewConstMetric(rawInputDataSize, prometheus.GaugeValue, queryStats.RawInputDataSizeParsed, e.config.stackName, query.QueryId)
		ch <- prometheus.MustNewConstMetric(runningDriversCount, prometheus.GaugeValue, float64(queryStats.RunningDrivers), e.config.stackName, query.QueryId)
		ch <- prometheus.MustNewConstMetric(totalCpuTime, prometheus.GaugeValue, queryStats.TotalCpuTimeParsed, e.config.stackName, query.QueryId)
		ch <- prometheus.MustNewConstMetric(totalDrivers, prometheus.GaugeValue, float64(queryStats.TotalDrivers), e.config.stackName, query.QueryId)
		ch <- prometheus.MustNewConstMetric(totalMemoryReservation, prometheus.GaugeValue, queryStats.TotalMemoryReservationParsed, e.config.stackName, query.QueryId)
		ch <- prometheus.MustNewConstMetric(totalScheduledTime, prometheus.GaugeValue, queryStats.TotalScheduledTimeParsed, e.config.stackName, query.QueryId)
		ch <- prometheus.MustNewConstMetric(userMemoryReservation, prometheus.GaugeValue, queryStats.UserMemoryReservationParsed, e.config.stackName, query.QueryId)
		ch <- prometheus.MustNewConstMetric(fullyBlocked, prometheus.GaugeValue, float64(boolToInt(queryStats.FullyBlocked)), e.config.stackName, query.QueryId)
	}
}

func prometheusExporterStart(host string, port string, stackName string, listenAddress string, apiPrefix string) {
	http.Handle("/metrics", promhttp.Handler())
	log.Info("Listening on ", listenAddress)
	config := Config{
		host:      host,
		port:      port,
		stackName: stackName,
		apiPrefix: apiPrefix,
	}
	prometheus.MustRegister(&clusterMetrics{
		config: config,
	})
	prometheus.MustRegister(&workersMetrics{
		config: config,
	})
	prometheus.MustRegister(&queryMetrics{
		config: config,
	})
	log.Fatal(http.ListenAndServe(listenAddress, nil))
}
