FROM go1.18

USER root

COPY ./etc /etc/
COPY orderrpc /usr/bin

WORKDIR /usr/bin/

# RUN apk add tzdata

EXPOSE 9002

USER normal

ENTRYPOINT ["/usr/bin/orderrpc", "-f","/etc/order.yaml"]