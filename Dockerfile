FROM ubuntu:20.04 AS build

## 设置时区
RUN apt-get -y update && DEBIAN_FRONTEND="noninteractive" apt -y install tzdata
ENV TZ=Asia/Shanghai
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone
# 容器环境变量添加，会覆盖默认的变量值  
ENV GO111MODULE=on\
    GOPROXY=https://goproxy.cn,direct

RUN mkdir /app
WORKDIR /app

RUN apt update && apt-get install -y --no-install-recommends ca-certificates
RUN update-ca-certificates

## 安装 ffmepg
RUN apt install -y ffmpeg

RUN apt install -y git wget sudo
RUN mkdir /app/tmp

## 安装go1.17.8
RUN chmod -R 777 /app/tmp&& cd /app/tmp
RUN wget https://go.dev/dl/go1.17.8.linux-amd64.tar.gz &&\
    tar -C /usr/local -xzf go1.17.8.linux-amd64.tar.gz &&\
    ## 软链接
    ln -s /usr/local/go/bin/* /usr/bin/
    # 设置环境变量
ENV GOPATH="$HOME/go" \
    PATH="$PATH:/usr/local/go/bin:$GOPATH/bin" \
    GOPROXY=https://goproxy.cn,direct
    ## 相当于source /bin/sh 中无source === /bin/bash -c "source ~/.bashrc"
# RUN . ~/.bashrc 

## 配置goav编译环境
RUN apt-get  install -y autoconf automake build-essential libass-dev libfreetype6-dev libsdl1.2-dev libtheora-dev libtool libva-dev libvdpau-dev libvorbis-dev libxcb1-dev libxcb-shm0-dev libxcb-xfixes0-dev pkg-config texi2html zlib1g-dev
RUN apt install -y libavdevice-dev libavfilter-dev libswscale-dev libavcodec-dev libavformat-dev libswresample-dev libavutil-dev
RUN apt-get install -y yasm
    ## 设置环境变量
ENV FFMPEG_ROOT=$HOME/ffmpeg \
    CGO_LDFLAGS="-L$FFMPEG_ROOT/lib/ -lavcodec -lavformat -lavutil -lswscale -lswresample -lavdevice -lavfilter" \
    CGO_CFLAGS="-I$FFMPEG_ROOT/include" \
    LD_LIBRARY_PATH=$HOME/ffmpeg/lib

## 编译安装gocv

RUN cd /app/tmp
RUN git clone https://github.com/hybridgroup/gocv.git \
    && cd gocv \
    && make install

## 设置ai应用环境变量
ENV IMG_NUM ""
##拉流地址 rtmp://aiedge.ndsl-lab.cn:8035/live/stream1
##推流地址 rtmp://aiedge.ndsl-lab.cn:8035/live/stream
ENV PULLSTREAM_URL ""

RUN mkdir /app/aiedge
COPY . /app/aiedge
CMD cd /app/aiedge/&& go run main.go
# CMD /app/aiedge/myapp





