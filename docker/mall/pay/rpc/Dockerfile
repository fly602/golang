FROM go1.18

WORKDIR /opt/mall/pay/rpc

COPY ./etc etc/
COPY payrpc bin/

EXPOSE 9003

USER normal

ENTRYPOINT ["bin/payrpc", "-f","etc/pay.yaml"]