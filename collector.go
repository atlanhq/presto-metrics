package main

import (
	"encoding/json"
	"github.com/prometheus/common/log"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type worker struct {
	AvailableProcessors string `json:"availableProcessors"`
}

type workers map[string]worker

func (w workers) collect(host string, port string) {
	resp, err := http.Get("http://" + host + ":" + port + "/v1/cluster/workerMemory")
	if err != nil {
		log.Errorf("%s", err)
		return
	}
	if resp.StatusCode != 200 {
		log.Errorf("%s", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("%s", err)
		return
	}

	_ = json.Unmarshal(body, &w)
}

type workerMetrics struct {
	Processors      string  `json:"processors"`
	HeapAvailable   int64   `json:"heapAvailable"`
	HeapUsed        int64   `json:"heapUsed"`
	NonHeapUsed     int64   `json:"nonHeapUsed"`
	ProcessCpuLoad  float64 `json:"processCpuLoad"`
	SystemCpuLoad   float64 `json:"systemCpuLoad"`
	Uptime          float64 `json:"uptime"`
	InternalAddress string  `json:"internalAddress"`
	NodeId          string  `json:"nodeId"`
	NodeVersion     struct {
		Version string `json:"version"`
	} `json:"nodeVersion"`
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
	TotalNodeMemory int64 `json:"totalNodeMemory"`
}

func (wm workerMetrics) collect(host string, port string, nodeId string) {
	resp, err := http.Get("http://" + host + ":" + port + "/v1/worker/" + nodeId + "/status")
	if err != nil {
		log.Errorf("%s", err)
		return
	}
	if resp.StatusCode != 200 {
		log.Errorf("%s", resp.StatusCode)
		log.Errorf("%s", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("%s", err)
		return
	}

	_ = json.Unmarshal(body, &wm)
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

func (cq clusterQueries) collect(host string, port string) {
	resp, err := http.Get("http://" + host + ":" + port + "/v1/query?state=RUNNING")
	if err != nil {
		log.Errorf("%s", err)
		return
	}
	if resp.StatusCode != 200 {
		log.Errorf("%s", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("%s", err)
		return
	}

	_ = json.Unmarshal(body, &cq)
}

type clusterMetrics struct {
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

func (cm clusterMetrics) collect(host string, port string) {
	resp, err := http.Get("http://" + host + ":" + port + "/v1/cluster")
	if err != nil {
		log.Errorf("%s", err)
		return
	}
	if resp.StatusCode != 200 {
		log.Errorf("%s", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("%s", err)
		return
	}

	_ = json.Unmarshal(body, &cm)
}

func main() {
	args := os.Args

	w := workers{}
	w.collect(args[1], args[2])
	for k, _ := range w {
		wm := workerMetrics{}
		wm.collect(args[1], args[2], strings.Split(k, " ")[0])
	}

	cq := clusterQueries{}
	cq.collect(args[1], args[2])

	cm := clusterMetrics{}
	cm.collect(args[1], args[2])
}
