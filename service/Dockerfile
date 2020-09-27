# from base image golang:latest
FROM envoyproxy/envoy-alpine-dev:latest

ARG service_name

# copy necessary binary and config files
COPY ./service/cert /cert
ADD ./${service_name}/envoy.yaml /etc/envoy.yaml
ADD ./service/entrypoint.sh /entrypoint.sh

# set workdir
WORKDIR /

CMD ["/entrypoint.sh"]

## old one
#FROM golang:latest
#ARG service_name
#COPY ./service/cert /cert
#ADD ./${service_name}/envoy.yaml /etc/envoy.yaml
#ADD ./service/entrypoint.sh /entrypoint.sh

#WORKDIR /

#CMD ["/entrypoint.sh"]
