FROM alpine

USER root

COPY ./etc /etc/
COPY payapi /usr/bin

WORKDIR /usr/bin/

EXPOSE 8003

ENTRYPOINT ["/usr/bin/payapi", "-f","/etc/user.yaml"]