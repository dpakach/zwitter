# from base image golang:latest
FROM envoyproxy/envoy-alpine-dev:latest

# copy necessary binary and config files
COPY bin /bin
COPY config /config
COPY cert /cert
ADD ./envoy.yaml /etc/envoy.yaml
ADD ./entrypoint.sh /entrypoint.sh

# set workdir
WORKDIR /

# Running the project
CMD ["/entrypoint.sh"]
