#!/bin/bash
# Tool needs some C Libs and therefore it is easier to build it in a docker Container
if [[ "$(docker images -q conntrack-build:14.04 2> /dev/null)" == "" ]]; then
  docker build -t conntrack-build:14.04 .
fi
docker run -it --rm -v $GOPATH:/root/go-work conntrack-build:14.04 /bin/bash -c "cd ~/go-work/src/github.com/schreibe72/conntrack; go build ."
