package main

import (
	"encoding/json"
	"github.com/prometheus/common/log"
	"io/ioutil"
	"net/http"
)

type worker struct {
	AvailableProcessors string `json:"availableProcessors"`
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
	resp, err := http.Get("http://" + host + ":" + port + "/v1/cluster/workerMemory")
	if err != nil {
		log.Errorf("%s", err)
		return workers{}, err
	}
	if resp.StatusCode != 200 {
		log.Errorf("%s", err)
		return workers{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("%s", err)
		return workers{}, err
	}

	_ = json.Unmarshal(body, &w)
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
		} `json:"pools"`
	} `json:"memoryInfo"`
	TotalNodeMemory int64 `json:"totalNodeMemory"`
}

func (wm workerMetrics) collect(host string, port string, nodeId string) (workerMetrics, error) {
	resp, err := http.Get("http://" + host + ":" + port + "/v1/worker/" + nodeId + "/status")
	if err != nil {
		log.Errorf("%s", err)
		return workerMetrics{}, err
	}
	if resp.StatusCode != 200 {
		log.Errorf("%s", resp.StatusCode)
		log.Errorf("%s", err)
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
		CompletedDrivers           int64   `json:"completedDrivers"`
		CreateTime                 string  `json:"createTime"`
		CumulativeUserMemory       float64 `json:"cumulativeUserMemory"`
		ElapsedTime                string  `json:"elapsedTime"`
		ExecutionTime              string  `json:"executionTime"`
		FullyBlocked               bool    `json:"fullyBlocked"`
		PeakTotalMemoryReservation string  `json:"peakTotalMemoryReservation"`
		PeakUserMemoryReservation  string  `json:"peakUserMemoryReservation"`
		QueuedDrivers              int64   `json:"queuedDrivers"`
		QueuedTime                 string  `json:"queuedTime"`
		RawInputDataSize           string  `json:"rawInputDataSize"`
		RawInputPositions          string  `json:"rawInputPositions"`
		RunningDrivers             int64   `json:"runningDrivers"`
		TotalCpuTime               string  `json:"totalCpuTime"`
		TotalDrivers               int64   `json:"totalDrivers"`
		TotalMemoryReservation     string  `json:"totalMemoryReservation"`
		TotalScheduledTime         string  `json:"totalScheduledTime"`
		UserMemoryReservation      string  `json:"userMemoryReservation"`
	} `json:"queryStats"`
	State string `json:"state"`
}

type clusterQueries []clusterQuery

func (cq clusterQueries) collect(host string, port string) (clusterQueries, error) {
	resp, err := http.Get("http://" + host + ":" + port + "/v1/query?state=RUNNING")
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
	return cq, nil
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

func (cm ClusterMetrics) collect(host string, port string) (ClusterMetrics, error) {
	resp, err := http.Get("http://" + host + ":" + port + "/v1/cluster")
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
