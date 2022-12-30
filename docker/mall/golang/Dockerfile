FROM golang:1.18

LABEL maintainer="fuleyi <fuleyi@outlook.com>"

ENV GOPROXY https://goproxy.cn,direct

USER root

RUN sed -i 's/deb.debian.org/mirrors.tuna.tsinghua.edu.cn/' /etc/apt/sources.list && \
    sed -i 's/security.debian.org/mirrors.tuna.tsinghua.edu.cn/' /etc/apt/sources.list && \
    sed -i 's/security-cdn.debian.org/mirrors.tuna.tsinghua.edu.cn/' /etc/apt/sources.list && \  
    apt-get update && \
    apt-get upgrade -y && \
    apt-get install -y --no-install-recommends \
    curl \
    zip \
    unzip \
    git \
    openssh-server \
    bash-completion \
    vim \
    tzdata \
    iputils-ping

# 安装 goctl
USER root
RUN GOPROXY=https://goproxy.cn/,direct go install github.com/zeromicro/go-zero/tools/goctl@latest

RUN echo 'if ! shopt -oq posix; then \n\
  if [ -f /usr/share/bash-completion/bash_completion ]; then \n\
    . /usr/share/bash-completion/bash_completion \n\
  elif [ -f /etc/bash_completion ]; then \n\
    . /etc/bash_completion \n\
  fi \n\
fi'\
>> /etc/bash.bashrc

# 安装 protoc
USER root
COPY protoc-21.12-linux-x86_64.zip /tmp/protoc.zip
RUN unzip -d /tmp/protoc /tmp/protoc.zip && \
    mv /tmp/protoc/bin/protoc $GOPATH/bin

# 安装 protoc-gen-go
USER root
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

# 安装 protoc-gen-go-grpc
USER root
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

RUN goctl completion bash > /etc/bash_completion.d/goctl

# $GOPATH/bin添加到环境变量中
ENV PATH $GOPATH/bin:$PATH
ENV GOPATH $GOPATH:/usr

# 清理垃圾
USER root
RUN apt-get clean && \
    rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/* && \
    rm /var/log/lastlog /var/log/faillog

# 设置工作目录
WORKDIR /usr/src/