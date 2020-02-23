package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"reflect"
	"strings"
)

type CloudWatch struct {
	*client.Client
}

func (c *CloudWatch) PutClusterMetricData(host string, port string, namespace string, stackName string) error {
	svc := cloudwatch.New(session.Must(session.NewSession()),
		aws.NewConfig().WithRegion("ap-south-1"))

	// cluster level metrics
	metricInput := new(cloudwatch.PutMetricDataInput)
	metricInput.SetNamespace(namespace)
	var metricsData []*cloudwatch.MetricDatum
	cm := ClusterMetrics{}
	cm, err := cm.collect(host, port)

	v := reflect.ValueOf(cm)
	typeOfS := v.Type()
	values := make([]interface{}, v.NumField())
	for i := 0; i < v.NumField(); i++ {
		values[i] = v.Field(i).Interface()
		metricData := new(cloudwatch.MetricDatum)
		metricData.SetMetricName(typeOfS.Field(i).Name)
		metricData.SetValue(v.Field(i).Interface().(float64))
		metricData.SetUnit("None")
		dimension := new(cloudwatch.Dimension)
		dimension.SetName("prestoStackName")
		dimension.SetValue(stackName)
		var dimensions []*cloudwatch.Dimension
		dimensions = append(dimensions, dimension)
		metricData.SetDimensions(dimensions)
		metricsData = append(metricsData, metricData)
	}
	metricInput.SetMetricData(metricsData)
	_, err = svc.PutMetricData(metricInput)
	if err != nil {
		_ = fmt.Errorf("%s", err)
		return err
	}
	return nil
}

func (c *CloudWatch) PutWorkerMetricData(host string, port string, namespace string, stackName string) error {
	svc := cloudwatch.New(session.Must(session.NewSession()),
		aws.NewConfig().WithRegion("ap-south-1"))

	// worker level metrics
	metricInput := new(cloudwatch.PutMetricDataInput)
	metricInput.SetNamespace(namespace)
	var metricsData []*cloudwatch.MetricDatum
	workers := workers{}
	workers, _ = workers.collect(host, port)
	var clusterGeneralPoolTotalMemory []float64
	var clusterGeneralPoolFreeMemory []float64
	var clusterGeneralPoolReservedMemory []float64
	var clusterGeneralPoolRevocableMemory []float64

	var clusterUserCPUUtilisation []float64
	var clusterSystemCPUUtilisation []float64

	for k, _ := range workers {
		wm := workerMetrics{}
		workerId := strings.Split(k, " ")[0]
		wm, err := wm.collect(host, port, workerId)
		if err != nil {
			fmt.Println(err)
			return err
		}
		clusterGeneralPoolFreeMemory = append(clusterGeneralPoolFreeMemory, float64(wm.MemoryInfo.Pools.General.FreeBytes))
		clusterGeneralPoolTotalMemory = append(clusterGeneralPoolTotalMemory, float64(wm.MemoryInfo.Pools.General.MaxBytes))
		clusterGeneralPoolReservedMemory = append(clusterGeneralPoolReservedMemory, float64(wm.MemoryInfo.Pools.General.ReservedBytes))
		clusterGeneralPoolRevocableMemory = append(clusterGeneralPoolRevocableMemory, float64(wm.MemoryInfo.Pools.General.ReservedRevocableBytes))

		clusterUserCPUUtilisation = append(clusterUserCPUUtilisation, wm.ProcessCpuLoad)
		clusterSystemCPUUtilisation = append(clusterSystemCPUUtilisation, wm.SystemCpuLoad)

		v := reflect.ValueOf(wm)
		typeOfS := v.Type()
		values := make([]interface{}, v.NumField())
		for i := 0; i < v.NumField(); i++ {
			values[i] = v.Field(i).Interface()
			metricData := new(cloudwatch.MetricDatum)
			metricData.SetMetricName(typeOfS.Field(i).Name)
			switch g := values[i].(type) {
			case float64:
				metricData.SetValue(values[i].(float64))
			case int64:
				metricData.SetValue(float64(values[i].(int64)))
			default:
				fmt.Println(g)
				continue
			}
			metricData.SetUnit("None")
			var dimensions []*cloudwatch.Dimension
			dimension := new(cloudwatch.Dimension)
			dimension.SetName("prestoStackName")
			dimension.SetValue(stackName)
			dimensions = append(dimensions, dimension)
			dimension = new(cloudwatch.Dimension)
			dimension.SetName("prestoWorkerId")
			dimension.SetValue(workerId)
			dimensions = append(dimensions, dimension)
			metricData.SetDimensions(dimensions)
			metricsData = append(metricsData, metricData)
		}
	}

	metricInput.SetMetricData(metricsData)
	_, err := svc.PutMetricData(metricInput)
	if err != nil {
		_ = fmt.Errorf("%s", err)
		return err
	}

	// cluster memory metrics
	metricInput = new(cloudwatch.PutMetricDataInput)
	metricInput.SetNamespace(namespace)
	var clusterMemoryMetricsData []*cloudwatch.MetricDatum

	clusterMemoryMetrics := ClusterMemoryMetrics{
		Sum(clusterGeneralPoolFreeMemory),
		Sum(clusterGeneralPoolTotalMemory),
		Sum(clusterGeneralPoolReservedMemory),
		Sum(clusterGeneralPoolRevocableMemory),
		Median(clusterGeneralPoolFreeMemory),
		Sum(clusterGeneralPoolFreeMemory) / float64(len(workers)),
	}
	v := reflect.ValueOf(clusterMemoryMetrics)
	typeOfS := v.Type()
	for i := 0; i < v.NumField(); i++ {
		metricData := new(cloudwatch.MetricDatum)
		metricData.SetMetricName(typeOfS.Field(i).Name)
		metricData.SetValue(v.Field(i).Interface().(float64))
		metricData.SetUnit("None")
		var dimensions []*cloudwatch.Dimension
		dimension := new(cloudwatch.Dimension)
		dimension.SetName("prestoStackName")
		dimension.SetValue(stackName)
		dimensions = append(dimensions, dimension)
		metricData.SetDimensions(dimensions)
		clusterMemoryMetricsData = append(clusterMemoryMetricsData, metricData)
	}
	metricInput.SetMetricData(clusterMemoryMetricsData)
	_, err = svc.PutMetricData(metricInput)
	if err != nil {
		_ = fmt.Errorf("%s", err)
		return err
	}

	// cluster CPU metrics
	metricInput = new(cloudwatch.PutMetricDataInput)
	metricInput.SetNamespace(namespace)
	var clusterCPUMetricsData []*cloudwatch.MetricDatum

	clusterCPUMetrics := ClusterCPUMetrics{
		ClusterUserCPUUtilisation:        Sum(clusterUserCPUUtilisation),
		ClusterSystemCPUUtilisation:      Sum(clusterSystemCPUUtilisation),
		MedianWorkerUserCPUUtilisation:   Median(clusterUserCPUUtilisation),
		MedianWorkerSystemCPUUtilisation: Median(clusterSystemCPUUtilisation),
		MeanWorkerUserCPUUtilisation:     Sum(clusterUserCPUUtilisation) / float64(len(workers)),
		MeanWorkerSystemCPUUtilisation:   Sum(clusterSystemCPUUtilisation) / float64(len(workers)),
	}
	v = reflect.ValueOf(clusterCPUMetrics)
	typeOfS = v.Type()
	for i := 0; i < v.NumField(); i++ {
		metricData := new(cloudwatch.MetricDatum)
		metricData.SetMetricName(typeOfS.Field(i).Name)
		metricData.SetValue(v.Field(i).Interface().(float64))
		metricData.SetUnit("None")
		var dimensions []*cloudwatch.Dimension
		dimension := new(cloudwatch.Dimension)
		dimension.SetName("prestoStackName")
		dimension.SetValue(stackName)
		dimensions = append(dimensions, dimension)
		metricData.SetDimensions(dimensions)
		clusterCPUMetricsData = append(clusterCPUMetricsData, metricData)
	}
	metricInput.SetMetricData(clusterCPUMetricsData)
	_, err = svc.PutMetricData(metricInput)
	if err != nil {
		_ = fmt.Errorf("%s", err)
		return err
	}

	return nil
}

func (c *CloudWatch) PutMetricData(host string, port string, namespace string, stackName string) error {
	err := c.PutClusterMetricData(host, port, namespace, stackName)
	if err != nil {
		fmt.Println(err)
	}

	err = c.PutWorkerMetricData(host, port, namespace, stackName)
	if err != nil {
		fmt.Println(err)
	}
	return nil
}

func cloudwatchAgentStart(host string, port string, namespace string, stackName string) {
	c := CloudWatch{}
	err := c.PutMetricData(host, port, namespace, stackName)
	if err != nil {
		fmt.Errorf("%s", err)
	}
}
