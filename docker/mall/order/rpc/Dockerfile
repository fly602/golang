FROM alpine

USER root

COPY ./etc /etc/
COPY orderrpc /usr/bin

WORKDIR /usr/bin/

EXPOSE 9002

ENTRYPOINT ["/usr/bin/orderrpc", "-f","/etc/order.yaml"]