FROM alpine

USER root

COPY ./etc /etc/
COPY userapi /usr/bin

WORKDIR /usr/bin/

EXPOSE 8000

ENTRYPOINT ["/usr/bin/userapi", "-f","/etc/user.yaml"]