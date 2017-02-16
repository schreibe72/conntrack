FROM ubuntu:14.04

RUN apt-get update;\
apt-get install -y wget build-essential tcpdump libpcap-dev;\
cd /usr/local;\
wget https://storage.googleapis.com/golang/go1.7.5.linux-amd64.tar.gz;\
tar xfvz go1.7.5.linux-amd64.tar.gz

ENV PATH=$PATH:/usr/local/go/bin
ENV GOPATH=/root/go-work
