FROM envoyproxy/envoy-alpine-dev:latest

# copy necessary binary and config files
COPY bin /bin
COPY config /config
COPY cert /cert
ADD ./envoy.yaml /etc/envoy.yaml
ADD ./entrypoint.sh /entrypoint.sh

EXPOSE 8889
EXPOSE 8888

# set workdir
WORKDIR /

# set entrypoint
CMD ["/entrypoint.sh"]
