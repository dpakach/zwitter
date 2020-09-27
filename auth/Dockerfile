FROM envoyproxy/envoy-alpine-dev:latest

# copy necessary binary and config files
COPY bin /bin
COPY config /config
COPY cert /cert
ADD ./envoy.yaml /etc/envoy.yaml
ADD ./entrypoint.sh /entrypoint.sh

EXPOSE 9990
EXPOSE 9999

# set workdir
WORKDIR /

# set entrypoint 
CMD ["/entrypoint.sh"]
