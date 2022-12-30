FROM alpine

USER root

COPY ./etc /etc/
COPY payapi /usr/bin

WORKDIR /usr/bin/

RUN apk add tzdata

EXPOSE 8003

ENTRYPOINT ["/usr/bin/payapi", "-f","/etc/pay.yaml"]