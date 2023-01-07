FROM go1.18

USER root

COPY ./etc /etc/
COPY payapi /usr/bin

WORKDIR /usr/bin/

EXPOSE 8003

USER normal

ENTRYPOINT ["/usr/bin/payapi", "-f","/etc/pay.yaml"]