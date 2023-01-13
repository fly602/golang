FROM go1.18

WORKDIR /opt/mall/product/rpc

COPY ./etc etc/
COPY productrpc bin/

EXPOSE 9001

USER normal

ENTRYPOINT ["bin/productrpc", "-f","etc/product.yaml"]