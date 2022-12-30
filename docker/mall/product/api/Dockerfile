FROM alpine

USER root

COPY ./etc /etc/
COPY productapi /usr/bin

RUN apk add tzdata

WORKDIR /usr/bin/

EXPOSE 8001

ENTRYPOINT ["/usr/bin/productapi", "-f","/etc/product.yaml"]