FROM alpine

USER root

COPY ./etc /etc/
COPY userrpc /usr/bin

WORKDIR /usr/bin/

RUN apk add tzdata

EXPOSE 9000

ENTRYPOINT ["/usr/bin/userrpc", "-f","/etc/user.yaml"]