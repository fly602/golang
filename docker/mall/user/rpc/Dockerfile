FROM go1.18

USER root

COPY ./etc /etc/
COPY userrpc /usr/bin

WORKDIR /usr/bin/

EXPOSE 9000

USER normal

ENTRYPOINT ["/usr/bin/userrpc", "-f","/etc/user.yaml"]