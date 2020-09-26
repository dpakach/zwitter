# from base image golang:latest
# FROM envoyproxy/envoy-alpine-dev:latest
FROM alpine:3.7

ARG service_name

# copy necessary binary and config files
# COPY ./bin /zwitter-bin
# COPY config /config
COPY ./service/cert /cert
ADD ./${service_name}/envoy.yaml /etc/envoy.yaml
ADD ./service/entrypoint.sh /entrypoint.sh

# set workdir
WORKDIR /

CMD ["/entrypoint.sh"]
