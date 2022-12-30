FROM alpine

USER root

COPY ./etc /etc/
COPY payrpc /usr/bin

WORKDIR /usr/bin/

RUN apk add tzdata

EXPOSE 9003

ENTRYPOINT ["/usr/bin/payrpc", "-f","/etc/pay.yaml"]