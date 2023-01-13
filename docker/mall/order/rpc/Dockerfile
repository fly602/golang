FROM go1.18

WORKDIR /opt/mall/order/rpc

COPY ./etc etc/
COPY orderrpc bin/

EXPOSE 9002

USER normal

ENTRYPOINT ["bin/orderrpc", "-f","etc/order.yaml"]