FROM alpine

USER root

COPY ./etc /etc/
COPY productrpc /usr/bin

RUN apk add tzdata

WORKDIR /usr/bin/

EXPOSE 9001

ENTRYPOINT ["/usr/bin/productrpc", "-f","/etc/product.yaml"]