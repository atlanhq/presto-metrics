package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/prestodb/presto-go-client/presto"
	"github.com/prometheus/common/log"
	"io/ioutil"
	"net/http"
	"strings"
)

type worker struct {
	NodeId        string
	HttpUri       string
	NodeVersion   string
	IsCoordinator string
	State         string
}

type workers map[string]worker

type ClusterMemoryMetrics struct {
	ClusterGeneralPoolFreeMemory       float64
	ClusterGeneralPoolTotalMemory      float64
	ClusterGeneralPoolReservedMemory   float64
	ClusterGeneralPoolRevocableMemory  float64
	MedianWorkersGeneralPoolFreeMemory float64
	MeanWorkerGeneralFreePoolMemory    float64
}

type ClusterCPUMetrics struct {
	ClusterUserCPUUtilisation        float64
	ClusterSystemCPUUtilisation      float64
	MedianWorkerUserCPUUtilisation   float64
	MedianWorkerSystemCPUUtilisation float64
	MeanWorkerUserCPUUtilisation     float64
	MeanWorkerSystemCPUUtilisation   float64
}

func (w workers) collect(host string, port string) (workers, error) {
	db, err := sql.Open("presto",
		fmt.Sprintf("http://presto_metrics@%s:%s?catalog=system&schema=runtime", host, port))
	if err != nil {
		return nil, err
	}
	data, err := db.Query("select * from nodes")
	if err != nil {
		return nil, err
	}
	for data.Next() {
		worker := worker{}
		err := data.Scan(&worker.NodeId, &worker.HttpUri, &worker.NodeVersion, &worker.IsCoordinator, &worker.State)
		if err != nil {
			fmt.Println("Error ", err)
			continue
		}
		w[worker.NodeId] = worker
	}
	return w, nil
}

type workerMetrics struct {
	Processors      int64   `json:"processors"`
	HeapAvailable   int64   `json:"heapAvailable"`
	HeapUsed        int64   `json:"heapUsed"`
	NonHeapUsed     int64   `json:"nonHeapUsed"`
	ProcessCpuLoad  float64 `json:"processCpuLoad"`
	SystemCpuLoad   float64 `json:"systemCpuLoad"`
	Uptime          string  `json:"uptime"`
	InternalAddress string  `json:"internalAddress"`
	NodeId          string  `json:"nodeId"`
	NodeVersion     struct {
		Version string `json:"version"`
	} `json:"nodeVersion"`
	MemoryInfo struct {
		Pools struct {
			General struct {
				FreeBytes                int64 `json:"freeBytes"`
				MaxBytes                 int64 `json:"maxBytes"`
				ReservedBytes            int64 `json:"reservedBytes"`
				ReservedRevocableBytes   int64 `json:"reservedRevocableBytes"`
				QueryMemoryAllocationMap map[string][]struct {
					Allocation int64  `json:"allocation"`
					Tag        string `json:"tag"`
				} `json:"queryMemoryAllocations"`
				QueryMemoryReservations map[string]int64 `json:"queryMemoryReservations"`
			} `json:"general"`
			Reserved struct {
				FreeBytes              int64 `json:"freeBytes"`
				MaxBytes               int64 `json:"maxBytes"`
				ReservedBytes          int64 `json:"reservedBytes"`
				ReservedRevocableBytes int64 `json:"reservedRevocableBytes"`
			} `json:"reserved"`
			TotalNodeMemory int64 `json:"totalNodeMemory"`
		} `json:"pools"`
	} `json:"memoryInfo"`
	TotalNodeMemory int64 `json:"totalNodeMemory"`
}

func (wm workerMetrics) collect(host string, port string, nodeId string, apiPrefix string) (workerMetrics, error) {
	resp, err := http.Get("http://" + host + ":" + port + "/" + apiPrefix + "worker/" + nodeId + "/status")
	if err != nil {
		log.Errorf("%s", err)
		return workerMetrics{}, err
	}
	if resp.StatusCode != 200 {
		log.Errorf("%s", resp.StatusCode)
		log.Errorf("%s", err)
		log.Errorf(nodeId)
		return workerMetrics{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("%s", err)
		return workerMetrics{}, err
	}

	_ = json.Unmarshal(body, &wm)
	return wm, nil
}

type clusterQuery struct {
	MemoryPool string `json:"memoryPool"`
	Query      string `json:"query"`
	QueryId    string `json:"queryId"`
	QueryStats struct {
		CompletedDrivers                 int64   `json:"completedDrivers"`
		CreateTime                       string  `json:"createTime"`
		CumulativeUserMemory             float64 `json:"cumulativeUserMemory"`
		ElapsedTime                      string  `json:"elapsedTime"`
		ElapsedTimeParsed                float64
		ExecutionTime                    string `json:"executionTime"`
		ExecutionTimeParsed              float64
		FullyBlocked                     bool   `json:"fullyBlocked"`
		PeakTotalMemoryReservation       string `json:"peakTotalMemoryReservation"`
		PeakTotalMemoryReservationParsed float64
		PeakUserMemoryReservation        string `json:"peakUserMemoryReservation"`
		PeakUserMemoryReservationParsed  float64
		QueuedDrivers                    int64  `json:"queuedDrivers"`
		QueuedTime                       string `json:"queuedTime"`
		QueuedTimeParsed                 float64
		RawInputDataSize                 string `json:"rawInputDataSize"`
		RawInputDataSizeParsed           float64
		RawInputPositions                string `json:"rawInputPositions"`
		RunningDrivers                   int64  `json:"runningDrivers"`
		TotalCpuTime                     string `json:"totalCpuTime"`
		TotalCpuTimeParsed               float64
		TotalDrivers                     int64  `json:"totalDrivers"`
		TotalMemoryReservation           string `json:"totalMemoryReservation"`
		TotalMemoryReservationParsed     float64
		TotalScheduledTime               string `json:"totalScheduledTime"`
		TotalScheduledTimeParsed         float64
		UserMemoryReservation            string `json:"userMemoryReservation"`
		UserMemoryReservationParsed      float64
	} `json:"queryStats"`
	State string `json:"state"`
}

type clusterQueries []clusterQuery

func (cq clusterQueries) collect(host string, port string, apiPrefix string) (clusterQueries, error) {
	resp, err := http.Get("http://" + host + ":" + port + "/" + apiPrefix + "query?state=RUNNING")
	if err != nil {
		log.Errorf("%s", err)
		return clusterQueries{}, err
	}
	if resp.StatusCode != 200 {
		log.Errorf("%s", err)
		return clusterQueries{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("%s", err)
		return clusterQueries{}, err
	}

	_ = json.Unmarshal(body, &cq)
	cqq := cq
	for index, query := range cq {
		queryStats := query.QueryStats
		queryStats.ElapsedTimeParsed, err = fromHumanDuration(queryStats.ElapsedTime)
		queryStats.ExecutionTimeParsed, _ = fromHumanDuration(queryStats.ExecutionTime)
		queryStats.PeakTotalMemoryReservationParsed, _ = fromHumanSize(queryStats.PeakTotalMemoryReservation)
		queryStats.PeakUserMemoryReservationParsed, _ = fromHumanSize(queryStats.PeakUserMemoryReservation)
		queryStats.QueuedTimeParsed, _ = fromHumanDuration(queryStats.QueuedTime)
		queryStats.RawInputDataSizeParsed, _ = fromHumanSize(queryStats.RawInputDataSize)
		queryStats.TotalCpuTimeParsed, _ = fromHumanDuration(queryStats.TotalCpuTime)
		queryStats.TotalMemoryReservationParsed, _ = fromHumanSize(queryStats.TotalMemoryReservation)
		queryStats.TotalScheduledTimeParsed, _ = fromHumanDuration(queryStats.TotalScheduledTime)
		queryStats.UserMemoryReservationParsed, _ = fromHumanSize(queryStats.UserMemoryReservation)
		cqq[index].QueryStats = queryStats
	}
	return cqq, nil
}

type ClusterMetrics struct {
	RunningQueries   float64 `json:"runningQueries"`
	BlockedQueries   float64 `json:"blockedQueries"`
	QueuedQueries    float64 `json:"queuedQueries"`
	ActiveWorkers    float64 `json:"activeWorkers"`
	RunningDrivers   float64 `json:"runningDrivers"`
	ReservedMemory   float64 `json:"reservedMemory"`
	TotalInputRows   float64 `json:"totalInputRows"`
	TotalInputBytes  float64 `json:"totalInputBytes"`
	TotalCpuTimeSecs float64 `json:"totalCpuTimeSecs"`
}

func (cm ClusterMetrics) collect(host string, port string, apiPrefix string) (ClusterMetrics, error) {
	var resp *http.Response
	var err error
	if strings.Contains(apiPrefix, "ui") {
		resp, err = http.Get("http://" + host + ":" + port + "/" + apiPrefix + "stats")
	} else {
		resp, err = http.Get("http://" + host + ":" + port + "/" + apiPrefix + "cluster")
	}
	if err != nil {
		log.Errorf("%s", err)
		return ClusterMetrics{}, err
	}
	if resp.StatusCode != 200 {
		log.Errorf("%s", err)
		return ClusterMetrics{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("%s", err)
		return ClusterMetrics{}, err
	}

	_ = json.Unmarshal(body, &cm)
	return cm, err
}
