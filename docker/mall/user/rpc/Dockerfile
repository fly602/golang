FROM alpine

USER root

COPY ./etc /etc/
COPY userrpc /usr/bin

WORKDIR /usr/bin/

EXPOSE 9000

ENTRYPOINT ["/usr/bin/userrpc", "-f","/etc/user.yaml"]