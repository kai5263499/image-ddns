FROM ubuntu:19.04

LABEL MAINTAINER="Wes Widner <kai5263499@gmail.com>"

ENV CGO_ENABLED=1 CGO_CPPFLAGS="-I/usr/include"
ENV GOPATH=/go
ENV DEBIAN_FRONTEND=noninteractive
ENV PATH=$GOPATH/bin:/usr/local/go/bin:$PATH

COPY . /go/src/github.com/kai5263499/image-ddns

WORKDIR /go/src/github.com/kai5263499/image-ddns

RUN echo "Install apt packages" && \
    apt-get update && \
    apt-get install -y \
    tesseract-ocr \
    libtesseract-dev \
    libtesseract4 \
    libhyperscan-dev \
    ragel \
    curl \
    cmake \
    pkg-config \
    g++

RUN echo "Install golang" && \
	curl -sLO https://dl.google.com/go/go1.13.3.linux-amd64.tar.gz && \
	tar -xf go1.13.3.linux-amd64.tar.gz && \
	mv go /usr/local && \
	rm -rf go1.13.3.linux-amd64.tar.gz

RUN echo "Caching golang modules" && \
	go mod download
