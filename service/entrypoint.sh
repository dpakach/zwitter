#!/bin/sh

/zwitter-bin/server &
envoy -c /etc/envoy.yaml --service-cluster service${SERVICE_NAME}
