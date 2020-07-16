# Presto Metrics 
![Go](https://github.com/atlanhq/presto-metrics/workflows/Go/badge.svg)

This is an agent which scraps metrics from Presto and pushes it to AWS Cloudwatch and also serves them on a prometheus compatible endpoint. 

## Table of Contents

- [Compatibility](https://github.com/atlanhq/presto-metrics#Compatibility)
- [Build and Run](https://github.com/atlanhq/presto-metrics#Build%20and%20Run)
- [Metrics](https://github.com/atlanhq/presto-metrics#Metrics)
- [Contribute](https://github.com/atlanhq/presto-metrics#Contribute)

### Compatibility

 [Presto](https://prestosql.io):v330

Golang: 1.13.4

### Build and Run

```bash
go build .
```

We recommend to install the agent inside the presto coordinator. See [releases page](https://github.com/atlanhq/presto-metrics/releases) to download pre-built binaries. 

#### Run prometheus exporter

```bash
./presto_metrics --web.service-name=prometheus \
                --web.presto-host=localhost \
                --web.presto-port=8080 \
                --web.stack-name=<label_for_metrics_dimension> \
                --web.listen-address=<port_to_expose_metrics_on>
```

#### Run cloudwatch agent

```bash
./presto_metrics --web.service-name=cloudwatch \
                --web.presto-host=localhost \
                --web.presto-port=8080 \
                --web.stack-name=<label_for_metrics_dimension> \
                --web.cloudwatch-namespace=<cloudwatch_namespace> \
                --web.cloudwatch-region=<cloudwatch_region_to_push_metrics>
```

### Metrics

```bash
# HELP atlan_presto_active_workers Active workers of the presto cluster.
# TYPE atlan_presto_active_workers gauge
atlan_presto_active_workers{prestoStackName="presto-oss-107"} 2
# HELP atlan_presto_blocked_queries Blocked queries of the presto cluster.
# TYPE atlan_presto_blocked_queries gauge
atlan_presto_blocked_queries{prestoStackName="presto-oss-107"} 0
# HELP atlan_presto_cluster_general_pool_free_memory total free general pool memory of cluster.
# TYPE atlan_presto_cluster_general_pool_free_memory gauge
atlan_presto_cluster_general_pool_free_memory{prestoStackName="presto-oss-107"} 1.1274289152e+10
# HELP atlan_presto_cluster_general_pool_reserved_memory total general pool reserved memory of cluster.
# TYPE atlan_presto_cluster_general_pool_reserved_memory gauge
atlan_presto_cluster_general_pool_reserved_memory{prestoStackName="presto-oss-107"} 0
# HELP atlan_presto_cluster_general_pool_revocable_memory total general pool revocable memory of cluster.
# TYPE atlan_presto_cluster_general_pool_revocable_memory gauge
atlan_presto_cluster_general_pool_revocable_memory{prestoStackName="presto-oss-107"} 0
# HELP atlan_presto_cluster_general_pool_total_memory total general pool memory of cluster.
# TYPE atlan_presto_cluster_general_pool_total_memory gauge
atlan_presto_cluster_general_pool_total_memory{prestoStackName="presto-oss-107"} 1.1274289152e+10
# HELP atlan_presto_cluster_system_cpu_utilisation cluster system cpu utlisation
# TYPE atlan_presto_cluster_system_cpu_utilisation gauge
atlan_presto_cluster_system_cpu_utilisation{prestoStackName="presto-oss-107"} 0.13328566168251338
# HELP atlan_presto_cluster_user_cpu_utilisation cluster user cpu utlisation
# TYPE atlan_presto_cluster_user_cpu_utilisation gauge
atlan_presto_cluster_user_cpu_utilisation{prestoStackName="presto-oss-107"} 0.13061522731158193
# HELP atlan_presto_mean_workers_general_pool_memory mean workers general pool memory
# TYPE atlan_presto_mean_workers_general_pool_memory gauge
atlan_presto_mean_workers_general_pool_memory{prestoStackName="presto-oss-107"} 3.758096384e+09
# HELP atlan_presto_mean_workers_system_cpu_utilisation mean workers system cpu utlisation
# TYPE atlan_presto_mean_workers_system_cpu_utilisation gauge
atlan_presto_mean_workers_system_cpu_utilisation{prestoStackName="presto-oss-107"} 0.04442855389417113
# HELP atlan_presto_mean_workers_user_cpu_utilisation mean workers user cpu utlisation
# TYPE atlan_presto_mean_workers_user_cpu_utilisation gauge
atlan_presto_mean_workers_user_cpu_utilisation{prestoStackName="presto-oss-107"} 0.04353840910386064
# HELP atlan_presto_median_workers_general_pool_memory median workers general pool memory
# TYPE atlan_presto_median_workers_general_pool_memory gauge
atlan_presto_median_workers_general_pool_memory{prestoStackName="presto-oss-107"} 3.758096384e+09
# HELP atlan_presto_median_workers_system_cpu_utilisation median workers system cpu utlisation
# TYPE atlan_presto_median_workers_system_cpu_utilisation gauge
atlan_presto_median_workers_system_cpu_utilisation{prestoStackName="presto-oss-107"} 0.030181086519114688
# HELP atlan_presto_median_workers_user_cpu_utilisation median workers user cpu utlisation
# TYPE atlan_presto_median_workers_user_cpu_utilisation gauge
atlan_presto_median_workers_user_cpu_utilisation{prestoStackName="presto-oss-107"} 0.03118712273641851
# HELP atlan_presto_queued_queries Queued queries of the presto cluster.
# TYPE atlan_presto_queued_queries gauge
atlan_presto_queued_queries{prestoStackName="presto-oss-107"} 0
# HELP atlan_presto_reserved_memory Reserved memory of the presto cluster.
# TYPE atlan_presto_reserved_memory gauge
atlan_presto_reserved_memory{prestoStackName="presto-oss-107"} 0
# HELP atlan_presto_running_drivers Running drivers of the presto cluster.
# TYPE atlan_presto_running_drivers gauge
atlan_presto_running_drivers{prestoStackName="presto-oss-107"} 0
# HELP atlan_presto_running_queries Running requests of the presto cluster.
# TYPE atlan_presto_running_queries gauge
atlan_presto_running_queries{prestoStackName="presto-oss-107"} 0
# HELP atlan_presto_total_cpu_time_secs Total cpu time of the presto cluster.
# TYPE atlan_presto_total_cpu_time_secs gauge
atlan_presto_total_cpu_time_secs{prestoStackName="presto-oss-107"} 0
# HELP atlan_presto_total_input_bytes Total input bytes of the presto cluster.
# TYPE atlan_presto_total_input_bytes gauge
atlan_presto_total_input_bytes{prestoStackName="presto-oss-107"} 1377
# HELP atlan_presto_total_input_rows Total input rows of the presto cluster.
# TYPE atlan_presto_total_input_rows gauge
atlan_presto_total_input_rows{prestoStackName="presto-oss-107"} 27
# HELP atlan_presto_worker_general_pool_free_bytes general pool free bytes
# TYPE atlan_presto_worker_general_pool_free_bytes gauge
atlan_presto_worker_general_pool_free_bytes{prestoStackName="presto-oss-107",prestoWorkerId="i-020e347885a9783bf"} 3.758096384e+09
atlan_presto_worker_general_pool_free_bytes{prestoStackName="presto-oss-107",prestoWorkerId="i-03ada9c86d2f19254"} 3.758096384e+09
atlan_presto_worker_general_pool_free_bytes{prestoStackName="presto-oss-107",prestoWorkerId="i-0540e97929f83f8bd"} 3.758096384e+09
# HELP atlan_presto_worker_general_pool_max_bytes general pool max bytes
# TYPE atlan_presto_worker_general_pool_max_bytes gauge
atlan_presto_worker_general_pool_max_bytes{prestoStackName="presto-oss-107",prestoWorkerId="i-020e347885a9783bf"} 3.758096384e+09
atlan_presto_worker_general_pool_max_bytes{prestoStackName="presto-oss-107",prestoWorkerId="i-03ada9c86d2f19254"} 3.758096384e+09
atlan_presto_worker_general_pool_max_bytes{prestoStackName="presto-oss-107",prestoWorkerId="i-0540e97929f83f8bd"} 3.758096384e+09
# HELP atlan_presto_worker_general_pool_reserved_bytes general pool reserved bytes
# TYPE atlan_presto_worker_general_pool_reserved_bytes gauge
atlan_presto_worker_general_pool_reserved_bytes{prestoStackName="presto-oss-107",prestoWorkerId="i-020e347885a9783bf"} 0
atlan_presto_worker_general_pool_reserved_bytes{prestoStackName="presto-oss-107",prestoWorkerId="i-03ada9c86d2f19254"} 0
atlan_presto_worker_general_pool_reserved_bytes{prestoStackName="presto-oss-107",prestoWorkerId="i-0540e97929f83f8bd"} 0
# HELP atlan_presto_worker_general_pool_reserved_revocable_bytes general pool reserved revocable bytes
# TYPE atlan_presto_worker_general_pool_reserved_revocable_bytes gauge
atlan_presto_worker_general_pool_reserved_revocable_bytes{prestoStackName="presto-oss-107",prestoWorkerId="i-020e347885a9783bf"} 0
atlan_presto_worker_general_pool_reserved_revocable_bytes{prestoStackName="presto-oss-107",prestoWorkerId="i-03ada9c86d2f19254"} 0
atlan_presto_worker_general_pool_reserved_revocable_bytes{prestoStackName="presto-oss-107",prestoWorkerId="i-0540e97929f83f8bd"} 0
# HELP atlan_presto_worker_heap_available heap available for worker
# TYPE atlan_presto_worker_heap_available gauge
atlan_presto_worker_heap_available{prestoStackName="presto-oss-107",prestoWorkerId="i-020e347885a9783bf"} 5.36870912e+09
atlan_presto_worker_heap_available{prestoStackName="presto-oss-107",prestoWorkerId="i-03ada9c86d2f19254"} 5.36870912e+09
atlan_presto_worker_heap_available{prestoStackName="presto-oss-107",prestoWorkerId="i-0540e97929f83f8bd"} 5.36870912e+09
# HELP atlan_presto_worker_heap_used heap used for worker
# TYPE atlan_presto_worker_heap_used gauge
atlan_presto_worker_heap_used{prestoStackName="presto-oss-107",prestoWorkerId="i-020e347885a9783bf"} 2.23394816e+08
atlan_presto_worker_heap_used{prestoStackName="presto-oss-107",prestoWorkerId="i-03ada9c86d2f19254"} 1.56223552e+08
atlan_presto_worker_heap_used{prestoStackName="presto-oss-107",prestoWorkerId="i-0540e97929f83f8bd"} 1.80819456e+08
# HELP atlan_presto_worker_non_heap_used non heap used for worker
# TYPE atlan_presto_worker_non_heap_used gauge
atlan_presto_worker_non_heap_used{prestoStackName="presto-oss-107",prestoWorkerId="i-020e347885a9783bf"} 1.64377432e+08
atlan_presto_worker_non_heap_used{prestoStackName="presto-oss-107",prestoWorkerId="i-03ada9c86d2f19254"} 1.37667376e+08
atlan_presto_worker_non_heap_used{prestoStackName="presto-oss-107",prestoWorkerId="i-0540e97929f83f8bd"} 1.38011408e+08
# HELP atlan_presto_worker_process_cpu_load user cpu load
# TYPE atlan_presto_worker_process_cpu_load gauge
atlan_presto_worker_process_cpu_load{prestoStackName="presto-oss-107",prestoWorkerId="i-020e347885a9783bf"} 0.03118712273641851
atlan_presto_worker_process_cpu_load{prestoStackName="presto-oss-107",prestoWorkerId="i-03ada9c86d2f19254"} 0.07720588235294118
atlan_presto_worker_process_cpu_load{prestoStackName="presto-oss-107",prestoWorkerId="i-0540e97929f83f8bd"} 0.022222222222222223
# HELP atlan_presto_worker_processors total processors a worker has
# TYPE atlan_presto_worker_processors gauge
atlan_presto_worker_processors{prestoStackName="presto-oss-107",prestoWorkerId="i-020e347885a9783bf"} 2
atlan_presto_worker_processors{prestoStackName="presto-oss-107",prestoWorkerId="i-03ada9c86d2f19254"} 2
atlan_presto_worker_processors{prestoStackName="presto-oss-107",prestoWorkerId="i-0540e97929f83f8bd"} 2
# HELP atlan_presto_worker_system_cpu_load system cpu load
# TYPE atlan_presto_worker_system_cpu_load gauge
atlan_presto_worker_system_cpu_load{prestoStackName="presto-oss-107",prestoWorkerId="i-020e347885a9783bf"} 0.030181086519114688
atlan_presto_worker_system_cpu_load{prestoStackName="presto-oss-107",prestoWorkerId="i-03ada9c86d2f19254"} 0.08088235294117647
atlan_presto_worker_system_cpu_load{prestoStackName="presto-oss-107",prestoWorkerId="i-0540e97929f83f8bd"} 0.022222222222222223
# HELP atlan_presto_worker_total_node_memory total node memory
# TYPE atlan_presto_worker_total_node_memory gauge
atlan_presto_worker_total_node_memory{prestoStackName="presto-oss-107",prestoWorkerId="i-020e347885a9783bf"} 0
atlan_presto_worker_total_node_memory{prestoStackName="presto-oss-107",prestoWorkerId="i-03ada9c86d2f19254"} 0
atlan_presto_worker_total_node_memory{prestoStackName="presto-oss-107",prestoWorkerId="i-0540e97929f83f8bd"} 0
```

### Contribute

1. Fork it
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Add some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create new Pull Request
