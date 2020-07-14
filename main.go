package main

import (
	"github.com/prometheus/common/log"
	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/prometheus/common/version"
	"time"
)

func main() {
	var (
		listenAddress       = kingpin.Flag("web.listen-address", "Address on which to expose metrics and web interface.").Default(":9483").String()
		serviceName         = kingpin.Flag("web.service-name", "service name to run cloudwatch/prometheus").Enum("cloudwatch", "prometheus")
		prestoHost          = kingpin.Flag("web.presto-host", "presto host").Default("localhost").String()
		prestoPort          = kingpin.Flag("web.presto-port", "presto port").Default("8080").String()
		stackName           = kingpin.Flag("web.stack-name", "presto stack name").Default("atlan-local-test-presto-stack").String()
		cloudWatchNamespace = kingpin.Flag("web.cloudwatch-namespace", "cloudwatch namespace").Default("presto").String()
		apiPrefix           = kingpin.Flag("web.api-prefix", "API prefix to use for making requests to presto").Default("v1/").String()
		cloudwatchRegion    = kingpin.Flag("web.cloudwatch-region", "cloudwatch region").Default("ap-south-1").String()
	)
	log.AddFlags(kingpin.CommandLine)
	kingpin.Version(version.Print("presto_exporter"))
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()
	if *serviceName == "cloudwatch" {
		for {
			cloudwatchAgentStart(*prestoHost, *prestoPort, *cloudWatchNamespace, *stackName, *apiPrefix, *cloudwatchRegion)
			time.Sleep(10 * time.Second)
		}
	} else {
		if *serviceName == "prometheus" {
			prometheusExporterStart(*prestoHost, *prestoPort, *stackName, *listenAddress, *apiPrefix)
		}
	}
}
