FROM alpine

USER root

COPY ./etc /etc/
COPY productapi /usr/bin

WORKDIR /usr/bin/

EXPOSE 8001

ENTRYPOINT ["/usr/bin/productapi", "-f","/etc/product.yaml"]