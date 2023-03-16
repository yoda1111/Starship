FROM ubuntu:22.04

ENV DEBIAN_FRONTEND=noninteractive
RUN apt-get update -y
RUN apt-get upgrade -y

RUN apt-get install -y libelf-dev wget

ENV WASI_VERSION 14
ENV WASI_VERSION_FULL=${WASI_VERSION}.0
RUN wget https://github.com/WebAssembly/wasi-sdk/releases/download/wasi-sdk-${WASI_VERSION}/wasi-sdk-${WASI_VERSION_FULL}-linux.tar.gz
RUN tar xvf wasi-sdk-${WASI_VERSION_FULL}-linux.tar.gz -C /opt
RUN rm -rf wasi-sdk-${WASI_VERSION_FULL}-linux.tar.gz && mv /opt/wasi-sdk-${WASI_VERSION_FULL} /opt/wasi-sdk

RUN mkdir -p /opt/starship/include
COPY common/* /opt/starship/include/