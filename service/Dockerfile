#
# alpine build stage
#

# use golang alpine image for build
FROM golang:1-alpine as build

# service name
ARG service_name

# set workdir
RUN mkdir -p /app/${service_name}
WORKDIR /app

# copy necessary files for build
COPY ./pkg /app/pkg
COPY ./auth /app/auth
COPY ./posts /app/posts
COPY ./users /app/users
COPY ./web /app/web
COPY ./media /app/media

COPY ./go.mod /app
COPY ./go.sum /app

RUN ls -la /app
RUN ls -la /app/${service_name}

# build
RUN go build /app/${service_name}/server/main.go

#
# make docker image
#

# from base image golang:latest
FROM envoyproxy/envoy-alpine-dev:latest

ARG service_name

# copy necessary binary and config files
COPY --from=build /app/main /app/main
COPY ./service/cert /cert
ADD ./${service_name}/envoy.yaml /etc/envoy.yaml
ADD ./service/entrypoint.sh /entrypoint.sh

# set workdir
WORKDIR /

CMD ["/entrypoint.sh"]
