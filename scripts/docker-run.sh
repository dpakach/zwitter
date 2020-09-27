#!/bin/bash

docker run --rm \
  -e GODEBUG='x509ignoreCN=0' \
  -v $(pwd)/auth/config:/config \
  -v $(pwd)/posts/bin:/zwitter-bin \
  -p 9999:9999 \
  -p 9990:9990 \
  -d --name=zauth \
  zwiter-auth

docker run --rm \
  -e GODEBUG='x509ignoreCN=0' \
  -v $(pwd)/posts/config:/config \
  -v $(pwd)/posts/bin:/zwitter-bin \
  -p 7777:7777 \
  -p 7778:7778 \
  -d --name=zposts \
  zwiter-posts

docker run --rm \
  -e GODEBUG='x509ignoreCN=0' \
  -v $(pwd)/users/config:/config \
  -v $(pwd)/posts/bin:/zwitter-bin \
  -p 8888:8888 \
  -p 8889:8889 \
  -d --name=zusers \
  zwiter-users
