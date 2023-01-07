FROM go1.18

USER root

COPY ./etc /etc/
COPY payrpc /usr/bin

WORKDIR /usr/bin/

EXPOSE 9003

USER normal

ENTRYPOINT ["/usr/bin/payrpc", "-f","/etc/pay.yaml"]