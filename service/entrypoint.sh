#!/bin/sh

envoy -c /etc/envoy.yaml --service-cluster service${SERVICE_NAME} &> /dev/null &
/zwitter-bin/server
