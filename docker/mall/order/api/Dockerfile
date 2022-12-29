FROM alpine

USER root

COPY ./etc /etc/
COPY orderapi /usr/bin

WORKDIR /usr/bin/

EXPOSE 8002

ENTRYPOINT ["/usr/bin/orderapi", "-f","/etc/user.yaml"]