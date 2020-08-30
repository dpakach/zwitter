# from base image golang:latest
FROM golang:latest

# copy necessary binary and config files
COPY bin /bin
COPY config /config
COPY cert /cert

# set workdir
WORKDIR /

# Running the project
CMD ["/bin/server"]
