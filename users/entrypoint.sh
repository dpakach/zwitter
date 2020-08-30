#!/bin/sh

/bin/server &
envoy -c /etc/envoy.yaml --service-cluster users
