#!/bin/sh

/root/presto_metrics  --web.service-name="${service_name:-prometheus}" \
                      --web.presto-host="${presto_host:-localhost}" \
                      --web.presto-port="${presto_port:-8080}" \
                      --web.stack-name="${stack_name:-atlan-local-test-presto-stack}" \
                      --web.cloudwatch-namespace="${cloudwatch_namespace:-presto}"