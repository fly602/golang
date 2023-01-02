FROM go1.18

USER root

COPY ./etc /etc/
COPY orderapi /usr/bin

WORKDIR /usr/bin/

# RUN apk add tzdata

EXPOSE 8002

ENTRYPOINT ["/usr/bin/orderapi", "-f","/etc/order.yaml"]